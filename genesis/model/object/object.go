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

package object

import (
	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
)

// Object marks that instance has ClassID (string representation of class).
type Object interface {
	GetClassID() string
	GetClass() Proxy
}

// BaseObject is a base implementation of Object interface.
type BaseObject struct {
}

// GetClassID return string representation of object's class.
func (bo *BaseObject) GetClassID() string {
	return class.ObjectID
}

func (bo *BaseObject) GetClass() Proxy {
	return nil
}

// Parent allows to create objects (smart contracts) inside itself as children.
type Parent interface {
	Callable
	GetChildStorage() storage.Storage     // Storage for child references
	AddChild(child Child) (string, error) // return key for GetChild func
	GetChild(key string) (Child, error)   // child type reference
	GetContext() []string                 // Parent give information about context references to its children
	GetContextStorage() storage.Storage   // Storage for context references
}

// Child allows to be created inside object (smart contract).
type Child interface {
	Callable
	GetParent() Parent
}
