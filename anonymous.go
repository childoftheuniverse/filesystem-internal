package internal

import (
	"golang.org/x/net/context"
	"io"
)

/*
AnonymousFile represents an inmemory buffer which is not bound to any path.
*/
type AnonymousFile struct {
	contents []byte
	position int
}

/*
NewAnonymousFile simply creates a new empty anonymous file.
*/
func NewAnonymousFile() *AnonymousFile {
	return &AnonymousFile{contents: make([]byte, 0)}
}

/*
Read reads up to len(data) bytes from the file and puts them into data.
*/
func (f *AnonymousFile) Read(ctx context.Context, data []byte) (int, error) {
	var rdata []byte
	var newEnd int
	var length int

	if f.position >= len(f.contents) {
		return 0, io.EOF
	}

	if f.position+len(data) >= len(f.contents) {
		newEnd = len(f.contents)
	} else {
		newEnd = f.position + len(data)
	}

	length = newEnd - f.position

	rdata = f.contents[f.position:newEnd]
	copy(data, rdata)

	f.position = newEnd

	return length, nil
}

/*
Write appends the specified data to the end of the buffer. This is not always
correct, but it's easy enough to make sense for most purposes.
*/
func (f *AnonymousFile) Write(ctx context.Context, data []byte) (int, error) {
	f.contents = append(f.contents, data...)
	return len(data), nil
}

/*
Len returns the current total length of the anonymous file.
This is not part of the standard file API, but it is certainly useful for
testing.
*/
func (f *AnonymousFile) Len() int {
	return len(f.contents)
}

/*
Close resets the position. That's about all it does. This is so the same
object can be reused in tests.
*/
func (f *AnonymousFile) Close(ctx context.Context) error {
	f.position = 0
	return nil
}
