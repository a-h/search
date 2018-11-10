package searcher

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

// Summary of activity.
type Summary struct {
	StartTime   time.Time
	EndTime     time.Time
	Directories int
	Files       int
	BytesRead   int64
}

// TimeTaken by the operation.
func (s Summary) TimeTaken() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

func (s Summary) String() string {
	if s.BytesRead > 0 {
		return fmt.Sprintf("Visited %d directories, %d files and read %v in %v",
			s.Directories,
			s.Files,
			bytefmt.ByteSize(uint64(s.BytesRead)),
			s.TimeTaken())
	}
	return fmt.Sprintf("Visited %d directories and %d files in %v",
		s.Directories,
		s.Files,
		s.TimeTaken())
}
