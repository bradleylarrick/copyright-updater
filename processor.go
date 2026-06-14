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
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	handlers "natuna.org/copyright/handlers"
)

type FileHandler interface {
	Format(src *os.File, dest *os.File) error
	AddProtected([]string)
}

type Processor struct {
	Handlers map[string]FileHandler
}

// Returns a Processor with the handler map populated.
func NewProcessor() *Processor {
	return &Processor{
		Handlers: map[string]FileHandler{
			".apt":        handlers.AptHandler{},
			".bash":       handlers.HashtagHandler{},
			".bat":        handlers.BatHandler{},
			".cs":         handlers.JavaHandler{},
			".css":        handlers.JavaHandler{},
			".csv":        handlers.HashtagHandler{},
			".go":         handlers.JavaHandler{},
			".gradle":     handlers.JavaHandler{},
			".groovy":     handlers.JavaHandler{},
			".html":       handlers.XmlHandler{},
			".java":       handlers.JavaHandler{},
			".js":         handlers.JavaHandler{},
			".md":         handlers.XmlHandler{},
			".mk":         handlers.HashtagHandler{},
			".properties": handlers.HashtagHandler{},
			".py":         handlers.HashtagHandler{},
			".rs":         handlers.JavaHandler{},
			".sh":         handlers.HashtagHandler{},
			".toml":       handlers.HashtagHandler{},
			".txt":        handlers.HashtagHandler{},
			".vm":         handlers.VmHandler{},
			".xaml":       handlers.XmlHandler{},
			".xmi":        handlers.XmlHandler{},
			".xml":        handlers.XmlHandler{},
			".xsd":        handlers.XmlHandler{},
			".yaml":       handlers.HashtagHandler{},
		},
	}
}

/*
 * Processes the given file.
 *
 * Returns an error if the file could not be processed.
 */
func (p Processor) ProcessFile(path string, name string) error {

	var destPath string
	if len(destDir) > 0 {
		destPath = filepath.Join(destDir, path)
	} else {
		destPath = path
	}

	err := validateDestPath(destPath)
	if err != nil {
		return err
	}

	fullSrc := filepath.Join(path, name)
	fullDest := filepath.Join(destPath, name)
	exclude := IsExcluded(fullSrc)
	if exclude {
		if isVerbose {
			fmt.Printf("Skipping excluded file: %s\n", fullSrc)
		}
		return nil
	}

	// Find the extension, which determines the updater to use
	ext := findExtension(name, fullSrc)

	handler, ok := p.Handlers[ext]
	if !ok {
		if isVerbose {
			fmt.Printf("Skipping %s\n", fullSrc)
		}
		return nil
	} else if isPreview {
		fmt.Println(fullSrc)
		return nil
	} else {
		if isVerbose {
			fmt.Printf("Processing %s -> %s\n", fullSrc, fullDest)
		}

		srcFile, err := os.Open(fullSrc)
		if err != nil {
			return err
		}

		tempFile, err := os.CreateTemp(destPath, "copyright-*.tmp")
		if err != nil {
			return err
		}

		err = handler.Format(srcFile, tempFile)
		srcFile.Close()
		tempFile.Close()
		if err == nil {
			err = os.Rename(tempFile.Name(), fullDest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to rename temp file: %v\n", err)
				return err
			}
		}
		return err
	}

}

/**
 * Finds the extension (which determines the updater to use) of the file based on its name.
 */
func findExtension(name string, fullSrc string) string {
	ext := filepath.Ext(name)
	// If the file has no extension and is a script file, use ".sh" as the extension.
	if len(ext) == 0 && isScriptFile(fullSrc) {
		return ".sh"
	}

	// If the file is a Makefile or Jenkinsfile, use the appropriate extension.
	if strings.HasPrefix(strings.ToLower(name), "makefile") {
		ext = ".mk"
	} else if strings.HasPrefix(strings.ToLower(name), "jenkinsfile") {
		ext = ".java"
	} else {
		pieces := strings.Split(name, ".")
		ndx := len(pieces) - 1
		if len(pieces) > 2 && strings.EqualFold(pieces[ndx], "vm") && !strings.EqualFold(pieces[ndx-1], "txt") {
			ndx--
		}
		ext = "." + pieces[ndx]
	}

	return ext
}

/*
 * Validates the destination path, creating it if it does not exist.
 */
func validateDestPath(destPath string) error {
	if _, err := os.Stat(destPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

/*
 * Checks if the file is a script file based on the first line starting with "#!".
 */
func isScriptFile(fullSrc string) bool {
	srcFile, err := os.Open(fullSrc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(4)
	}
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#!") {
			return true
		}
	}

	return false
}
