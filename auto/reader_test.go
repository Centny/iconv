package auto

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	bys1, err := ReadFileAsUtf8("../gbk.txt")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(bys1))
	defer os.Remove("utf8_tmp.txt")
	err = ioutil.WriteFile("utf8_tmp.txt", bys1, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
	bys2, err := ReadFileAsUtf8("utf8_tmp.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys1) != string(bys2) {
		t.Error("error")
		return
	}
	bys3, err := ReadFileAsUtf8("../gbk2.txt")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(bys3))
	//
	//test unknow
	// os.Remove("unknow.txt")
	// ioutil.WriteFile("unknow.txt", []byte{0x0}, os.ModePerm)
	// _, err = ReadFileAsUtf8("unknow.txt")
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }

	//
	//test error
	_, err = ReadFileAsUtf8("xkkdss")
	if err == nil {
		t.Error(err)
		return
	}
	_, err = ReadFileAs("../gbk.txt", "xxdss")
	if err == nil {
		t.Error(err)
		return
	}
}
