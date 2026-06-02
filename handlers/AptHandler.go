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

package handlers

import (
	"os"
	"strings"
)

type AptHandler struct{}

var (
	aptHeader    = "~~"
	aptFooter    = "~~"
	aptPrefix    = "~~"
	aptProtected = []string{}
)

func (AptHandler) Format(src *os.File, dest *os.File) error {
	err := startProcess(src, dest, aptHeader, aptFooter, aptPrefix)
	if err != nil {
		return err
	}

	if len(aptProtected) > 0 {
		findProtected(aptProtected)
	}

	skipAptHeader()
	findHeader()
	writeCopyright()
	return endProcess()
}

/*
 * Load additional protected patterns.
 */
func (AptHandler) AddProtected(protected []string) {
	addProtected(&aptProtected, protected)
}

/*
 * APT files can have a custom header that needs to be skipped. We look for the first blank line
 * or comment line and keep everything above it.
 */
func skipAptHeader() {
	for i := ndx; i < len(lines); i++ {
		line := lines[i]
		if isBlank(line) || strings.HasPrefix(line, aptHeader) {
			ndx = i
			break
		} else {
			writeLine(line)
		}
	}
}
