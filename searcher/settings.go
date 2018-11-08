package searcher

// Settings for the searcher.
type Settings struct {
	IncludeNames       []string
	ExcludeNames       []string
	IncludeDirectories bool
	IncludeFiles       bool
	IncludeText        string
	ExcludeText        string
	TextSearch         func(name, text string) (ok bool, bytesRead int64, err error)
}
