# search

Search your file system without having to remember how to string together complex `find` and `grep` commands with `xargs`.

```
$ search .
Directory: /Users/adrian/go/src/github.com/a-h/search
Include Names: []
Exclude Names: [.git]
Include Directories: true
Include Files: true
Include Text:
Exclude Text:
/Users/adrian/go/src/github.com/a-h/search
/Users/adrian/go/src/github.com/a-h/search/README.md
/Users/adrian/go/src/github.com/a-h/search/contains
/Users/adrian/go/src/github.com/a-h/search/contains/contains.go
/Users/adrian/go/src/github.com/a-h/search/main.go
/Users/adrian/go/src/github.com/a-h/search/searcher
/Users/adrian/go/src/github.com/a-h/search/searcher/matcher.go
/Users/adrian/go/src/github.com/a-h/search/searcher/searcher.go
/Users/adrian/go/src/github.com/a-h/search/searcher/searcher_test.go
/Users/adrian/go/src/github.com/a-h/search/searcher/settings.go
/Users/adrian/go/src/github.com/a-h/search/searcher/summary.go
Visited 32 directories and 46 files in 2.390197ms
```

```
$ search -text github.com .
Directory: /Users/adrian/go/src/github.com/a-h/search
Include Names: []
Exclude Names: [.git]
Include Directories: true
Include Files: true
Include Text: github.com
Exclude Text:
/Users/adrian/go/src/github.com/a-h/search/README.md
/Users/adrian/go/src/github.com/a-h/search/main.go
/Users/adrian/go/src/github.com/a-h/search/searcher/searcher_test.go
Visited 32 directories, 46 files and read 10.3K in 2.455713ms
```

## Install

```
go get github.com/a-h/search
cd search
go install
```

## Get help

```bash
search
```

```
usage: search [<args>] directory
  -directories
        set to true to include directories (default true)
  -exclude-names string
        path names to exclude (default ".git")
  -exclude-text string
        text to search for
  -files
        set to true to include files (default true)
  -names string
        path names to search for
  -print-settings
        prints out the settings during search (default true)
  -print-summary
        prints out a summary after search (default true)
  -text string
        text to search for
```

## List all files and directories in the current subdirectory and lower

```
search .
```

## Find text in all .go files

```
search -names *.go -text github.com/a-h .
```

## Find directories called `search`

```
search -names search -files=false .
```

## Find directories that start with `s`

```
search -names "s*" -files=false .
```