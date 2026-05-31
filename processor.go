package main

import (
	"fmt"
	"os"
	"path/filepath"

	handlers "natuna.org/copyright/handlers"
)

type FileHandler interface {
	Format(src *os.File, dest *os.File, copyright *[]string) error
}

type Processor struct {
	handlers map[string]FileHandler
}

// Returns a Processor with the handler map populated.
func NewProcessor() *Processor {
	return &Processor{
		handlers: map[string]FileHandler{
			".bash":   handlers.HashtagHandler{},
			".cs":     handlers.JavaHandler{},
			".css":    handlers.JavaHandler{},
			".csv":    handlers.HashtagHandler{},
			".go":     handlers.JavaHandler{},
			".groovy": handlers.JavaHandler{},
			".html":   handlers.XmlHandler{},
			".java":   handlers.JavaHandler{},
			".js":     handlers.JavaHandler{},
			".py":     handlers.HashtagHandler{},
			".sh":     handlers.HashtagHandler{},
			".toml":   handlers.HashtagHandler{},
			".xmi":    handlers.XmlHandler{},
			".xml":    handlers.XmlHandler{},
			".xsd":    handlers.XmlHandler{},
		},
	}
}

// Processes the given file
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
	ext := filepath.Ext(name)
	exclude := IsExcluded(fullSrc)
	if exclude {
		if isVerbose {
			fmt.Printf("Skipping excluded file: %s\n", fullSrc)
		}
		return nil
	}

	// fmt.Printf("srcDir = %s, fullSrc = %s\n", srcDir, fullSrc)
	// fmt.Printf("base = %s\n", filepath.Base(fullSrc))
	handler, ok := p.handlers[ext]
	if !ok {
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

		err = handler.Format(srcFile, tempFile, &copyright.copyright)
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

/*
 * Validates the destination path, creating it if it does not exist.
 */
func validateDestPath(destPath string) error {
	if _, err := os.Stat(destPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
