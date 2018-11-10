package searcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/a-h/search/contains"
)

// New creates a new Searcher with the Settings.
func New(s Settings) Searcher {
	m := matcher{
		Settings: s,
	}
	return Searcher{
		Settings:       s,
		Walker:         filepath.Walk,
		PathMatcher:    m.isPathMatch,
		ContentMatcher: contains.TextInFile,
	}
}

// ErrCancelled is returned when the search is cancelled by the context timing out or being cancelled.
var ErrCancelled = errors.New("search cancelled early")

// Searcher searches through files.
type Searcher struct {
	Settings       Settings
	Walker         func(root string, walkFn filepath.WalkFunc) error
	PathMatcher    func(ctx context.Context, path string, isDir bool) (matched bool, err error)
	ContentMatcher func(ctx context.Context, path string, text string) (matched bool, bytesRead int64, err error)
}

// Walk the directory. Pass open channels for paths and errors.
// The paths and errors channels will be closed by this function.
func (s *Searcher) Walk(ctx context.Context, directory string, paths chan string, errors chan error) (summary Summary, err error) {
	summary = Summary{
		StartTime: time.Now(),
	}
	// If it's not a content search, we don't need to waste memory doing concurrent work.
	if !s.Settings.IsContentSearch() {
		summary.Files, summary.Directories, err = s.filterPaths(ctx, directory, paths, errors)
		close(paths)
		close(errors)
		return
	}
	// Connect non-concurrent path filtering to the concurrent content reading.
	var wg sync.WaitGroup

	// Start up several content filters.
	contentFilterInput := make(chan string)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			read := s.filterContent(ctx, contentFilterInput, paths, errors)
			summary.BytesRead += read
			wg.Done()
		}()
	}

	walkErrors := make(chan error)
	wg.Add(1)
	go func() {
		for wErr := range walkErrors {
			errors <- wErr
		}
		close(errors)
		wg.Done()
	}()

	// Start passing paths to the content filters.
	wg.Add(1)
	go func() {
		summary.Files, summary.Directories, err = s.filterPaths(ctx, directory, contentFilterInput, walkErrors)
		close(contentFilterInput)
		close(walkErrors)
		wg.Done()
	}()

	wg.Wait()
	close(paths)
	summary.EndTime = time.Now()
	return
}

func (s *Searcher) filterPaths(ctx context.Context, directory string,
	paths chan string, errors chan error) (files, directories int, err error) {
	searcher := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot search: %v", err)
		}
		select {
		case <-ctx.Done():
			return ErrCancelled
		default:
			// Continue.
		}
		isDir := info.IsDir()
		if isDir {
			directories++
		} else {
			files++
		}
		pathMatched, err := s.PathMatcher(ctx, path, isDir)
		if err != nil {
			errors <- err
			return nil
		}
		if !pathMatched {
			return nil
		}
		if s.Settings.IsContentSearch() && isDir {
			return nil
		}
		paths <- path
		return nil
	}
	err = s.Walker(directory, searcher)
	return
}

func (s *Searcher) filterContent(ctx context.Context, inputPaths chan string, matchedPaths chan string, errors chan error) (bytesRead int64) {
	for path := range inputPaths {
		matched, read, err := s.ContentMatcher(ctx, path, s.Settings.IncludeText)
		if err != nil {
			errors <- err
			continue
		}
		if matched {
			matchedPaths <- path
		}
		bytesRead += read
	}
	return
}
