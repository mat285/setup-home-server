package task

import (
	"io"
	"strings"
	"sync"
)

type CaptureWriter struct {
	lock    sync.Mutex
	out     io.Writer
	buf     []byte
	prepend []byte
}

func NewCaptureWriter(out io.Writer, prepend []byte) *CaptureWriter {
	return &CaptureWriter{
		out:     out,
		buf:     make([]byte, 0),
		prepend: prepend,
	}
}

func (cw *CaptureWriter) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	cw.lock.Lock()
	defer cw.lock.Unlock()
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		bline := []byte(line + "\n")
		cw.out.Write(append(cw.prepend, bline...))
		cw.buf = append(cw.buf, bline...)
	}
	return len(buf), nil
	idx := containsNewline(buf)
	total := 0
	for idx >= 0 {
		n, err := cw.out.Write(buf[:idx])
		total += n
		if n > 0 {
			cw.buf = append(cw.buf, buf[:n]...)
		}
		if err != nil {
			return total, err
		}
		cw.out.Write(cw.prepend)
		buf = buf[idx+1:]
		idx = containsNewline(buf)
	}
	n, err := cw.out.Write(buf)
	total += n
	if n > 0 {
		cw.buf = append(cw.buf, buf[:n]...)
	}
	return total, err
}

func (cw *CaptureWriter) Captured() []byte {
	cw.lock.Lock()
	defer cw.lock.Unlock()
	return cw.buf
}

func containsNewline(buf []byte) int {
	for i, b := range buf {
		if b == '\n' {
			return i
		}
	}
	return -1
}
