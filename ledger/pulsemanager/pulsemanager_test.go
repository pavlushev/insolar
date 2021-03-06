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

package pulsemanager_test

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/stretchr/testify/assert"
)

func TestPulseManager_Current(t *testing.T) {
	ledger, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	var err error

	pm := ledger.GetPulseManager()

	pulse, err := pm.Current()
	assert.NoError(t, err)
	assert.Equal(t, *core.GenesisPulse, *pulse)
}
