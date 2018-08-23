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

package packet

import (
	"github.com/insolar/insolar/network/hostnetwork/host"
)

// Builder allows lazy building of packets.
// Each operation returns new copy of a builder.
type Builder struct {
	actions []func(packet *Packet)
}

// NewBuilder returns empty packet builder.
func NewBuilder() Builder {
	return Builder{}
}

// Build returns configured packet.
func (cb Builder) Build() (packet *Packet) {
	packet = &Packet{}
	for _, action := range cb.actions {
		action(packet)
	}
	return
}

// Sender sets packet sender.
func (cb Builder) Sender(node *host.Host) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Sender = node
	})
	return cb
}

// Receiver sets packet receiver.
func (cb Builder) Receiver(node *host.Host) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Receiver = node
	})
	return cb
}

// Type sets packet type.
func (cb Builder) Type(packetType packetType) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Type = packetType
	})
	return cb
}

// Request adds request data to packet.
func (cb Builder) Request(request interface{}) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Data = request
	})
	return cb
}

// Response adds response data to packet
func (cb Builder) Response(response interface{}) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Data = response
		packet.IsResponse = true
	})
	return cb
}

// Error adds error description to packet.
func (cb Builder) Error(err error) Builder {
	cb.actions = append(cb.actions, func(packet *Packet) {
		packet.Error = err
	})
	return cb
}
