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

type XmlHandler struct{}

var (
	xmlHeader    = "<!--"
	xmlFooter    = " -->"
	xmlPrefix    = " "
	xmlProtected = []string{"<?xml version", "<!DOCTYPE"}
)

func (XmlHandler) Format(src *os.File, dest *os.File, copyright *Copyright) error {
	err := startProcess(src, dest, xmlHeader, xmlFooter, xmlPrefix)
	if err != nil {
		return err
	}

	if len(xmlProtected) > 0 {
		findProtected(xmlProtected)
	}
	findHeader()
	writeCopyright(copyright)
	return endProcess()
}
