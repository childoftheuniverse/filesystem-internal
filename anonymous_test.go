package internal

import (
	"bytes"
	"golang.org/x/net/context"
	"testing"
)

const (
	TESTDATA = "This is some test data!"
)

/*
Test that any data written into an anonymous file object is the same when it
is read back.
*/
func TestReadWrite(t *testing.T) {
	var ctx = context.Background()
	var file = NewAnonymousFile()
	var data []byte
	var n int
	var err error

	n, err = file.Write(ctx, []byte(TESTDATA))
	if err != nil {
		t.Errorf("Write returned an unexpected error: %v", err)
	}
	if n != len(TESTDATA) {
		t.Errorf("Unexpected write length: %d (expected: %d)",
			n, len(TESTDATA))
	}

	// Use Close() to reset to the beginning of the file.
	file.Close(ctx)

	data = make([]byte, len(TESTDATA))
	n, err = file.Read(ctx, data)

	if err != nil {
		t.Errorf("Read returned an unexpected error: %v", err)
	}
	if n != len(TESTDATA) {
		t.Errorf("Unexpected length read back: %d (expected: %d)", len(data),
			len(TESTDATA))
	}

	if !bytes.Equal(data, []byte(TESTDATA)) {
		t.Errorf("Unexpected data read: expected \"%v\", got \"%v\"",
			data, []byte(TESTDATA))
	}
}
