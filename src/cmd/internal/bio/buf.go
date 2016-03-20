// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bio implements seekable buffered I/O.
package bio

import (
	"bufio"
	"io"
	"log"
	"os"
)

const EOF = -1

// Reader implements a seekable buffered io.Reader.
type Reader struct {
	f *os.File
	r *bufio.Reader
}

// Writer implements a seekable buffered io.Writer.
type Writer struct {
	f *os.File
	w *bufio.Writer
}

// Reader returns this Reader's underlying bufio.Reader.
func (r *Reader) Reader() *bufio.Reader { return r.r }

// Writer returns this Writer's underlying bufio.Writer.
func (w *Writer) Writer() *bufio.Writer { return w.w }

// Create creates the file named name and returns a Writer
// for that file.
func Create(name string) (*Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &Writer{f: f, w: bufio.NewWriter(f)}, nil
}

// Open returns a Reader for the file named name.
func Open(name string) (*Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Reader{f: f, r: bufio.NewReader(f)}, nil
}

// BufWriter returns a Writer on top of w.
// TODO(dfc) remove this method and replace caller with bufio.Writer.
func BufWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// BufWriter returns a Reader on top of r.
// TODO(dfc) remove this method and replace caller with bufio.Reader.
func BufReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *Writer) WriteString(p string) (int, error) {
	return w.w.WriteString(p)
}

func (r *Reader) Seek(offset int64, whence int) int64 {
	if whence == 1 {
		offset -= int64(r.r.Buffered())
	}
	off, err := r.f.Seek(offset, whence)
	if err != nil {
		log.Fatalf("seeking in output: %v", err)
	}
	r.r.Reset(r.f)
	return off
}

func (w *Writer) Seek(offset int64, whence int) int64 {
	if err := w.w.Flush(); err != nil {
		log.Fatalf("writing output: %v", err)
	}
	off, err := w.f.Seek(offset, whence)
	if err != nil {
		log.Fatalf("seeking in output: %v", err)
	}
	return off
}

func (r *Reader) Offset() int64 {
	off, err := r.f.Seek(0, 1)
	if err != nil {
		log.Fatalf("seeking in output [0, 1]: %v", err)
	}
	off -= int64(r.r.Buffered())
	return off
}

func (w *Writer) Offset() int64 {
	if err := w.w.Flush(); err != nil {
		log.Fatalf("writing output: %v", err)
	}
	off, err := w.f.Seek(0, 1)
	if err != nil {
		log.Fatalf("seeking in output [0, 1]: %v", err)
	}
	return off
}

func (w *Writer) Flush() error {
	return w.w.Flush()
}

func (w *Writer) WriteByte(c byte) error {
	return w.w.WriteByte(c)
}

func Bread(r *Reader, p []byte) int {
	n, err := io.ReadFull(r.r, p)
	if n == 0 {
		if err != nil && err != io.EOF {
			n = -1
		}
	}
	return n
}

func Bgetc(r *Reader) int {
	c, err := r.r.ReadByte()
	if err != nil {
		if err != io.EOF {
			log.Fatalf("reading input: %v", err)
		}
		return EOF
	}
	return int(c)
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *Reader) Peek(n int) ([]byte, error) {
	return r.r.Peek(n)
}

func Brdline(r *Reader, delim int) string {
	s, err := r.r.ReadBytes(byte(delim))
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}
	return string(s)
}

func (r *Reader) Close() error {
	return r.f.Close()
}

func (w *Writer) Close() error {
	err := w.w.Flush()
	err1 := w.f.Close()
	if err == nil {
		err = err1
	}
	return err
}