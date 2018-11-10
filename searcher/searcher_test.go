package searcher

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSearcher(t *testing.T) {
	file := fileInfo{
		fIsDir: false,
	}
	dir := fileInfo{
		fIsDir: true,
	}
	walker := func(root string, walkFn filepath.WalkFunc) error {
		walkFn("/root/file1.txt", file, nil)
		walkFn("/root/sub", dir, nil)
		walkFn("/root/sub/file2.txt", file, nil)
		return nil
	}
	matchAll := func(ctx context.Context, path string, isDir bool) (matched bool, err error) {
		return true, nil
	}
	matchNone := func(ctx context.Context, path string, isDir bool) (matched bool, err error) {
		return false, nil
	}
	var errMatching = errors.New("error matching")
	matchError := func(ctx context.Context, path string, isDir bool) (matched bool, err error) {
		return false, errMatching
	}
	matchAllContent := func(ctx context.Context, path string, text string) (matched bool, bytesRead int64, err error) {
		return true, 10, nil
	}
	tests := []struct {
		name           string
		walker         func(root string, walkFn filepath.WalkFunc) error
		settings       Settings
		pathMatcher    func(ctx context.Context, path string, isDir bool) (matched bool, err error)
		contentMatcher func(ctx context.Context, path string, text string) (matched bool, bytesRead int64, err error)
		expected       []string
		expectedErrors []error
	}{
		{
			name:        "if we match everything, we get all the paths",
			walker:      walker,
			pathMatcher: matchAll,
			expected: []string{
				"/root/file1.txt",
				"/root/sub",
				"/root/sub/file2.txt",
			},
		},
		{
			name:   "a text match can be successful",
			walker: walker,
			settings: Settings{
				IncludeText: "test",
			},
			pathMatcher:    matchAll,
			contentMatcher: matchAllContent,
			expected: []string{
				"/root/file1.txt",
				"/root/sub/file2.txt",
			},
		},
		{
			name:        "if we match nothing, we get nothing",
			walker:      walker,
			pathMatcher: matchNone,
			expected:    nil,
		},
		{
			name:           "if the matcher errors, we receive it",
			walker:         walker,
			pathMatcher:    matchError,
			expected:       nil,
			expectedErrors: []error{errMatching, errMatching, errMatching},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := Searcher{
				PathMatcher:    tc.pathMatcher,
				ContentMatcher: tc.contentMatcher,
				Walker:         tc.walker,
				Settings:       tc.settings,
			}
			var wg sync.WaitGroup
			pc := make(chan string)
			errc := make(chan error)
			wg.Add(1)
			go func() {
				var err error
				s.Walk(context.Background(), "/", pc, errc)
				if err != nil && err != ErrCancelled {
					t.Errorf("failed to walk: %v", err)
				}
				wg.Done()
			}()

			var paths []string
			wg.Add(1)
			go func() {
				for p := range pc {
					paths = append(paths, p)
				}
				wg.Done()
			}()
			var errors []error
			wg.Add(1)
			go func() {
				for err := range errc {
					errors = append(errors, err)
				}
				wg.Done()
			}()
			wg.Wait()
			if !reflect.DeepEqual(tc.expected, paths) {
				t.Errorf("expected paths: %v, got %v", tc.expected, paths)
			}
			if !reflect.DeepEqual(tc.expectedErrors, errors) {
				t.Errorf("expected errors: %v, got %v", tc.expected, errors)
			}
		})
	}
}

type fileInfo struct {
	fName  string
	fSize  int64
	fMode  os.FileMode
	fIsDir bool
}

func (fi fileInfo) Name() string {
	return fi.fName
}

func (fi fileInfo) Size() int64 {
	return fi.fSize
}

func (fi fileInfo) Mode() os.FileMode {
	return fi.fMode
}

func (fi fileInfo) ModTime() time.Time {
	return time.Time{}
}

func (fi fileInfo) IsDir() bool {
	return fi.fIsDir
}

func (fi fileInfo) Sys() interface{} {
	return true
}

func TestMatcherPaths(t *testing.T) {
	type previousInput struct {
		path  string
		isDir bool
	}

	tests := []struct {
		name           string
		settings       Settings
		previousInputs []previousInput
		inputPath      string
		isDir          bool
		expected       bool
		expectedErr    error
	}{
		{
			name:      "files are not matched by default",
			settings:  Settings{},
			inputPath: "/code/hello.txt",
			expected:  false,
		},
		{
			name:      "directories are not matched by default",
			settings:  Settings{},
			inputPath: "/code/hello",
			isDir:     true,
			expected:  false,
		},
		{
			name: "any file name is allowed if no names are included",
			settings: Settings{
				IncludeFiles: true,
			},
			inputPath: "/code/hello.txt",
			isDir:     false,
			expected:  true,
		},
		{
			name: "file names can be wildcard included positively",
			settings: Settings{
				IncludeFiles: true,
				IncludeNames: []string{"*.txt"},
			},
			inputPath: "/code/hello.txt",
			isDir:     false,
			expected:  true,
		},
		{
			name: "file names can be wildcard included negatively",
			settings: Settings{
				IncludeFiles: true,
				IncludeNames: []string{"*.go"},
			},
			inputPath: "/code/hello.txt",
			isDir:     false,
			expected:  false,
		},
		{
			name: "file names can be wildcard excluded positively",
			settings: Settings{
				IncludeFiles: true,
				ExcludeNames: []string{"*.txt"},
			},
			inputPath: "/code/hello.txt",
			isDir:     false,
			expected:  false,
		},
		{
			name: "file names can be wildcard excluded negatively",
			settings: Settings{
				IncludeFiles: true,
				ExcludeNames: []string{"*.go"},
			},
			inputPath: "/code/hello.txt",
			isDir:     false,
			expected:  true,
		},
		{
			name: "directories can't match text",
			settings: Settings{
				IncludeFiles:       true,
				IncludeDirectories: true,
				IncludeText:        "test",
			},
			inputPath: "/code",
			isDir:     true,
			expected:  false,
		},
		{
			name: "find a directory by name",
			settings: Settings{
				IncludeFiles:       false,
				IncludeDirectories: true,
				IncludeNames:       []string{"search"},
				ExcludeNames:       []string{".git"},
			},
			inputPath: "/Users/adrian/go/src/github.com/a-h/search",
			isDir:     true,
			expected:  true,
		},
		{
			name: "directories can be ignored",
			settings: Settings{
				IncludeFiles:       false,
				IncludeDirectories: true,
				ExcludeNames:       []string{".git"},
			},
			inputPath: "/Users/adrian/go/src/github.com/a-h/search/.git",
			isDir:     true,
			expected:  false,
		},
		{
			name: "subdirectories can be ignored",
			settings: Settings{
				IncludeFiles:       false,
				IncludeDirectories: true,
				ExcludeNames:       []string{".git"},
			},
			previousInputs: []previousInput{
				previousInput{
					path:  "/Users/adrian/go/src/github.com/a-h/search/.git",
					isDir: true},
			},
			inputPath: "/Users/adrian/go/src/github.com/a-h/search/.git/test",
			isDir:     true,
			expected:  false,
		},
		{
			name: "subdirectories can be ignored by files",
			settings: Settings{
				IncludeFiles:       true,
				IncludeDirectories: false,
				ExcludeNames:       []string{".git"},
			},
			previousInputs: []previousInput{
				previousInput{
					path:  "/Users/adrian/go/src/github.com/a-h/search/.git",
					isDir: true},
			},
			inputPath: "/Users/adrian/go/src/github.com/a-h/search/.git/test",
			isDir:     false,
			expected:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := matcher{
				Settings: test.settings,
			}
			for _, pi := range test.previousInputs {
				m.isPathMatch(context.Background(), pi.path, pi.isDir)
			}
			actual, actualErr := m.isPathMatch(context.Background(), test.inputPath, test.isDir)
			if !reflect.DeepEqual(actualErr, test.expectedErr) {
				t.Fatalf("expected error '%v', got '%v'", test.expectedErr, actualErr)
			}
			if actual != test.expected {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}
