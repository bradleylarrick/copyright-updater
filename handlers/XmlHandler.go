package handlers

import "os"

type XmlHandler struct{}

var (
	xmlHeader = "<!--"
	xmlFooter = " -->"
	xmlPrefix = " "
)

func (XmlHandler) Format(src *os.File, dest *os.File, copyright *[]string) error {
	err := startProcess(src, dest, xmlHeader, xmlFooter, xmlPrefix)
	if err != nil {
		return err
	}
	writeCopyright(copyright)
	return endProcess()
}
