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

import "os"

type VmHandler struct{}

var (
	vmHeader    = "##"
	vmFooter    = "##"
	vmPrefix    = "##"
	vmProtected = []string{}
)

func (VmHandler) Format(src *os.File, dest *os.File) error {
	err := startProcess(src, dest, vmHeader, vmFooter, vmPrefix)
	if err != nil {
		return err
	}

	if len(vmProtected) > 0 {
		findProtected(vmProtected)
	}
	findHeader()
	writeCopyright()
	return endProcess()
}

/*
 * Load additional protected patterns.
 */
func (VmHandler) AddProtected(protected []string) {
	addProtected(&hashProtected, protected)
}
