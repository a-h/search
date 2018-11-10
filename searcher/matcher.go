package searcher

import (
	"context"
	"fmt"
	"path"
	"strings"
)

type matcher struct {
	Settings    Settings
	bannedPaths []string
}

// IsMatch determines whether the file / path matches the settings.
func (m *matcher) isMatch(ctx context.Context, path string, isDir bool) (matched bool, bytesRead int64, err error) {
	for _, bp := range m.bannedPaths {
		if strings.HasPrefix(path, bp) {
			return
		}
	}
	if len(m.Settings.ExcludeNames) > 0 {
		if ok := matches(m.Settings.ExcludeNames, path); ok {
			if isDir {
				m.bannedPaths = append(m.bannedPaths, path)
			}
			return
		}
	}
	if isDir && !m.Settings.IncludeDirectories {
		return
	}
	if !isDir && !m.Settings.IncludeFiles {
		return
	}
	if len(m.Settings.IncludeNames) > 0 {
		if ok := matches(m.Settings.IncludeNames, path); !ok {
			return
		}
	}
	if len(m.Settings.IncludeText) > 0 && isDir {
		return
	}
	if len(m.Settings.IncludeText) > 0 && !isDir {
		ok, r, tErr := m.Settings.TextSearch(ctx, path, m.Settings.IncludeText)
		if tErr != nil {
			err = fmt.Errorf("%v: %v", path, tErr)
			return
		}
		bytesRead = r
		if !ok {
			return
		}
	}
	matched = true
	return
}

func matches(patterns []string, name string) bool {
	for _, pattern := range patterns {
		_, suffix := path.Split(name)
		if m, _ := path.Match(pattern, suffix); m {
			return true
		}
	}
	return false
}
