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
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/stretchr/testify/assert"
)

type invalidScopeReference struct{}

func (r *invalidScopeReference) GetClassID() string {
	return class.ReferenceID
}

func (r *invalidScopeReference) GetRecord() string {
	return "145"
}

func (r *invalidScopeReference) GetDomain() string {
	return "123"
}

func (r *invalidScopeReference) GetScope() ScopeType {
	return ScopeType(10000)
}

func (r *invalidScopeReference) String() string {
	return fmt.Sprintf("#%s.#%s", "145", "123")
}

func (r *invalidScopeReference) GetReference() Reference {
	return r
}

func (r *invalidScopeReference) SetReference(ref Reference) {
}

func (r *invalidScopeReference) GetParent() Parent {
	return nil
}

func TestNewHandler(t *testing.T) {
	mockParent := &mockParentProxy{}
	handler := NewResolverHandler(mockParent)

	assert.Equal(t, &ResolverHandler{
		globalResolver: GlobalResolver,
		childResolver: &childResolver{
			parent: mockParent,
		},
		contextResolver: &contextResolver{
			parent: mockParent,
		},
	}, handler)
}

func TestHandler_GetObject_Not_Reference(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewResolverHandler(mockParent)

	obj, err := resolverHandler.GetObject("not reference", "mockChild")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestHandler_GetObject_GlobalScope(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewResolverHandler(nil)
	newMap := make(map[string]Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	ref, _ := NewReference("123", "1", GlobalScope)
	(*GlobalResolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolverHandler.GetObject(ref, factory)

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ChildScope(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewResolverHandler(mockParent)
	ref, _ := NewReference("1", "1", ChildScope)

	obj, err := resolverHandler.GetObject(ref, factory)

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ContextScope(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolverHandler := NewResolverHandler(mockParent)
	ref, _ := NewReference("1", record, ContextScope)

	obj, err := resolverHandler.GetObject(ref, factory)

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_default(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewResolverHandler(mockParent)
	ref := &invalidScopeReference{}

	obj, err := resolverHandler.GetObject(ref, factory)

	assert.EqualError(t, err, "unknown scope type: 10000")
	assert.Nil(t, obj)
}

func TestHandler_SetGlobalMap(t *testing.T) {
	resolverHandler := NewResolverHandler(nil)
	resolverHandler.globalResolver.globalInstanceMap = nil

	newMap := make(map[string]Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, resolverHandler.globalResolver.globalInstanceMap)
}
