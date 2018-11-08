package searcher

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// New creates a new Searcher with the Settings.
func New(s Settings) Searcher {
	m := matcher{
		Settings: s,
	}
	return Searcher{
		Walker:  filepath.Walk,
		Matcher: m.isMatch,
	}
}

// ErrCancelled is returned when the search is cancelled by the context timing out or being cancelled.
var ErrCancelled = errors.New("search cancelled early")

// Searcher searches through files.
type Searcher struct {
	Walker  func(root string, walkFn filepath.WalkFunc) error
	Matcher func(path string, isDir bool) (matched bool, bytesRead int64, err error)
}

// Walk the directory specified in the settings.
func (s Searcher) Walk(ctx context.Context, directory string, paths chan string, errors chan error) (summary Summary, err error) {
	start := time.Now()
	searcher := func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return ErrCancelled
		default:
			// Continue.
		}
		isDir := info.IsDir()
		if isDir {
			summary.Directories++
		} else {
			summary.Files++
		}
		matched, read, err := s.Matcher(path, isDir)
		summary.BytesRead += read
		if err != nil {
			errors <- err
		}
		if matched {
			paths <- path
		}
		return nil
	}
	err = s.Walker(directory, searcher)
	summary.TimeTaken = time.Now().Sub(start)
	return
}
