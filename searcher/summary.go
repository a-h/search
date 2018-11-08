package searcher

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

// Summary of activity.
type Summary struct {
	Directories int
	Files       int
	BytesRead   int64
	TimeTaken   time.Duration
}

func (s Summary) String() string {
	if s.BytesRead > 0 {
		return fmt.Sprintf("Visited %d directories, %d files and read %v in %v",
			s.Directories,
			s.Files,
			bytefmt.ByteSize(uint64(s.BytesRead)),
			s.TimeTaken)
	}
	return fmt.Sprintf("Visited %d directories and %d files in %v",
		s.Directories,
		s.Files,
		s.TimeTaken)
}
