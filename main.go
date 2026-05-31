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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
)

var (
	isPreview       bool
	isVerbose       bool
	template        string
	sourceDirs      []string
	excludePatterns []string
	srcDir          string
	destDir         string
	processor       *Processor
	copyright       *Copyright
)

func main() {
	ret, ok := processCommandLine(os.Args[1:])
	if !ok {
		os.Exit(ret)
	}

	processor = NewProcessor()
	year := time.Now().Format("2006")
	copyright = NewCopyright(template, year)
	// for _, line := range copyright.copyright {
	// 	fmt.Println(line)
	// }
	searchDirectories()
}

// Process the given source directories for files to process.
func searchDirectories() {
	for _, dir := range sourceDirs {
		srcDir = filepath.Clean(dir)
		info, err := os.Stat(srcDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if !info.IsDir() {
			fmt.Fprintf(os.Stderr, "Invalid directory: %s\n", srcDir)
			continue
		}
		err = searchDirectory(srcDir)
		if err != nil {
			os.Exit(2)
		}
	}
}

// Traverse the given directory for files to process.
func searchDirectory(path string) error {
	if isVerbose {
		fmt.Printf("Processing directory %s ...\n", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		fullname := filepath.Join(path, entry.Name())
		if strings.HasPrefix(name, ".") {
			if isVerbose {
				fmt.Printf("Skipping %s\n", fullname)
			}
			continue
		}
		if entry.IsDir() {
			if IsExcluded(fullname) {
				if isVerbose {
					fmt.Printf("Skipping excluded directory: %s\n", fullname)
				}
				continue
			}

			err = searchDirectory(fullname)
			if err != nil {
				return err
			}
		} else {
			err = processor.ProcessFile(path, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}
		}
	}

	return nil
}

// Returns true if the path is excluded by the exclusions list.
func IsExcluded(path string) bool {
	for _, pattern := range excludePatterns {
		match, _ := doublestar.PathMatch(pattern, path)
		if match {
			return true
		}
	}
	return false
}

// processCommandLine processes the command line arguments and returns the
// true and the matcher string when successful.  It returns false if the
// -help option was specified, or if no matcher string specified
func processCommandLine(cmdLine []string) (int, bool) {
	helpFlag := flag.Bool("help", false, "print this message and exit")
	verboseFlag := flag.Bool("v", false, "set verbose logging")
	previewFlag := flag.Bool("p", false, "only list files that will be updated")
	destArg := flag.String("d", "", "destination directory (defaults to source)")
	templateArg := flag.String("t", ".copyright.txt", "a copyright template file (defaults to .copyright.txt)")
	excludedList := flag.String("e", "", "a list of directory patterns to exclude")

	flag.CommandLine.Parse(cmdLine)

	if *helpFlag {
		printUsage(os.Stdout)
		return 0, false
	}

	isVerbose = *verboseFlag
	isPreview = *previewFlag
	destDir = *destArg
	template = *templateArg
	populateExclusions(*excludedList)

	args := flag.Args()
	// at least 1 argument expected
	if len(args) < 1 {
		printUsage(os.Stderr)
		return 1, false
	}

	sourceDirs = args[0:]
	return 0, true
}

// Parses the excluded list argument and populates the excludePatterns slice.
func populateExclusions(excludedList string) {
	if len(excludedList) > 0 {
		for path := range strings.SplitSeq(excludedList, ",") {
			excluded, err := filepath.Localize(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid excluded path: %s\n", path)
				continue
			}
			excludePatterns = append(excludePatterns, excluded)
		}
	}
}

// printUsage prints the usage message for the program
func printUsage(stream *os.File) {
	fmt.Fprintln(stream, "usage: copyright [options] <source> ...")
	fmt.Fprintln(stream, "  -d string")
	fmt.Fprintln(stream, "\tdestination directory (defaults to overwrite")
	fmt.Fprintln(stream, "  -help")
	fmt.Fprintln(stream, "\tprint this message and exit")
	fmt.Fprintln(stream, "  -p\tonly list files that will be updated")
	fmt.Fprintln(stream, "  -t string")
	fmt.Fprintln(stream, "\ta copyright template file")
	fmt.Fprintln(stream, "  -v\tset verbose logging")
}
