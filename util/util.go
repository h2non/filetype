package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
)

// ReadMSOfficeXMLFile Read the contents of the xml file
// contained in the ms office file according to the specified file name.
func ReadMSOfficeXMLFile(content []byte, filename string) ([]byte, error) {
	br := bytes.NewReader(content)
	zipr, err := zip.NewReader(br, int64(len(content)))
	if err != nil {
		return nil, err
	}

	var file *zip.File
	for _, f := range zipr.File {
		if f.FileInfo().Name() == filename {
			file = f
			break
		}
	}

	if file == nil {
		return nil, fmt.Errorf("The specified file could not be found: %s", filename)
	}

	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
