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

package preprocessor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/stretchr/testify/assert"
)

var contractNames = []string{"wallet", "member", "allowance"}
var pathWithContracts = "../../../genesis/experiment/"

func contractPath(name string) string {
	return pathWithContracts + name + "/" + name + ".go"
}

func MakeTestName(file string, contractType string) string {
	return fmt.Sprintf("Generate contract %s from '%s'", contractType, file)
}

func TestGenerateProxiesForRealSmartContracts(t *testing.T) {
	for _, name := range contractNames {
		file := contractPath(name)
		t.Run(MakeTestName(file, "proxy"), func(t *testing.T) {
			var buf bytes.Buffer
			err := GenerateContractProxy(file, "testRef", &buf)
			assert.NoError(t, err)

			code, err := ioutil.ReadAll(&buf)
			assert.NoError(t, err)
			assert.NotEqual(t, len(code), 0)
		})
	}
}

func TestGenerateWrappersForRealSmartContracts(t *testing.T) {
	for _, name := range contractNames {
		file := contractPath(name)
		t.Run(MakeTestName(file, "wrapper"), func(t *testing.T) {
			var buf bytes.Buffer
			err := GenerateContractWrapper(file, &buf)
			assert.NoError(t, err)

			code, err := ioutil.ReadAll(&buf)
			assert.NoError(t, err)
			assert.NotEqual(t, len(code), 0)
		})
	}
}

func TestCompilingRealSmartContracts(t *testing.T) {
	iccDir := "../../../cmd/icc"

	_, err := exec.Command("go", "build", "-o", iccDir+"/icc", iccDir).CombinedOutput()
	assert.NoError(t, err)

	contracts := make(map[string]string)
	for _, name := range contractNames {
		code, err := ioutil.ReadFile(contractPath(name))
		assert.NoError(t, err)
		contracts[name] = string(code)
	}

	am := testutil.NewTestArtifactManager()
	cb := testutil.NewContractBuilder(am, iccDir+"/icc")
	err = cb.Build(contracts)
	assert.NoError(t, err)
}
