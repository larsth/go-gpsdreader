//Package gpsdscanner scans from a stream of bytes and returns a byte 
//slice, which contains one gpsd JSON document
package gpsdscanner

import (
	"io"
	"sync"

	"github.com/larsth/linescanner"
)

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
		return nil, err
	}
	return r, nil
}

func (s *Scanner) Scan() ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.lineScanner.Scan() == false {
		if s.lineScanner.Err() != nil {
			return nil, s.lineScanner.Err()
		}
	}
	return s.lineScanner.Bytes(), nil
}
