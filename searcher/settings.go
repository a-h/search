package searcher

// Settings for the searcher.
type Settings struct {
	IncludeNames       []string
	ExcludeNames       []string
	IncludeDirectories bool
	IncludeFiles       bool
	IncludeText        string
	ExcludeText        string
}

// IsContentSearch returns true if the content of files must be evaluated.
func (s Settings) IsContentSearch() bool {
	return len(s.IncludeText) > 0 || len(s.ExcludeText) > 0
}
