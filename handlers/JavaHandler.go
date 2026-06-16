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

type JavaHandler struct{}

var (
	javaHeader    = "/*"
	javaFooter    = " */"
	javaPrefix    = " *"
	JavaProtected = []string{}

	slashes = "//"
)

func (JavaHandler) Format(src *os.File, dest *os.File) error {
	err := startProcess(src, dest, javaHeader, javaFooter, javaPrefix)
	if err != nil {
		return err
	}

	if strings.HasPrefix(lines[ndx], slashes) {
		setSlashPrefixes()
	}

	if len(JavaProtected) > 0 {
		findProtected(JavaProtected)
	}

	findHeader()
	writeCopyright()
	return endProcess()
}

/*
 * Load additional protected patterns.
 */
func (JavaHandler) AddProtected(protected []string) {
	addProtected(&JavaProtected, protected)
}

/*
 * Set the prefixes to slash prefixes.
 */
func setSlashPrefixes() {
	header = slashes
	footer = slashes
	prefix = slashes
}
