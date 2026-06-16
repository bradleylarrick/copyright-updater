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
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	handlers "natuna.org/copyright/handlers"
)

func TestLoadExtensions(t *testing.T) {
	config := &Configuration{
		Extensions: []Extension{
			{Extension: ".abc", Processor: "apt", Protected: []string{"intro"}},
			{Extension: ".xyz", Processor: "java", Protected: []string{"package"}},
			{Extension: ".zyx", Processor: "zyx", Protected: []string{}},
			{Extension: ".123", Processor: "bat", Protected: []string{}},
			{Extension: ".456", Processor: "hashtag", Protected: []string{}},
			{Extension: ".789", Processor: "xml", Protected: []string{}},
		},
	}

	b, err := toml.Marshal(config)
	assert.NoError(t, err)
	fmt.Println(string(b))

	loadExtensions(config)
	assert.Equal(t, handlers.AptHandler{}, processor.Handlers[".abc"])
	assert.Equal(t, []string{"intro"}, handlers.AptProtected)
	assert.Equal(t, handlers.JavaHandler{}, processor.Handlers[".xyz"])
	assert.Equal(t, []string{"package"}, handlers.JavaProtected)
	assert.Equal(t, handlers.BatHandler{}, processor.Handlers[".123"])
	assert.Equal(t, handlers.HashtagHandler{}, processor.Handlers[".456"])
	assert.Equal(t, handlers.XmlHandler{}, processor.Handlers[".789"])
}

func TestLoadExclusions(t *testing.T) {
	config := &Configuration{
		Exclusions: []string{"path/to/exclude", "dir/to/exclude/*"},
	}

	loadExclusions(config)
	testPath, _ := filepath.Localize("path/to/exclude")
	assert.Contains(t, excludedPaths, testPath)
	testPath, _ = filepath.Localize("dir/to/exclude")
	assert.Contains(t, excludedDirs, testPath)
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
