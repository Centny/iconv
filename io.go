package iconv

import "io"

//Reader is the iconv reader which supported convert the data to speical encoding.
type Reader struct {
	Base io.Reader
	left int
	rlen int
	buf  []byte
	conv *ICONV
}

//NewReader is the default creator for Reader
func NewReader(base io.Reader, from, to string, bsize int) (*Reader, error) {
	conv, err := NewICONV(from, to)
	if err != nil {
		return nil, err
	}
	return &Reader{
		Base: base,
		buf:  make([]byte, bsize),
		conv: conv,
	}, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	r.rlen, err = r.Base.Read(r.buf[r.left:])
	if err != nil {
		return
	}
	tlen := r.rlen + r.left
	r.left, n, err = r.conv.Convert(r.buf[0:tlen], p)
	if err == nil && r.left > 0 {
		copy(r.buf, r.buf[tlen-r.left:tlen])
	}
	return
}

//Close will free converter
func (r *Reader) Close() error {
	return r.conv.Release()
}

//Writer is the iconv writer which supported convert the data to speical encoding.
type Writer struct {
	Base io.Writer
	left int
	ibuf []byte
	obuf []byte
	conv *ICONV
}

//NewWriter is the default creator for Writer
func NewWriter(base io.Writer, from, to string, bsize int) (*Writer, error) {
	conv, err := NewICONV(from, to)
	if err != nil {
		return nil, err
	}
	return &Writer{
		Base: base,
		ibuf: make([]byte, bsize),
		obuf: make([]byte, bsize),
		conv: conv,
	}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var poffset = 0
	var pleft = len(p)
	var clen, tlen int
	for {
		if pleft > 0 {
			clen = copy(w.ibuf[w.left:], p[poffset:])
			pleft -= clen
			poffset += clen
		} else {
			clen = 0
		}
		tlen = w.left + clen
		w.left, n, err = w.conv.Convert(w.ibuf[0:tlen], w.obuf)
		if err != nil || n < 1 {
			break
		}
		_, err = w.Base.Write(w.obuf[0:n])
		if err != nil {
			break
		}
		if w.left > 0 {
			copy(w.ibuf, w.ibuf[tlen-w.left:tlen])
		}
	}
	return len(p) - w.left, err
}

//Close will free converter
func (w *Writer) Close() error {
	return w.conv.Release()
}
