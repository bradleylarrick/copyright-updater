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
	"regexp"
	"strings"
)

var (
	scanner  *bufio.Scanner
	writer   *bufio.Writer
	header   string
	footer   string
	prefix   string
	lines    []string
	previous string
	ndx      int
)

/*
 * Starts the process of updating the copyright header in the source file.
 */
func startProcess(srcFile *os.File, destFile *os.File, hdr string, ftr string, pref string) error {
	header = hdr
	footer = ftr
	prefix = pref
	previous = ""
	ndx = 0
	lines = lines[:ndx] // clear the lines slice

	scanner = bufio.NewScanner(srcFile)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	writer = bufio.NewWriter(destFile)
	return nil
}

/*
 * Checks the first line against the list of protected lines and, if a match is found,
 * writes the line to the output file and moves to the next line.
 */
func findProtected(protectedList []string) {
	if ndx < len(lines) {
		for _, protected := range protectedList {
			if len(lines[ndx]) > 0 && strings.HasPrefix(lines[ndx], protected) {
				writeLine(lines[ndx])
				ndx++
				break
			}
		}
	}
}

/*
 * findHeader finds the start of the copyright header in the source file and,
 * if found, skips to the end of the header.
 */
func findHeader() {
	headerStart := false
	headerEnd := -1
	hasCopyright := false
	for i := ndx; i < len(lines); i++ {
		line := lines[i]
		if !headerStart {
			if len(line) == 0 {
				continue // skip blank lines before the header
			} else if strings.HasPrefix(line, header) {
				headerStart = true
			} else {
				break // No header block found
			}
		}

		if !hasCopyright {
			if strings.Contains(strings.ToLower(line), "copyright") {
				hasCopyright = true
				findPreviousYear(line)
			}
		}

		// We need to check for the footer line before checkning for a prefixed line to avoid
		//  matching comment lines as footers when the file has single-character comment prefixes.
		if isCommentFooter(i, line) {
			headerEnd = i
			break
		}
	}

	if headerEnd >= 0 && hasCopyright {
		ndx = headerEnd + 1
	}
}

/*
 * Writes the copyright header to the output file.
 */
func writeCopyright() error {
	text := GetCopyright(previous)
	writeLine(header)
	for _, line := range text {
		trimmed := strings.TrimRight(prefix+" "+line, " \t\n")
		writeLine(trimmed)
	}
	writeLine(footer)
	return nil
}

/*
 * Copies the remaining lines from the source file to the output file.
 */
func endProcess() error {
	for i := ndx; i < len(lines); i++ {
		line := lines[i]
		// Add a blank line after the header if the next line is a package declaration
		if i == ndx && strings.HasPrefix(line, "package") {
			writeLine("")
		}
		writeLine(line)
	}

	err := writer.Flush()
	if err != nil {
		return err
	}

	// Check for reading phase errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(4)
	}

	return nil
}

// Write a single line to the output file.
func writeLine(line string) {
	_, err := writer.WriteString(line + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing line: %s\n", err)
		os.Exit(4)
	}
}

/*
 * Load additional protected patterns.
 */
func addProtected(protected *[]string, new []string) {
	for _, pro := range new {
		*protected = append(*protected, pro)
	}
	// fmt.Printf("New Protected: %s\n", *protected)
}

/*
 * Finds the previous year in the given line and stores it in the previous variable.
 */
func findPreviousYear(line string) {
	re := regexp.MustCompile(`\b(\d{4})?-?(\d{4})\b`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 0 {
		if matches[1] != "" {
			previous = matches[1]
		} else {
			previous = matches[2]
		}
	}
}

/*
 * Returns true if the given line is a comment footer and the next line doesn't have a comment prefix.
 * This protects us in files with single-character comment prefixes.
 */
func isCommentFooter(index int, line string) bool {
	ret := false
	if strings.HasSuffix(line, strings.TrimSpace(footer)) {
		ret = true
		if strings.EqualFold(footer, prefix) {
			if index < len(lines)-1 {
				ret = !strings.HasPrefix(lines[index+1], prefix)
			}
		}
	}
	return ret
}

/*
 * Returns true if the given line is empty
 */
func isBlank(line string) bool {
	trimmed := strings.TrimRight(line, " ")
	return len(trimmed) == 0
}
