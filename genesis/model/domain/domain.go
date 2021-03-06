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

package domain

import (
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/object"
)

// Domain is a collection of objects and related policies (construction, referencing, logical consensus etc).
type Domain interface {
	contract.SmartContract
}

type InternalDomain interface {
	contract.InternalSmartContract
}

// BaseDomain is a base implementation of Domain interface.
type BaseDomain struct {
	contract.BaseSmartContract
	//class object.Factory
	Name string
}

// NewBaseDomain creates new BaseDomain instance.
func NewBaseDomain(parent object.Parent, class object.Factory, name string) *BaseDomain {
	return &BaseDomain{
		BaseSmartContract: *contract.NewBaseSmartContract(parent, class.(object.Proxy)),
		//class:             class,
		Name: name,
	}
}

// GetClassID return string representation of object's class.
func (d *BaseDomain) GetClassID() string {
	return class.DomainID
}

/*func (d *BaseDomain) GetClass() object.Factory {
	return d.class
}*/

// GetName return name of domain.
func (d *BaseDomain) GetName() string {
	return d.Name
}
