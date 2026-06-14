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
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindExtension(t *testing.T) {
	tests := []struct {
		name     string
		fullSrc  string
		expected string
	}{
		{
			name:     "example.sh",
			fullSrc:  "test/example.sh",
			expected: ".sh",
		},
		{
			name:     "WithHeader.java",
			fullSrc:  "test/src/main/WithHeader.java",
			expected: ".java",
		},
		{
			name:     "example.txt",
			fullSrc:  "test/example.txt",
			expected: ".txt",
		},
		{
			name:     "Makefile",
			fullSrc:  "test/Makefile",
			expected: ".mk",
		},
		{
			name:     "Jenkinsfile",
			fullSrc:  "test/Jenkinsfile",
			expected: ".java",
		},
		{
			name:     "example.txt.vm",
			fullSrc:  "test/example.txt.vm",
			expected: ".vm",
		},
		{
			name:     "loader",
			fullSrc:  "test/loader",
			expected: ".sh",
		},
		{
			name:     "index.apt.vm",
			fullSrc:  "test/src/site/apt/index.apt.vm",
			expected: ".apt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := findExtension(test.name, test.fullSrc)
			assert.Equal(t, test.expected, result, "findExtension(%s) = %s; want %s", test.fullSrc, result, test.expected)
		})
	}
}

func TestValidateDestPath(t *testing.T) {
	tests := []struct {
		name     string
		dest     string
		expected error
	}{
		{
			name:     "existing destination",
			dest:     "handlers",
			expected: nil,
		},
		{
			name:     "new valid destination",
			dest:     "temp",
			expected: nil,
		},
		{
			name:     "invalid destination",
			dest:     "main.go/temp",
			expected: &os.PathError{Op: "mkdir", Path: "main.go", Err: syscall.ERROR_PATH_NOT_FOUND},
		},
	}

	os.RemoveAll("temp")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validateDestPath(test.dest)
			assert.Equal(t, test.expected, result, "validateDestPath(%s) = %t; want %t", test.dest, result, test.expected)
		})
	}
}
