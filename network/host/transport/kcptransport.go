/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package transport

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/network/host/message"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/xtaci/kcp-go"
)

type kcpTransport struct {
	listener *kcp.Listener

	received chan *message.Message
	sequence *uint64

	disconnectStarted  chan bool
	disconnectFinished chan bool

	mutex   *sync.RWMutex
	futures map[message.RequestID]Future

	proxy      relay.Proxy
	blockCrypt kcp.BlockCrypt
}

// NewKCPTransport creates kcpTransport.
func NewKCPTransport(conn net.PacketConn, proxy relay.Proxy) (Transport, error) {
	return newKCPTransport(conn, proxy)
}

func newKCPTransport(conn net.PacketConn, proxy relay.Proxy) (*kcpTransport, error) {
	crypt, err := kcp.NewNoneBlockCrypt([]byte{})

	if err != nil {
		return nil, err
	}

	lis, err := kcp.ServeConn(crypt, 0, 0, conn)
	if err != nil {
		return nil, err
	}

	transport := &kcpTransport{
		listener: lis,

		received: make(chan *message.Message),
		sequence: new(uint64),

		disconnectStarted:  make(chan bool),
		disconnectFinished: make(chan bool),

		mutex:   &sync.RWMutex{},
		futures: make(map[message.RequestID]Future),

		proxy:      proxy,
		blockCrypt: crypt,
	}

	return transport, nil
}

// SendRequest sends request message and returns future.
func (t *kcpTransport) SendRequest(msg *message.Message) (Future, error) {
	if !msg.IsValid() {
		return nil, errors.New("invalid message")
	}

	msg.RequestID = t.generateID()

	future := t.createFuture(msg)

	err := t.sendMessage(msg)
	if err != nil {
		future.Cancel()
		return nil, err
	}

	return future, nil
}

// SendResponse sends response message.
func (t *kcpTransport) SendResponse(requestID message.RequestID, msg *message.Message) error {
	msg.RequestID = requestID

	return t.sendMessage(msg)
}

// Start starts networking.
func (t *kcpTransport) Start() error {
	for {
		if session, err := t.listener.AcceptKCP(); err == nil {
			go t.handleAcceptedConnection(session)
		} else {
			<-t.disconnectFinished
			return err
		}
	}
}

// Stop stops networking.
func (t *kcpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	err := t.listener.Close()
	if err != nil {
		log.Println("Failed to close socket:", err.Error())
	}

	t.disconnectStarted <- true
	close(t.disconnectStarted)
}

// Close closes message channels.
func (t *kcpTransport) Close() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	close(t.received)
	close(t.disconnectFinished)
}

// Messages returns incoming messages channel.
func (t *kcpTransport) Messages() <-chan *message.Message {
	return t.received
}

// Stopped checks if networking is stopped already.
func (t *kcpTransport) Stopped() <-chan bool {
	return t.disconnectStarted
}

func (t *kcpTransport) socketDialTimeout(addr string, timeout time.Duration) (*kcp.UDPSession, error) {
	return kcp.DialWithOptions(addr, t.blockCrypt, 0, 0)
}

func (t *kcpTransport) generateID() message.RequestID {
	id := AtomicLoadAndIncrementUint64(t.sequence)
	return message.RequestID(id)
}

func (t *kcpTransport) createFuture(msg *message.Message) Future {
	newFuture := NewFuture(msg.RequestID, msg.Receiver, msg, func(f Future) {
		t.mutex.Lock()
		defer t.mutex.Unlock()

		delete(t.futures, f.Request().RequestID)
	})

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.futures[msg.RequestID] = newFuture

	return newFuture
}

func (t *kcpTransport) getFuture(msg *message.Message) Future {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.futures[msg.RequestID]
}

func (t *kcpTransport) sendMessage(msg *message.Message) error {
	var recvAddress string
	if t.proxy.ProxyNodesCount() > 0 {
		recvAddress = t.proxy.GetNextProxyAddress()
	}
	if len(recvAddress) == 0 {
		recvAddress = msg.Receiver.Address.String()
	}

	session, err := t.socketDialTimeout(recvAddress, time.Second)
	if err != nil {
		return err
	}

	data, err := message.SerializeMessage(msg)
	if err != nil {
		return err
	}

	_, err = session.Write(data)

	return err
}

func (t *kcpTransport) getRemoteAddress(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func (t *kcpTransport) handleAcceptedConnection(session *kcp.UDPSession) {
	for {
		err := session.SetDeadline(time.Now().Add(time.Millisecond * 50))
		if err != nil {
			log.Println(err.Error())
		}
		// Wait for Messages
		msg, err := message.DeserializeMessage(session)
		if err != nil {
			// TODO should we penalize this Node somehow ? Ban it ?
			// if err.Error() != "EOF" {
			// }
			return
		}
		msg.RemoteAddress = t.getRemoteAddress(session)
		t.handleMessage(msg)
	}
}

func (t *kcpTransport) handleMessage(msg *message.Message) {
	if msg.IsResponse {
		t.processResponse(msg)
	} else {
		t.processRequest(msg)
	}
}

func (t *kcpTransport) processResponse(msg *message.Message) {
	future := t.getFuture(msg)
	if future != nil && !shouldProcessMessage(future, msg) {
		future.SetResult(msg)
	}
	future.Cancel()
}

func (t *kcpTransport) processRequest(msg *message.Message) {
	if msg.IsValid() {
		t.received <- msg
	}
}