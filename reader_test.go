package gpsdscanner

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type tNew struct {
	Init        func(td *tNew)
	Input       []byte
	BytesReader *bytes.Reader
	Reader      io.Reader
	WantScanner *Scanner
	WantErr     error
}

func errorTest(got error, want error, i int) (ok bool, s string) {
	var (
		wantStr string
		gotStr  string
	)
	if want == nil {
		wantStr = `<nil>`
	} else {
		wantStr = want.Error()
	}
	if got == nil {
		gotStr = `<nil>`
	} else {
		gotStr = got.Error()
	}
	if strings.Compare(wantStr, gotStr) != 0 {
		format := "%s%s\n, %s%s\"\n%s %d]"
		s1 := `Got the error "`
		s2 := `"... but want the error "`
		s3 := `The test error ocurred in test: tdNew[`
		s = fmt.Sprintf(format, s1, gotStr, s2, wantStr, s3, i)
		ok = false
	} else {
		s = ""
		ok = true
	}
	return
}

func byteSliceTest(got []byte, want []byte, i int) (ok bool, s string) {
	var (
		wantStr string
		gotStr  string
	)
	if want == nil {
		wantStr = "<empty>"
	} else {
		wantStr = string(want)
		if len(wantStr) == 0 {
			wantStr = "<empty>"
		}
	}
	if got == nil {
		gotStr = "<empty>"
	} else {
		gotStr = string(got)
		if len(gotStr) == 0 {
			gotStr = "<empty>"
		}
	}
	if strings.Compare(wantStr, gotStr) != 0 {
		format := "%s%s\n, %s%s\"\n%s %d]"
		s1 := `Got the byte slice "`
		s2 := `"... but want the byte slice "`
		s3 := `The test error ocurred in test: tdNew[`
		s = fmt.Sprintf(format, s1, gotStr, s2, wantStr, s3, i)
		ok = false
	} else {
		s = ""
		ok = true
	}
	return
}

var tdNew []*tNew = []*tNew{
	//Test 0:
	&tNew{
		BytesReader: nil,
		Input:       nil,
		Init:        nil,
		Reader:      nil,
		WantScanner: nil,
		WantErr: fmt.Errorf("%s%s",
			`Error while creating a  "github.com/larsth/linescanner".`,
			`Linescanner : Nil io.Reader`),
	},
	//Test 1:
	&tNew{
		Input: []byte(``),
		//Init function initializes td.BytesReader, and td.Reader
		Init: func(td *tNew) {
			td.BytesReader = bytes.NewReader(td.Input)
			td.Reader = io.Reader(td.BytesReader)
		},
		WantScanner: &Scanner{},
		WantErr:     nil,
	},
}

func TestNew(t *testing.T) {
	var (
		gotScanner *Scanner
		gotErr     error
	)
	for i, td := range tdNew {
		if td.Init != nil {
			td.Init(td)
		}
		gotScanner, gotErr = New(td.Reader)

		//*Scanner
		if (gotScanner == nil) && td.WantScanner != nil {
			format := "%s\nThe test error ocurred in tdNew[%d]\nWant scanner is: \n\t%#v"
			s := "Got a *Scanner which is nil, but a want *Scanner, which is not nil"
			t.Errorf(format, s, i, td.WantScanner)
		}
		if (gotScanner != nil) && td.WantScanner == nil {
			format := "%s\nThe test error ocurred in tdNew[%d]\nGot scanner is: \n\t%#v"
			s := "Got a *Scanner which is not nil, but a want *Scanner, which is nil"
			t.Errorf(format, s, i, gotScanner)
		}
		//error
		if ok, s := errorTest(gotErr, td.WantErr, i); false == ok {
			t.Error(s)
		}
	}
}

type tScan struct {
	Init        func(td *tScan) error
	Input       []byte
	BytesReader *bytes.Reader
	WantP       []byte
	Scanner     *Scanner
	WantErr     error
}

var tdScan = []*tScan{
	//test 0:
	&tScan{
		Init: func(td *tScan) error {
			td.Scanner = &Scanner{}
			return nil
		},
		Input: nil,
		WantP: nil,
		WantErr: fmt.Errorf("%s",
			`The scanner is not initialized (nil io.Reader)`),
	},
	//test 0:
	&tScan{
		Init: func(td *tScan) error {
			var err error
			td.BytesReader = bytes.NewReader(td.Input)
			if td.Scanner, err = New(io.Reader(td.BytesReader)); err != nil {
				return err
			}
			return nil
		},
		Input:   []byte(``),
		WantP:   nil,
		WantErr: nil,
	},
}

func TestScan(t *testing.T) {
	var (
		gotErr error
		gotP   []byte
		i      int
		td     *tScan
	)
	for i, td = range tdScan {
		if td.Init != nil {
			if err := td.Init(td); err != nil {
				t.Fatal("Fatal error: Cannot initialize with New(io.Reader): ",
					"error from New function: ", err.Error())
			}
		}
		_, gotP, gotErr = td.Scanner.Scan()
		//p
		if ok, s := byteSliceTest(gotP, td.WantP, i); false == ok {
			t.Error(s)
		}
		//error
		if ok, s := errorTest(gotErr, td.WantErr, i); false == ok {
			t.Error(s)
		}
	}
}
