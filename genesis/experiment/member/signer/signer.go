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

package signer

import (
	"github.com/ugorji/go/codec"
)

type msg struct {
	Ref    string
	Method string
	Params []interface{}
	Seed   []byte
}

func Serialize(ref string, method string, params []interface{}, seed []byte) ([]byte, error) {
	message := msg{ref, method, params, seed}
	var serialized []byte
	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&serialized, ch).Encode(message)
	if err != nil {
		return nil, err
	}
	return serialized, nil
}
