package internal

import (
	"github.com/childoftheuniverse/filesystem"
	"golang.org/x/net/context"
	"net/url"
	"os"
	"strings"
)

func init() {
	var ifs = &internalFileSystem{
		knownFiles: make(map[string]*AnonymousFile),
	}

	filesystem.AddImplementation("internal", ifs)
}

/*
internalFileSystem allows to perform filesystem-like operations on inmemory
files.
*/
type internalFileSystem struct {
	knownFiles map[string]*AnonymousFile
}

/*
internalOpenFile determines whether a file with the known path exists already,
creates a new file with the name otherwise.
*/
func (i *internalFileSystem) internalOpenFile(ctx context.Context, path *url.URL, create bool) (*AnonymousFile, error) {
	var file *AnonymousFile
	var ok bool

	file, ok = i.knownFiles[path.Path]
	if !ok && !create {
		return nil, os.ErrNotExist
	}
	if !ok {
		file = NewAnonymousFile()
		i.knownFiles[path.Path] = file
	}

	return file, nil
}

func (i *internalFileSystem) OpenReader(ctx context.Context, path *url.URL) (filesystem.ReadCloser, error) {
	return i.internalOpenFile(ctx, path, false)
}

func (i *internalFileSystem) OpenWriter(ctx context.Context, path *url.URL) (filesystem.WriteCloser, error) {
	i.knownFiles[path.Path] = NewAnonymousFile()
	return i.knownFiles[path.Path], nil
}

func (i *internalFileSystem) OpenAppender(ctx context.Context, path *url.URL) (filesystem.WriteCloser, error) {
	return i.internalOpenFile(ctx, path, true)
}

func (i *internalFileSystem) ListEntries(ctx context.Context, path *url.URL) ([]string, error) {
	var matchingPaths = make(map[string]bool)
	var returnPaths []string
	var knownPath string

	for knownPath, _ = range i.knownFiles {
		if strings.HasPrefix(knownPath, path.Path) {
			var subpath = knownPath[len(path.Path):]

			for len(subpath) > 0 && subpath[0] == '/' {
				var i int
				subpath = subpath[1:]

				for i = 0; i < len(subpath); i++ {
					if subpath[i] == '/' {
						subpath = subpath[:i]
						break
					}
				}
			}

			if len(subpath) > 0 {
				matchingPaths[subpath] = true
			}
		}
	}

	for knownPath, _ = range matchingPaths {
		returnPaths = append(returnPaths, knownPath)
	}

	return returnPaths, nil
}

func (i *internalFileSystem) WatchFile(context.Context, *url.URL, filesystem.FileWatchFunc) (filesystem.CancelWatchFunc, chan error, error) {
	return nil, nil, filesystem.EUNSUPP
}

func (i *internalFileSystem) Remove(ctx context.Context, path *url.URL) error {
	delete(i.knownFiles, path.Path)
	return nil
}
