iconv binding by go
======

### Install

<br/>
`linux/unix/window/osx`
* install `libiconv` by `yum`,`brew` or compile by source

* install iconv by command

```
export CGO_CPPFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -liconv"
go get github.com/Centny/iconv
```

* test iconv

```
go test -v github.com/Centny/iconv
```

note: adding libsigar install path to LD_LIBRARY_PATH


### Example

`iconv`

```
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
```

`read all as utf8 (github.com/Centny/uchardet needed)`

```
	bys1, err := ReadFileAsUtf8("../gbk.txt")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(bys1))
```