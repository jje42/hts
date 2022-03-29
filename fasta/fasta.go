package fasta

import (
	"bufio"
	"io"
)

type Record struct {
	Header   string
	Sequence string
}

type Scanner struct {
	b          *bufio.Reader
	rec        *Record
	nextHeader string
	err        error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		b: bufio.NewReader(r),
	}
}

func (s *Scanner) Scan() bool {
	for {
		line, err := s.b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return s.rec != nil
			}
			s.err = err
			return false
		}
		if line[0] == '>' {
			if s.rec == nil {
				s.rec = &Record{Header: line[1:]}
			} else {
				s.nextHeader = line[1:]
				return true
			}
		}
		s.rec.Sequence += line
	}
}

func (s *Scanner) Record() *Record {
	r := s.rec
	if s.nextHeader != "" {
		s.rec = &Record{Header: s.nextHeader}
		s.nextHeader = ""
	} else {
		s.rec = nil
	}
	return r
}

func (s *Scanner) Err() error {
	return s.err
}
