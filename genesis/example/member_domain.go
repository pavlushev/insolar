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

package example

import (
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
)

// MemberDomainName is a name for member domain.
const MemberDomainName = "MemberDomain"

// MemberDomain is a contract that allows to add new members to system.
type MemberDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateMember is used to create new member as a child to domain storage.
	CreateMember() (string, error)
	// GetMember returns member from its record in domain storage.
	GetMember(string) (Member, error)
}

type memberDomain struct {
	domain.BaseDomain
	memberFactoryReference object.Reference
}

// newMemberDomain creates new instance of MemberDomain.
func newMemberDomain(parent object.Parent, class object.Factory) (MemberDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	mf, _ := NewMemberFactory(parent).(*memberFactory)

	md := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, class, MemberDomainName),
	}
	// Add memberFactory as a child
	record, err := md.AddChild(mf)
	if err != nil {
		return nil, err
	}
	md.memberFactoryReference, err = object.NewReference("", record, object.ChildScope)
	if err != nil {
		return nil, err
	}
	return md, nil
}

// GetClassID returns string representation of MemberDomain's class.
func (md *memberDomain) GetClassID() string {
	return class.MemberDomainID
}

// CreateMember creates new member as a child to domain storage.
func (md *memberDomain) CreateMember() (string, error) {
	// Get child by memberFactoryRecord
	r := md.GetResolver()
	child, err := r.GetObject(md.memberFactoryReference, nil)
	if err != nil {
		return "", err
	}
	// Check if it Factory
	mf, ok := child.(object.Factory)
	if !ok {
		return "", fmt.Errorf("child by reference `%s` is not Factory instance", md.memberFactoryReference)
	}
	// Create member
	member, err := mf.Create(md)
	if err != nil {
		return "", err
	}
	if member == nil {
		return "", fmt.Errorf("factory returns nil")
	}

	record, err := md.AddChild(member)
	if err != nil {
		return "", err
	}

	return record, nil
}

// GetMember returns member from its record in domain storage.
func (md *memberDomain) GetMember(record string) (Member, error) {
	r := md.GetResolver()
	ref, err := object.NewReference("", record, object.ChildScope)
	if err != nil {
		return nil, err
	}
	child, err := r.GetObject(md.memberFactoryReference, nil)
	if err != nil {
		return nil, err
	}
	// Check if it Factory
	mf, ok := child.(object.Factory)
	if !ok {
		return nil, fmt.Errorf("child by reference `%s` is not Factory instance", md.memberFactoryReference)
	}

	member, err := r.GetObject(ref, mf)
	if err != nil {
		return nil, err
	}
	return member.(Member), nil
}

type memberDomainProxy struct {
	contract.BaseSmartContractProxy
}

// newMemberDomainProxy creates new proxy and associates it with new instance of MemberDomain.
func newMemberDomainProxy(parent object.Parent, class object.Factory) (*memberDomainProxy, error) {
	instance, err := newMemberDomain(parent, class)
	if err != nil {
		return nil, err
	}
	return &memberDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: instance,
		},
	}, nil
}

// CreateMember is a proxy call for instance method.
func (mdp *memberDomainProxy) CreateMember() (string, error) {
	return mdp.Instance.(MemberDomain).CreateMember()
}

// GetMember is a proxy call for instance method.
func (mdp *memberDomainProxy) GetMember(record string) (Member, error) {
	return mdp.Instance.(MemberDomain).GetMember(record)
}

type memberDomainFactory struct {
	object.BaseFactory
	parent object.Parent
}

// NewMemberDomainFactory creates new factory for MemberDomain.
func NewMemberDomainFactory(parent object.Parent) object.Factory {
	return &memberDomainFactory{
		parent: parent,
	}
}

// GetClassID returns string representation of MemberDomain's class.
func (mdf *memberDomainFactory) GetClassID() string {
	return class.MemberDomainID
}

func (mdf *memberDomainFactory) GetClass() object.Proxy {
	return mdf
}

// GetParent returns parent
func (mdf *memberDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// Create is a factory method for new MemberDomain instances.
func (mdf *memberDomainFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy, err := newMemberDomainProxy(parent, mdf)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
