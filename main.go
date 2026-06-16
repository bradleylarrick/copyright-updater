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
	"runtime"
	"strings"

	handlers "natuna.org/copyright/handlers"
)

var (
	version       string = "1.3.0 (" + runtime.GOOS + " " + runtime.GOARCH + ")"
	isPreview     bool
	isVerbose     bool
	templateFile  string
	sources       []string
	excludedDirs  []string
	excludedPaths []string
	source        string
	destDir       string
	processor     *Processor
)

/*
 * The Usage function prints the command line usage message.
 */
var Usage = func() {
	fmt.Fprintln(flag.CommandLine.Output(), "usage: copyright [options] <source> ...")
	flag.PrintDefaults()
}

func main() {
	ret, ok := processCommandLine(os.Args[1:])
	if !ok {
		os.Exit(ret)
	}

	// displayMemory()
	processor = NewProcessor()
	configure()
	processSources()
	// displayMemory()
}

/*
 * Load the configuration file and set up the copyright template.
 */
func configure() {
	err := loadConfigurationFile()
	if err != nil {
		os.Exit(1)
	}

	template := handlers.LoadTemplate(templateFile, Config.Copyright, isVerbose)
	if template == nil {
		fmt.Fprintln(os.Stderr, "No copyright template found.")
		os.Exit(1)
	} else if isVerbose {
		fmt.Println("Copyright template:")
		for _, line := range template {
			fmt.Println(line)
		}
		fmt.Println("Handlers:")
		fmt.Println(processor.Handlers)
	}
}

// Process the given source directories for files to process.
func processSources() {
	for _, src := range sources {
		source = filepath.Clean(src)
		info, err := os.Stat(source)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if info.IsDir() {
			err = searchDirectory(source)
		} else {
			dir, file := filepath.Split(source)
			err = processor.ProcessFile(dir, file)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}
}

// Traverse the given directory for files to process.
func searchDirectory(path string) error {
	if IsExcluded(path, excludedDirs) {
		if isVerbose {
			fmt.Printf("Skipping excluded directory: %s\n", path)
		}
		return nil
	}

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

/*
 *  Processes the command line arguments and returns a flag
 * indicating success and a exit value to use of not.
 */
func processCommandLine(cmdLine []string) (int, bool) {
	flag.NewFlagSet("copyright", flag.ExitOnError)
	flag.Usage = Usage
	destArg := flag.String("d", "", "destination `directory` (defaults to source)")
	excludedList := flag.String("e", "", "a list of directory `patterns` to exclude")
	helpFlag := flag.Bool("h", false, "print this message and exit")
	previewFlag := flag.Bool("p", false, "only list files that will be updated")
	templateArg := flag.String("t", "", "a copyright `template` file (default: .copyright.txt)")
	verboseFlag := flag.Bool("v", false, "set verbose logging")
	versionFlag := flag.Bool("version", false, "display the version and exit")

	flag.CommandLine.Parse(cmdLine)

	if *helpFlag {
		Usage()
		return 0, false // stop processing, but no error
	}

	if *versionFlag {
		program := filepath.Base(os.Args[0])
		fmt.Printf("%s %s\n", program, version)
		return 0, false // stop processing, but no error
	}

	isPreview = *previewFlag
	isVerbose = *verboseFlag
	destDir = *destArg
	templateFile = *templateArg
	populateExclusions(*excludedList)

	args := flag.Args()
	// at least 1 argument expected
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "expected at least 1 source directory")
		Usage()
		return 1, false // stop processing and exit with error value 1
	}

	sources = args[0:]
	return 0, true // all's good
}

// Parses the excluded list argument and populates the excludePatterns slice.
func populateExclusions(excludedList string) {
	if len(excludedList) > 0 {
		for path := range strings.SplitSeq(excludedList, ",") {
			AddExclusion(path)
		}
	}
}

/*
 * displays memory usage statistics.
 */
func displayMemory() {
	if isVerbose {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc = %v kB\n", m.Alloc/1024)
		fmt.Printf("TotalAlloc = %v kB\n", m.TotalAlloc/1024)
		fmt.Printf("System = %v kB\n", m.Sys/1024)
	}
}
