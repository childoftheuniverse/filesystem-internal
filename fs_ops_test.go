package internal

import (
	"bytes"
	"github.com/childoftheuniverse/filesystem"
	"golang.org/x/net/context"
	"net/url"
	"os"
	"sort"
	"testing"
)

func TestWriteFileAndReadBack(t *testing.T) {
	var ctx = context.Background()
	var output filesystem.WriteCloser
	var input filesystem.ReadCloser
	var path *url.URL
	var contents = make([]byte, 3*len(TESTDATA))
	var n int
	var err error

	path, err = url.Parse("internal:///test/file/one")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	if output, err = filesystem.OpenWriter(ctx, path); err != nil {
		t.Fatalf("Error opening output file %s: %s", path, err)
	}

	if n, err = output.Write(ctx, []byte(TESTDATA)); err != nil {
		t.Errorf("Error writing testdata to %s: %s", path, err)
	}
	if n != len(TESTDATA) {
		t.Errorf("Mismatched length writing testdata to %s: %d (expected %d)",
			path, n, len(TESTDATA))
	}

	if n, err = output.Write(ctx, []byte(TESTDATA)); err != nil {
		t.Errorf("Error writing testdata to %s: %s", path, err)
	}
	if n != len(TESTDATA) {
		t.Errorf("Mismatched length writing testdata to %s: %d (expected %d)",
			path, n, len(TESTDATA))
	}

	if err = output.Close(ctx); err != nil {
		t.Errorf("Error closing output file %s: %s", path, err)
	}

	if input, err = filesystem.OpenReader(ctx, path); err != nil {
		t.Fatalf("Error opening input file %s: %s", path, err)
	}

	if n, err = input.Read(ctx, contents); err != nil {
		t.Errorf("Error reading %d bytes from %s: %s",
			len(contents), path, err)
	}
	if n != 2*len(TESTDATA) {
		t.Errorf("Mismatched length reading testdata %s: %d (expected %d)",
			path, n, 2*len(TESTDATA))
	}

	if err = input.Close(ctx); err != nil {
		t.Errorf("Error closing input file %s: %s", path, err)
	}

	if !bytes.Equal(contents[:2*len(TESTDATA)],
		append([]byte(TESTDATA), []byte(TESTDATA)...)) {
		t.Errorf("Mismatched contents: was \"%s\", expected \"%s\"",
			string(contents),
			string(append([]byte(TESTDATA), []byte(TESTDATA)...)))
	}

	if err = filesystem.Remove(ctx, path); err != nil {
		t.Errorf("Error removing %s: %s", path, err)
	}

	_, err = filesystem.OpenReader(ctx, path)
	if err == nil {
		t.Errorf("No error opening %s after deleting", path)
	} else if err != os.ErrNotExist {
		t.Errorf("Unexpected type of error opening %s after deleting: %s",
			path, err)
	}
}

func TestWriteFileAndReadBackRepeatedly(t *testing.T) {
	var ctx = context.Background()
	var output filesystem.WriteCloser
	var input filesystem.ReadCloser
	var path *url.URL
	var contents = make([]byte, 3*len(TESTDATA))
	var n int
	var i int
	var err error

	path, err = url.Parse("internal:///test/file/two")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	for i = 0; i < 5; i++ {
		if output, err = filesystem.OpenWriter(ctx, path); err != nil {
			t.Fatalf("Error opening output file %s: %s", path, err)
		}

		if n, err = output.Write(ctx, []byte(TESTDATA)); err != nil {
			t.Errorf("Error writing testdata to %s: %s", path, err)
		}
		if n != len(TESTDATA) {
			t.Errorf("Mismatched length writing testdata to %s: %d (expected %d)",
				path, n, len(TESTDATA))
		}

		if n, err = output.Write(ctx, []byte(TESTDATA)); err != nil {
			t.Errorf("Error writing testdata to %s: %s", path, err)
		}
		if n != len(TESTDATA) {
			t.Errorf("Mismatched length writing testdata to %s: %d (expected %d)",
				path, n, len(TESTDATA))
		}

		if err = output.Close(ctx); err != nil {
			t.Errorf("Error closing output file %s: %s", path, err)
		}
	}

	for i = 0; i < 5; i++ {
		if input, err = filesystem.OpenReader(ctx, path); err != nil {
			t.Fatalf("Error opening input file %s: %s", path, err)
		}

		if n, err = input.Read(ctx, contents); err != nil {
			t.Errorf("Error reading %d bytes from %s: %s",
				len(contents), path, err)
		}
		if n != 2*len(TESTDATA) {
			t.Errorf("Mismatched length reading testdata %s: %d (expected %d)",
				path, n, 2*len(TESTDATA))
		}

		if err = input.Close(ctx); err != nil {
			t.Errorf("Error closing input file %s: %s", path, err)
		}

		if !bytes.Equal(contents[:n],
			append([]byte(TESTDATA), []byte(TESTDATA)...)) {
			t.Errorf("Mismatched contents: was \"%s\", expected \"%s\"",
				string(contents),
				string(append([]byte(TESTDATA), []byte(TESTDATA)...)))
		}
	}

	if err = filesystem.Remove(ctx, path); err != nil {
		t.Errorf("Error removing %s: %s", path, err)
	}
}

func TestAppendFileAndReadBack(t *testing.T) {
	var ctx = context.Background()
	var output filesystem.WriteCloser
	var input filesystem.ReadCloser
	var path *url.URL
	var contents = make([]byte, 3*len(TESTDATA))
	var n int
	var i int
	var err error

	path, err = url.Parse("internal:///test/file/three")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	for i = 0; i < 2; i++ {
		if output, err = filesystem.OpenAppender(ctx, path); err != nil {
			t.Fatalf("Error opening output file %s: %s", path, err)
		}

		if n, err = output.Write(ctx, []byte(TESTDATA)); err != nil {
			t.Errorf("Error writing testdata to %s: %s", path, err)
		}
		if n != len(TESTDATA) {
			t.Errorf("Mismatched length writing testdata to %s: %d (expected %d)",
				path, n, len(TESTDATA))
		}

		if err = output.Close(ctx); err != nil {
			t.Errorf("Error closing output file %s: %s", path, err)
		}
	}

	if input, err = filesystem.OpenReader(ctx, path); err != nil {
		t.Fatalf("Error opening input file %s: %s", path, err)
	}

	if n, err = input.Read(ctx, contents); err != nil {
		t.Errorf("Error reading %d bytes from %s: %s",
			len(contents), path, err)
	}
	if n != 2*len(TESTDATA) {
		t.Errorf("Mismatched length reading testdata %s: %d (expected %d)",
			path, n, 2*len(TESTDATA))
	}

	if err = input.Close(ctx); err != nil {
		t.Errorf("Error closing input file %s: %s", path, err)
	}

	if !bytes.Equal(contents[:n],
		append([]byte(TESTDATA), []byte(TESTDATA)...)) {
		t.Errorf("Mismatched contents: was \"%s\", expected \"%s\"",
			string(contents),
			string(append([]byte(TESTDATA), []byte(TESTDATA)...)))
	}

	if err = filesystem.Remove(ctx, path); err != nil {
		t.Errorf("Error removing %s: %s", path, err)
	}
}

func TestListDirectoryContents(t *testing.T) {
	var ctx = context.Background()
	var output filesystem.WriteCloser
	var path *url.URL
	var files []string
	var i int
	var err error

	path, err = url.Parse("internal:///test/dir/one1")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	if output, err = filesystem.OpenAppender(ctx, path); err != nil {
		t.Fatalf("Error opening output file %s: %s", path, err)
	}

	if err = output.Close(ctx); err != nil {
		t.Errorf("Error closing output file %s: %s", path, err)
	}

	path, err = url.Parse("internal:///test/dir/two2")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	if output, err = filesystem.OpenAppender(ctx, path); err != nil {
		t.Fatalf("Error opening output file %s: %s", path, err)
	}

	if err = output.Close(ctx); err != nil {
		t.Errorf("Error closing output file %s: %s", path, err)
	}

	path, err = url.Parse("internal:///test/dir/three3")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	if output, err = filesystem.OpenAppender(ctx, path); err != nil {
		t.Fatalf("Error opening output file %s: %s", path, err)
	}

	if err = output.Close(ctx); err != nil {
		t.Errorf("Error closing output file %s: %s", path, err)
	}

	path, err = url.Parse("internal:///test/dir/subdir/four4")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	if output, err = filesystem.OpenAppender(ctx, path); err != nil {
		t.Fatalf("Error opening output file %s: %s", path, err)
	}

	if err = output.Close(ctx); err != nil {
		t.Errorf("Error closing output file %s: %s", path, err)
	}

	path, err = url.Parse("internal:///test/dir")
	if err != nil {
		t.Fatalf("Error parsing test path: %s", err)
	}

	files, err = filesystem.ListEntries(ctx, path)
	if err != nil {
		t.Errorf("Error listing contents of %s: %s", path, err)
	}

	if len(files) != 4 {
		t.Errorf("Wrong number of files in result: %v", files)
	}

	sort.Strings(files)

	if i = sort.SearchStrings(files, "one1"); files[i] != "one1" {
		t.Errorf("Could not find one1 in %v", files)
	}
	if i = sort.SearchStrings(files, "two2"); files[i] != "two2" {
		t.Errorf("Could not find two2 in %v", files)
	}
	if i = sort.SearchStrings(files, "three3"); files[i] != "three3" {
		t.Errorf("Could not find three3 in %v", files)
	}
	if i = sort.SearchStrings(files, "subdir"); files[i] != "subdir" {
		t.Errorf("Could not find subdir in %v", files)
	}
}
