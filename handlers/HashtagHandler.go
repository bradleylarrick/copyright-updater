package handlers

import "os"

type HashtagHandler struct{}

var (
	hashHeader = "#"
	hashFooter = "#"
	hashPrefix = "#"
)

func (HashtagHandler) Format(src *os.File, dest *os.File, copyright *[]string) error {
	err := startProcess(src, dest, hashHeader, hashFooter, hashPrefix)
	if err != nil {
		return err
	}

	findProtected([]string{"#!"})
	findHeader()
	writeCopyright(copyright)
	return endProcess()
}
