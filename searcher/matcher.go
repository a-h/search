package searcher

import (
	"context"
	"path"
	"strings"
)

type matcher struct {
	Settings    Settings
	bannedPaths []string
}

// isPathMatch determines whether the file / path matches the settings.
func (m *matcher) isPathMatch(ctx context.Context, path string, isDir bool) (matched bool, err error) {
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
	if m.Settings.IsContentSearch() && isDir {
		return
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
