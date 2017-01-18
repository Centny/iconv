package iconv

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	var utf8dat = []byte(`
链接：http://pan.baidu.com/s/1eR2hRYa 密码：hy07


此链接失效，加Q群有快速找资源教程。


最新影片交流Q群：585038005
    `)
	fmt.Println("utf8->", len(utf8dat))
	enc, err := NewICONV("utf8", "gbk")
	if err != nil {
		t.Error(err)
		return
	}
	gbkdat, err := enc.ConvertAll([]byte(utf8dat))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("gbk->", len(gbkdat))
	fmt.Println("\n\n\n\ngbk->utf8", len(gbkdat))
	enc2, err := NewICONV("gbk", "UTF8")
	if err != nil {
		t.Error(err)
		return
	}
	utf8dat2, err := enc2.ConvertAll([]byte(gbkdat))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("utf8->", len(utf8dat2))
	if string(utf8dat) != string(utf8dat2) {
		t.Error("error")
		return
	}
}
