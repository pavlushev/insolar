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

package core

import (
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
)

// InstanceDomainName is a name for instance domain.
const InstanceDomainName = "InstanceDomain"

// InstanceDomain is a contract that stores instances of other domains
type InstanceDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateInstance is used to create new instance as a child to domain storage.
	CreateInstance(object.Factory) (string, error)
	// GetInstance returns instance from its record in domain storage.
	GetInstance(string) (object.Proxy, error)
}

type instanceDomain struct {
	domain.BaseDomain
}

// newInstanceDomain creates new instance of InstanceDomain
func newInstanceDomain(parent object.Parent, class object.Factory) (*instanceDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	instDomain := &instanceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, class, InstanceDomainName),
	}
	return instDomain, nil
}

// GetClassID returns string representation of InstanceDomain's class.
func (instDom *instanceDomain) GetClassID() string {
	return class.InstanceDomainID
}

// CreateInstance creates new instance as a child to domain storage.
func (instDom *instanceDomain) CreateInstance(fc object.Factory) (string, error) {
	instance, err := fc.Create(instDom)
	if err != nil {
		return "", err
	}
	if instance == nil {
		return "", fmt.Errorf("factory returns nil")
	}

	record, err := instDom.ChildStorage.Set(instance)
	if err != nil {
		return "", err
	}

	return record, nil
}

// GetInstance returns instance from its record in domain storage.
func (instDom *instanceDomain) GetInstance(record string) (object.Proxy, error) {
	instance, err := instDom.ChildStorage.Get(record)
	if err != nil {
		return nil, err
	}

	result, ok := instance.(object.Proxy)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `Proxy` instance", record)
	}

	return result, nil
}

type instanceDomainProxy struct {
	contract.BaseSmartContractProxy
}

// newInstanceDomainProxy creates new proxy and associate it with new instance of InstanceDomain.
func newInstanceDomainProxy(parent object.Parent, class object.Factory) (*instanceDomainProxy, error) {
	instance, err := newInstanceDomain(parent, class)
	if err != nil {
		return nil, err
	}
	return &instanceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: instance,
		},
	}, nil
}

// CreateInstance is a proxy call for instance method.
func (idp *instanceDomainProxy) CreateInstance(fc object.Factory) (string, error) {
	return idp.Instance.(InstanceDomain).CreateInstance(fc)
}

// GetInstance is a proxy call for instance method.
func (idp *instanceDomainProxy) GetInstance(record string) (object.Proxy, error) {
	return idp.Instance.(InstanceDomain).GetInstance(record)
}

type instanceDomainFactory struct {
	object.BaseFactory
	parent object.Parent
}

// NewInstanceDomainFactory creates new factory for InstanceDomain.
func NewInstanceDomainFactory(parent object.Parent) object.Factory {
	return &instanceDomainFactory{
		parent: parent,
	}
}

// GetParent returns parent link
func (idf *instanceDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// GetClassID returns string representation of InstanceDomain's class.
func (idf *instanceDomainFactory) GetClassID() string {
	return class.InstanceDomainID
}

func (idf *instanceDomainFactory) GetClass() object.Proxy {
	return idf
}

// Create is factory method that used to create new InstanceDomain instances.
func (idf *instanceDomainFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy, err := newInstanceDomainProxy(parent, idf)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
