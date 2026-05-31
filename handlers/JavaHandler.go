package handlers

import "os"

type JavaHandler struct{}

var (
	javaHeader = "/*"
	javaFooter = " */"
	javaPrefix = " *"
)

func (JavaHandler) Format(src *os.File, dest *os.File, copyright *[]string) error {
	err := startProcess(src, dest, javaHeader, javaFooter, javaPrefix)
	if err != nil {
		return err
	}

	findHeader()
	writeCopyright(copyright)
	return endProcess()
}
