/*
 *    Copyright 2018 Insolar
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

package hostnetwork

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/routing"
	"github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/pkg/errors"
)

// RelayRequest sends relay request to target.
func RelayRequest(hostHandler hosthandler.HostHandler, command, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	var typedCommand packet.CommandType
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("RelayRequest: target for relay request not found")
		return err
	}

	switch command {
	case "start":
		typedCommand = packet.StartRelay
	case "stop":
		typedCommand = packet.StopRelay
	default:
		err = errors.New("RelayRequest: unknown command")
		return err
	}
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeRelay).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestRelay{Command: typedCommand}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// CheckOriginRequest send a request to check target host originality
func CheckOriginRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("CheckOriginRequest: target for relay request not found")
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeCheckOrigin).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestCheckOrigin{}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Debugln(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// AuthenticationRequest sends an authentication request.
func AuthenticationRequest(hostHandler hosthandler.HostHandler, command, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("AuthenticationRequest: target for auth request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	var authCommand packet.CommandType
	switch command {
	case "begin":
		authCommand = packet.BeginAuthentication
	case "revoke":
		authCommand = packet.RevokeAuthentication
	default:
		err = errors.New("AuthenticationRequest: unknown command")
		return err
	}
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeAuthentication).
		Sender(origin).
		Receiver(targetHost).
		Request(&packet.RequestAuthentication{Command: authCommand}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Debugln(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// ObtainIPRequest is request to self IP obtaining.
func ObtainIPRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("ObtainIPRequest: target for relay request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeObtainIP).
		Sender(origin).
		Receiver(targetHost).
		Request(&packet.RequestObtainIP{}).
		Build()

	future, err := hostHandler.SendRequest(request)

	if err != nil {
		log.Debugln(err.Error())
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// RelayOwnershipRequest sends a relay ownership request.
func RelayOwnershipRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("relayOwnershipRequest: target for relay request not found")
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeRelayOwnership).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestRelayOwnership{Ready: true}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

// CascadeSendMessage sends a message to the next cascade layer.
func CascadeSendMessage(hostHandler hosthandler.HostHandler, data core.Cascade, targetID string, method string, args [][]byte) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("cascadeSendMessage: couldn't find a target host")
	}

	request := packet.NewBuilder().Sender(hostHandler.HtFromCtx(ctx).Origin).Receiver(targetHost).Type(packet.TypeCascadeSend).
		Request(&packet.RequestCascadeSend{
			Data: data,
			RPC: packet.RequestDataRPC{
				Method: method,
				Args:   args,
			},
		}).Build()

	future, err := hostHandler.SendRequest(request)
	if err != nil {
		return err
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func CheckPublicKeyRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return errors.Wrap(err, "failed to build a context")
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return errors.Wrap(err, "failed to find a target host")
	}
	if !exist {
		return errors.Wrap(err, "couldn't find a target host")
	}

	request := packet.NewBuilder().Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).Type(packet.TypeCheckPublicKey).
		Request(&packet.RequestCheckPublicKey{NodeID: hostHandler.GetNodeID(), HostID: hostHandler.GetOriginHost().IDs[0]}).
		Build()

	future, err := hostHandler.SendRequest(request)
	if err != nil {
		return errors.Wrap(err, "failed to send an authorization request")
	}
	return checkResponse(hostHandler, future, targetID, request)
}

// ResendPulseToKnownHosts resends received pulse to all known hosts
func ResendPulseToKnownHosts(hostHandler hosthandler.HostHandler, hosts []*routing.RouteHost, pulse *packet.RequestPulse) {
	for _, host1 := range hosts {
		err := sendPulse(hostHandler, host1.Host, pulse)
		if err != nil {
			log.Debug("error resending pulse to host %s: %s", host1.ID, err.Error())
		}
	}
}

func sendPulse(hostHandler hosthandler.HostHandler, host *host.Host, pulse *packet.RequestPulse) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return errors.Wrap(err, "failed to send pulse")
	}
	request := packet.NewBuilder().Sender(hostHandler.HtFromCtx(ctx).Origin).Receiver(host).
		Type(packet.TypePulse).Request(pulse).Build()

	future, err := hostHandler.SendRequest(request)
	if err != nil {
		return errors.Wrap(err, "failed to send pulse")
	}
	return checkResponse(hostHandler, future, "", request)
}

func checkNodePrivRequest(hostHandler hosthandler.HostHandler, targetID string) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return err
	}
	if !exist {
		err = errors.New("checkNodePrivRequest: target for check node privileges request not found")
		return err
	}

	origin := hostHandler.HtFromCtx(ctx).Origin
	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeCheckNodePriv).Sender(origin).Receiver(targetHost).Request(&packet.RequestCheckNodePriv{RoleKey: "test string"}).Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return errors.Wrap(err, "Failed to SendRequest")
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func knownOuterHostsRequest(hostHandler hosthandler.HostHandler, targetID string, hosts int) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}
	targetHost, exist, err := hostHandler.FindHost(ctx, targetID)
	if err != nil {
		return errors.Wrap(err, "Failed to FindHost")
	}
	if !exist {
		err = errors.New("knownOuterHostsRequest: target for relay request not found")
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeKnownOuterHosts).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(targetHost).
		Request(&packet.RequestKnownOuterHosts{
			ID:         hostHandler.HtFromCtx(ctx).Origin.ID.String(),
			OuterHosts: hosts},
		).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return errors.Wrap(err, "Failed to SendRequest")
	}

	return checkResponse(hostHandler, future, targetID, request)
}

func SendActiveNodesRequest(hostHandler hosthandler.HostHandler, target *host.Host) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeActiveNodes).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(target).
		Request(&packet.RequestActiveNodes{}).
		Build()
	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return errors.Wrap(err, "Failed to SendRequest")
	}

	return checkResponse(hostHandler, future, target.ID.String(), request)
}

// SendRelayOwnership send a relay ownership request.
func SendRelayOwnership(hostHandler hosthandler.HostHandler, subnetIDs []string) {
	for _, id1 := range subnetIDs {
		err := RelayOwnershipRequest(hostHandler, id1)
		log.Errorln(err.Error())
	}
}

func sendRelayedRequest(hostHandler hosthandler.HostHandler, request *packet.Packet) {
	_, err := hostHandler.SendRequest(request)
	if err != nil {
		log.Debugln(err)
	}
}

func sendCheckSignedNonceRequest(hostHandler hosthandler.HostHandler, target *host.Host, nonce []byte) error {
	ctx, err := NewContextBuilder(hostHandler).SetDefaultHost().Build()
	if err != nil {
		return err
	}

	builder := packet.NewBuilder()
	request := builder.Type(packet.TypeCheckSignedNonce).
		Sender(hostHandler.HtFromCtx(ctx).Origin).
		Receiver(target).
		Request(&packet.RequestCheckSignedNonce{}).
		Build()

	future, err := hostHandler.SendRequest(request)

	if err != nil {
		return errors.Wrap(err, "Failed to SendRequest")
	}

	return checkResponse(hostHandler, future, target.ID.String(), request)
}

func checkResponse(hostHandler hosthandler.HostHandler, future transport.Future, targetID string, request *packet.Packet) error {
	var err error
	rsp, err := future.GetResult(hostHandler.GetPacketTimeout())
	if err != nil {
		return errors.Wrap(err, "checkResponse error")
	}
	switch request.Type {
	case packet.TypeKnownOuterHosts:
		response := rsp.Data.(*packet.ResponseKnownOuterHosts)
		err = handleKnownOuterHosts(hostHandler, response, targetID)
	case packet.TypeCheckOrigin:
		response := rsp.Data.(*packet.ResponseCheckOrigin)
		handleCheckOriginResponse(hostHandler, response, targetID)
	case packet.TypeAuthentication:
		response := rsp.Data.(*packet.ResponseAuthentication)
		err = handleAuthResponse(hostHandler, response, targetID)
	case packet.TypeObtainIP:
		response := rsp.Data.(*packet.ResponseObtainIP)
		err = handleObtainIPResponse(hostHandler, response, targetID)
	case packet.TypeRelayOwnership:
		response := rsp.Data.(*packet.ResponseRelayOwnership)
		handleRelayOwnership(hostHandler, response, targetID)
	case packet.TypeCheckNodePriv:
		response := rsp.Data.(*packet.ResponseCheckNodePriv)
		err = handleCheckNodePrivResponse(hostHandler, response)
	case packet.TypeRelay:
		response := rsp.Data.(*packet.ResponseRelay)
		err = handleRelayResponse(hostHandler, response, targetID)
	case packet.TypeCascadeSend:
		response := rsp.Data.(*packet.ResponseCascadeSend)
		if !response.Success {
			err = errors.New(response.Error)
		}
	case packet.TypePulse:
		response := rsp.Data.(*packet.ResponsePulse)
		if !response.Success {
			err = errors.New(response.Error)
		}
	case packet.TypeCheckPublicKey:
		response := rsp.Data.(*packet.ResponseCheckPublicKey)
		err = handleCheckPublicKeyResponse(hostHandler, response)
		if err == nil {
			err = sendCheckSignedNonceRequest(hostHandler, rsp.Sender, response.Nonce)
		}
	case packet.TypeCheckSignedNonce:
		err = handleCheckSignedNonceResponse(hostHandler, rsp.Data.(*packet.ResponseCheckSignedNonce))
	case packet.TypeActiveNodes:
		response := rsp.Data.(*packet.ResponseActiveNodes)
		err = handleActiveNodesResponse(hostHandler, response)
	}
	return err
}
