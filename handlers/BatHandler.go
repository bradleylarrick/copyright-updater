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

type BatHandler struct{}

var (
	batHeader    = "REM"
	batFooter    = "REM"
	batPrefix    = "REM"
	BatProtected = []string{"@echo"}
)

func (BatHandler) Format(src *os.File, dest *os.File) error {
	err := startProcess(src, dest, batHeader, batFooter, batPrefix)
	if err != nil {
		return err
	}

	if len(BatProtected) > 0 {
		findProtected(BatProtected)
	}
	findHeader()
	writeCopyright()
	return endProcess()
}

/*
 * Load additional protected patterns.
 */
func (BatHandler) AddProtected(protected []string) {
	addProtected(&BatProtected, protected)
}
