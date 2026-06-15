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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateExclusions(t *testing.T) {
	excludedList := "LICENSE*,test/*,temp/*"
	populateExclusions(excludedList)
	assert.Equal(t, 1, len(excludedPaths))
	assert.Equal(t, 2, len(excludedDirs))
	assert.Equal(t, `LICENSE*`, excludedPaths[0])
	assert.Equal(t, `test`, excludedDirs[0])
	assert.Equal(t, `temp`, excludedDirs[1])
}

func TestIsExcluded(t *testing.T) {
	tests := []struct {
		name     string
		fullSrc  string
		pattern  *[]string
		expected bool
	}{
		{
			name:     "LICENSE file",
			fullSrc:  "LICENSE",
			pattern:  &excludedPaths,
			expected: true,
		},
		{
			name:     "LICENSE.md file",
			fullSrc:  "LICENSE.md",
			pattern:  &excludedPaths,
			expected: true,
		},
		{
			name:     "LICENSE in test directory",
			fullSrc:  `test\LICENSE`,
			pattern:  &excludedPaths,
			expected: false,
		},
		{
			name:     "test directory",
			fullSrc:  "test",
			pattern:  &excludedDirs,
			expected: true,
		},
		{
			name:     "temp directory",
			fullSrc:  "temp",
			pattern:  &excludedDirs,
			expected: true,
		},
		{
			name:     "src directory",
			fullSrc:  "src",
			pattern:  &excludedDirs,
			expected: false,
		},
	}

	excludedList := "LICENSE*,test/*,temp/*"
	populateExclusions(excludedList)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsExcluded(test.fullSrc, *test.pattern))
		})
	}
}
