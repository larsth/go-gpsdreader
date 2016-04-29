//Package gpsdscanner scans from a stream of bytes and returns a byte
//slice, which contains one gpsd JSON document
package gpsdscanner

import (
	"io"
	"sync"

	"github.com/juju/errors"
	"github.com/larsth/linescanner"
)

var ErrNotInitialized = errors.New("The scanner is not initialized (nil io.Reader)")

type Scanner struct {
	mutex       sync.Mutex
	reader      io.Reader
	lineScanner *linescanner.LineScanner
}

func New(reader io.Reader) (*Scanner, error) {
	var err error

	r := new(Scanner)
	r.reader = reader
	if r.lineScanner, err = linescanner.New(reader); err != nil {
		annotatedErr := errors.Annotatef(err, "%s %s",
			"Error while creating a ",
			"\"github.com/larsth/linescanner\".Linescanner ")
		return nil, annotatedErr
	}
	return r, nil
}

func (s *Scanner) Scan() (ok bool, p []byte, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.reader == nil {
		return false, nil, ErrNotInitialized
	}

	ok = s.lineScanner.Scan()
	if err := s.lineScanner.Err(); err != nil {
		annotatedErr := errors.Annotate(err,
			"A \"github.com/larsth/linescanner\".Linescanner.Scan() error")
		return false, nil, annotatedErr
	}
	return true, s.lineScanner.Bytes(), nil
}
