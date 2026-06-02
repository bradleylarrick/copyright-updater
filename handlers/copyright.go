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
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	defaultTemplateFile = ".copyright.txt"
	currentYear         = time.Now().Format("2006")
)

type Copyright struct {
	Copyright []string
}

/*
 * Returns a new Copyright instance with the derived template.
 */
func NewCopyright(templateFile string, defaultTemplate []string, isVerbose bool) *Copyright {
	return &Copyright{Copyright: loadTemplate(templateFile, defaultTemplate, isVerbose)}
}

/*
 * Returns the copyright text as a slice of strings. If no previous copyright string is provided,
 * the ${year} placeholder is replaced with the current year. If a previous copyright string is provided,
 * the ${year} placeholder is replaced with '<previous>-<current year>'.
 */
func (Copyright) GetCopyright(copyright *Copyright, previous string) []string {
	var newYear string
	if previous != "" && !strings.EqualFold(previous, currentYear) {
		newYear = previous + "-" + currentYear
	} else {
		newYear = currentYear
	}
	var returnVal []string
	for _, line := range copyright.Copyright {
		returnVal = append(returnVal, strings.ReplaceAll(line, "${year}", newYear))
	}
	return returnVal
}

/*
 * Attempts to load the given copyright template. if no template is provided, the default template is used.
 * If the default template is not available, it uses the template from the global configuration file.
 */
func loadTemplate(templateFile string, defaultTemplate []string, isVerbose bool) []string {

	var copyright []string

	var file *os.File
	var err error
	// If a template file argument provided, fail if the file cannot be opened.
	if templateFile != "" {
		file, err = openFile(templateFile, true)
		if err != nil {
			os.Exit(4)
		}
	} else {
		// Don't fail if the default template file is not available; use the template from the global config instead.
		file, _ = openFile(defaultTemplateFile, false)
	}

	if file != nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			copyright = append(copyright, line)

			if scanner.Err() != nil {
				fmt.Fprintf(os.Stderr, "Failed: %s\n", scanner.Err().Error())
				os.Exit(1)
			}
		}
	} else {
		if isVerbose {
			fmt.Println("Loading copyright from global config...")
		}
		copyright = defaultTemplate
	}

	return copyright
}

/*
 * Opens the given file path. If an error occurs and failOnError is true, return the error.
 * Otherwise, return nil, nil.
 */
func openFile(path string, failOnError bool) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil && failOnError {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
		return nil, err
	}

	return file, nil
}
