/*
 * Copyright (c) 2026 Bradley Larrick. All rights reserved.
 *
 * Licensed under the Apache License v2.0
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/zenizh/go-capturer"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	processor = NewProcessor()
	result := m.Run()
	os.Exit(result)
}

func TestPopulateExclusions(t *testing.T) {
	excludedList := "LICENSE*,test/*,temp/*,../../tester/*"

	output := capturer.CaptureOutput(func() {
		populateExclusions(excludedList)
	})

	lines := strings.Split(output, "\n")
	assert.True(t, strings.HasPrefix(lines[0], "Invalid excluded path"))

	assert.Contains(t, excludedPaths, `LICENSE*`)
	assert.Contains(t, excludedDirs, `test`)
	assert.Contains(t, excludedDirs, `temp`)
	assert.NotContains(t, excludedDirs, `../../tester`)
}
