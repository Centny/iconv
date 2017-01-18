package iconv

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iconv.h>
#include <errno.h>
#cgo darwin CPPFLAGS: -I/usr/local/include
#cgo darwin LDFLAGS: -L/usr/local/lib -liconv
#cgo linux CPPFLAGS: -I/usr/local/include
#cgo linux LDFLAGS: -L/usr/local/lib -liconv


iconv_t iconv_open_v(const char *from, const char *to)
{
    iconv_t res = iconv_open(to, from);
    if (res == (iconv_t)-1)
    {
        return 0;
    }
    return res;
}


size_t iconv_v(iconv_t cd, char *in, size_t *inleft, char *out, size_t *outleft)
{
    return iconv(cd, &in, inleft, &out, outleft);
}

*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

//E2BIG is the iconv error for need more output buffer
var E2BIG = syscall.Errno(C.E2BIG)

//EILSEQ is the iconv error for invalid input
var EILSEQ = syscall.Errno(C.EILSEQ)

//EINVAL is the iconv error for need more input buffer
var EINVAL = syscall.Errno(C.EINVAL)

//ICONV is the converter struct.
type ICONV struct {
	cv C.iconv_t
}

//NewICONV create on converter by from encoding and target encoding
func NewICONV(from, to string) (*ICONV, error) {
	fromcode := C.CString(from)
	defer C.free(unsafe.Pointer(fromcode))
	tocode := C.CString(to)
	defer C.free(unsafe.Pointer(tocode))
	conv := C.iconv_open_v(fromcode, tocode)
	if conv == nil {
		return nil, fmt.Errorf("open fail with code(%v)", -1)
	}
	return &ICONV{
		cv: conv,
	}, nil
}

//Convert will convert the in data to out buffer and return the left bytes and avail bytes
func (i *ICONV) Convert(in, out []byte) (left, avail int, err error) {
	var indata = (*C.char)(C.CBytes(in))
	defer C.free(unsafe.Pointer(indata))
	var outdata = (*C.char)(C.CBytes(out))
	defer C.free(unsafe.Pointer(outdata))
	var inleft = C.size_t(len(in))
	var outleft = C.size_t(len(out))
	_, err = C.iconv_v(i.cv, indata, &inleft, outdata, &outleft)
	if err == nil || err == E2BIG || err == EINVAL {
		// fmt.Println("xx->", len(in), inleft, outleft, avail, err)
		err = nil
		avail = len(out) - int(outleft)
		left = int(inleft)
		var gdata = C.GoBytes(unsafe.Pointer(outdata), C.int(avail))
		copy(out, gdata)
	}
	return
}

//ConvertAll will convert all input to one output
func (i *ICONV) ConvertAll(in []byte) (out []byte, err error) {
	out = make([]byte, len(in))
	var left, avail, total = len(in), 0, 0
	var tlen = 0
	for x := 2; ; x++ {
		tlen = left
		left, avail, err = i.Convert(in[0:left], out[total:])
		total += avail
		if err != nil || avail < 1 || left < 1 {
			break
		}
		copy(in, in[tlen-left:tlen])
		nout := make([]byte, x*len(in))
		copy(nout, out)
		out = nout
	}
	out = out[0:total]
	return
}

//Release will free the converter
func (i *ICONV) Release() (err error) {
	if i.cv == nil {
		return
	}
	var code = C.iconv_close(i.cv)
	if code != 0 {
		err = fmt.Errorf("return code(%v)", code)
	}
	i.cv = nil
	return
}
