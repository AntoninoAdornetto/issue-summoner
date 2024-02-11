package utils

import "os"

type OSFileOpener struct{}

func (o OSFileOpener) Open(fileName string) (*os.File, error) {
	return os.Open(fileName)
}
