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
	"testing"

	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

type mockInstance struct{}

func (p *mockInstance) GetClassID() string {
	return "mockInstance"
}

func (p *mockInstance) GetClass() object.Proxy {
	return nil
}

func (p *mockInstance) GetParent() object.Parent {
	return nil
}

func TestBaseSmartContractProxy_GetClassID(t *testing.T) {
	proxy := &BaseSmartContractProxy{
		Instance: &mockInstance{},
	}
	assert.Equal(t, "mockInstance", proxy.GetClassID())
}

func TestBaseSmartContractProxy_GetParent(t *testing.T) {
	proxy := &BaseSmartContractProxy{
		Instance: &mockInstance{},
	}
	assert.Nil(t, proxy.GetParent())
}

func TestBaseSmartContractProxy_GetOrCreateComposite(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy := &BaseSmartContractProxy{
		Instance: NewBaseSmartContract(parent, factory),
	}
	compositeFactory := &BaseCompositeFactory{}
	composite := &BaseComposite{
		class: compositeFactory,
	}

	res, err := proxy.GetOrCreateComposite(compositeFactory)

	assert.NoError(t, err)
	assert.Equal(t, composite, res)
}

func TestBaseSmartContractProxy_GetOrCreateComposite_Error(t *testing.T) {
	parent := &mockParent{}
	proxy := &BaseSmartContractProxy{
		Instance: &mockChildProxy{
			parent: parent,
		},
	}

	compositeFactory := &BaseCompositeFactory{}

	_, err := proxy.GetOrCreateComposite(compositeFactory)

	assert.EqualError(t, err, "Instance is not Composing Container")
}
