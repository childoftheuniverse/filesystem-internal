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
Tell returns our current offset in the virtual file.
*/
func (f *AnonymousFile) Tell(ctx context.Context) (int64, error) {
	return int64(f.position), nil
}

/*
Seek changes the current position to the specified one.
*/
func (f *AnonymousFile) Seek(ctx context.Context, offset int64) error {
	if offset < 0 || offset > int64(len(f.contents)) {
		return io.EOF
	}

	f.position = int(offset)
	return nil
}

/*
Skip advances our position cursor by the specified amount of bytes.
*/
func (f *AnonymousFile) Skip(ctx context.Context, diff int64) error {
	if diff < 0 || int64(f.position)+diff > int64(len(f.contents)) {
		return io.EOF
	}

	f.position += int(diff)
	return nil
}

/*
Close resets the position. That's about all it does. This is so the same
object can be reused in tests.
*/
func (f *AnonymousFile) Close(ctx context.Context) error {
	f.position = 0
	return nil
}
