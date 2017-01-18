package auto

import (
	"io/ioutil"

	"fmt"

	"github.com/Centny/iconv"
	"github.com/Centny/uchardet"
)

//ReadFileAs read file as special encoding
func ReadFileAs(filename, encoding string) ([]byte, error) {
	bys, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	det := uchardet.NewChardet()
	defer det.Release()
	var code = det.Handle(bys)
	if code != 0 {
		return nil, fmt.Errorf("chardet handle return code(%v)", code)
	}
	coding := det.End()
	if len(coding) < 1 {
		return nil, fmt.Errorf("uchardet det encoding fail")
	}
	if coding == encoding {
		return bys, nil
	}
	enc, err := iconv.NewICONV(coding, encoding)
	if err != nil {
		return nil, fmt.Errorf("create iconv by from(%v),to(%v) fail with %v", coding, encoding, err)
	}
	defer enc.Release()
	return enc.ConvertAll(bys)
}

//ReadFileAsUtf8 read file and convert data to utf8
func ReadFileAsUtf8(filename string) ([]byte, error) {
	return ReadFileAs(filename, "UTF-8")
}
