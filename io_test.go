package iconv

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	file, err := os.Open("gbk.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	reader, err := NewReader(file, "gbk", "utf8", 10)
	if err != nil {
		t.Error(err)
		return
	}
	var rlen int
	var buf = make([]byte, 64)
	for {
		rlen, err = reader.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("%v", string(buf[0:rlen]))
	}
	if err != io.EOF {
		t.Error(err)
		return
	}
}

func TestWriter(t *testing.T) {
	file, err := os.Open("gbk.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	defer os.RemoveAll("gbk_tmp.txt")
	fileo, err := os.OpenFile("gbk_tmp.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
	defer fileo.Close()
	reader, err := NewReader(file, "gbk", "utf8", 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()
	writer, err := NewWriter(fileo, "utf8", "gbk", 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer writer.Close()
	utf8data := []byte{}
	var wlen, rlen int
	var buf = make([]byte, 64)
	for {
		rlen, err = reader.Read(buf)
		if err != nil {
			fmt.Println("reader->", err)
			break
		}
		utf8data = append(utf8data, buf[0:rlen]...)
		wlen, err = writer.Write(buf[0:rlen])
		if err != nil {
			fmt.Println("writer->", err)
			break
		}
		if rlen != wlen {
			t.Error("error")
			return
		}
		//fmt.Printf("%v", string(buf[0:rlen]))
	}
	if err != io.EOF {
		t.Error(err)
		return
	}
	file.Close()
	fileo.Close()
	//
	gbkbys, err := ioutil.ReadFile("gbk_tmp.txt")
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := NewICONV("gbk", "utf8")
	if err != nil {
		t.Error(err)
		return
	}
	utf8dat2, err := enc.ConvertAll(gbkbys)
	if err != nil {
		t.Error(err)
		return
	}
	if string(utf8data) != string(utf8dat2) {
		fmt.Println(len(utf8data), len(utf8dat2))
		fmt.Println(string(utf8dat2))
		t.Error("error")
		return
	}
	reader.Close()
	//
	//test error
	_, err = NewReader(file, "xkkd", "dsfsdf", 1000)
	if err == nil {
		t.Error("error")
		return
	}
	_, err = NewWriter(file, "xkkd", "dsfsdf", 1000)
	if err == nil {
		t.Error("error")
		return
	}
	ewriter, _ := NewWriter(&errw{}, "utf8", "gbk", 100)
	_, err = ewriter.Write(utf8data)
	if err == nil {
		t.Error("error")
		return
	}
}

type errw struct {
}

func (e *errw) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("%v", "error")
}
