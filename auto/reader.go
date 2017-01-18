package auto

import (
	"io/ioutil"

	"fmt"

	"github.com/Centny/iconv"
	"github.com/Centny/uchardet"
)

//ReadFileAs read file as special encoding
func ReadFileAs(filename, encoding string, try ...string) (res []byte, err error) {
	var bys []byte
	bys, err = ioutil.ReadFile(filename)
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
	var trylist []string
	if len(coding) > 0 {
		trylist = []string{coding}
	} else {
		trylist = try
	}
	var enc *iconv.ICONV
	for _, code := range trylist {
		enc, err = iconv.NewICONV(code, encoding)
		if err != nil {
			continue
		}
		res, err = enc.ConvertAll(bys)
		enc.Release()
		if err == nil {
			return
		}
	}
	return nil, fmt.Errorf("try convert from %v to %v fail", trylist, encoding)
}

//ReadFileAsUtf8 read file and convert data to utf8
func ReadFileAsUtf8(filename string) ([]byte, error) {
	return ReadFileAs(filename, "UTF-8", "GBK", "UTF-8")
}
