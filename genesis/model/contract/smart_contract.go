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

package contract

import (
	"fmt"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
)

type BaseCompositeCollection struct {
	parent  object.Parent
	storage []object.Composite
}

func NewBaseCompositeCollection(parent object.Parent) *BaseCompositeCollection {
	return &BaseCompositeCollection{
		storage: make([]object.Composite, 0),
		parent:  parent,
	}
}

func (bcc *BaseCompositeCollection) GetClassID() string {
	return class.CompositeCollectionID
}

func (bcc *BaseCompositeCollection) GetParent() object.Parent {
	return bcc.parent
}

func (bcc *BaseCompositeCollection) GetList() []object.Composite {
	return bcc.storage
}

func (bcc *BaseCompositeCollection) Add(composite object.Composite) {
	bcc.storage = append(bcc.storage, composite)
}

func (bcc *BaseCompositeCollection) GetInterfaceKey() string {
	return class.CompositeCollectionID
}

type BaseCompositeCollectionProxy struct {
	BaseSmartContractProxy
}

func (*BaseCompositeCollectionProxy) GetInterfaceKey() string {
	return class.AllowanceID
}

func (cp *BaseCompositeCollectionProxy) GetList() []object.Composite {
	return cp.Instance.(object.CompositeCollection).GetList()
}

func (cp *BaseCompositeCollectionProxy) Add(composite object.Composite) {
	cp.Instance.(object.CompositeCollection).Add(composite)
}

// SmartContract marks that object is smart contract.
// TODO: Composite work interface
type SmartContract interface {
	object.Child
}

// InternalSmartContract contains private method.
type InternalSmartContract interface {
	object.Child
	GetResolver() object.Resolver
}

// BaseSmartContract is a base implementation of ComposingContainer, Callable and TypedObject interfaces.
type BaseSmartContract struct {
	object.BaseObject
	CompositeMap   map[string]object.Reference
	ChildStorage   storage.Storage
	ContextStorage storage.Storage
	Parent         object.Parent
	resolver       object.Resolver
	class          object.Proxy
}

// NewBaseSmartContract creates new BaseSmartContract instance with empty CompositeMap, ChildStorage and specific parent.
func NewBaseSmartContract(parent object.Parent, class object.Proxy) *BaseSmartContract {
	// TODO: NewCompositeHolder
	return &BaseSmartContract{
		CompositeMap: make(map[string]object.Reference),
		ChildStorage: storage.NewMapStorage(),
		Parent:       parent,
		class:        class,
	}
}

// GetResolver return instance or create it if necessary.
func (sc *BaseSmartContract) GetResolver() object.Resolver {
	if sc.resolver == nil {
		sc.resolver = object.NewResolverHandler(sc)
	}
	return sc.resolver
}

// GetClassID return string representation of object's class.
func (sc *BaseSmartContract) GetClassID() string {
	return class.SmartContractID
}

func (sc *BaseSmartContract) GetClass() object.Proxy {
	return sc.class
}

// CreateComposite allows to create composites inside smart contract.
func (sc *BaseSmartContract) CreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	composite, err := compositeFactory.Create(sc)
	if err != nil {
		return nil, err
	}

	interfaceKey := compositeFactory.GetInterfaceKey()
	_, exist := sc.CompositeMap[interfaceKey]
	if exist {
		return nil, fmt.Errorf("delegate with name %s already exist", interfaceKey)
	}

	compositeAsChild, ok := composite.(object.Child)
	if !ok {
		return nil, fmt.Errorf("composite is not a Child")
	}

	record, err := sc.AddChild(compositeAsChild)
	if err != nil {
		return nil, err
	}

	ref, err := object.NewReference("", record, object.ChildScope)
	if err != nil {
		// TODO: Remove child and do test
		return nil, err
	}
	sc.CompositeMap[interfaceKey] = ref

	return composite, nil
}

// GetComposite return composite by its key (if its exist inside smart contract).
func (sc *BaseSmartContract) GetComposite(key string, class object.CompositeFactory) (object.Composite, error) {
	ref, exist := sc.CompositeMap[key]
	if !exist {
		return nil, fmt.Errorf("delegate with name %s does not exist", key)
	}

	compObject, err := sc.GetResolver().GetObject(ref, class)
	if err != nil {
		return nil, err
	}

	composite, ok := compObject.(object.Composite)
	if !ok {
		return nil, fmt.Errorf("GetComposite: object is not Composite")
	}

	return composite, nil
}

// GetOrCreateComposite return composite by its key if its exist inside smart contract and create new one otherwise.
func (sc *BaseSmartContract) GetOrCreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	composite, err := sc.GetComposite(compositeFactory.GetInterfaceKey(), compositeFactory)
	if err != nil {
		composite, err = sc.CreateComposite(compositeFactory)
		if err != nil {
			return nil, err
		}
		return composite, nil
	}
	return composite, nil
}

// GetChildStorage return storage with children of smart contract.
func (sc *BaseSmartContract) GetChildStorage() storage.Storage {
	return sc.ChildStorage
}

// AddChild add new child to smart contract's ChildStorage.
func (sc *BaseSmartContract) AddChild(child object.Child) (string, error) {
	key, err := sc.ChildStorage.Set(child)
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetChild get child from smart contract's ChildStorage.
func (sc *BaseSmartContract) GetChild(key string) (object.Child, error) {
	child, err := sc.ChildStorage.Get(key)
	if err != nil {
		return nil, err
	}
	return child.(object.Child), nil
}

// GetContextStorage return storage with objects, which smart contract's children will have access to.
func (sc *BaseSmartContract) GetContextStorage() storage.Storage {
	return sc.ContextStorage
}

// GetContext return list of keys in ContextStorage.
func (sc *BaseSmartContract) GetContext() []string {
	return sc.GetContextStorage().GetKeys()
}

// GetParent return parent of smart contract.
func (sc *BaseSmartContract) GetParent() object.Parent {
	return sc.Parent
}
