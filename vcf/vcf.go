package vcf

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type VCF struct {
	file   string
	Header Header
}

func New(file string) (VCF, error) {
	header, err := readHeaderFromFile(file)
	if err != nil {
		return VCF{}, fmt.Errorf("unable to create VCF: %w", err)
	}
	return VCF{
		file:   file,
		Header: header,
	}, nil
}

type Scanner struct {
	vcf        VCF
	cmd        *exec.Cmd
	stdout     io.ReadCloser
	token      Variant
	err        error
	scanner    *bufio.Scanner
	scanCalled bool
	done       bool
}

func NewScanner(v VCF, loc ...string) (*Scanner, error) {
	var err error
	s := &Scanner{vcf: v}
	s.cmd = exec.Command("bcftools", "view", "-H", v.file)
	s.stdout, err = s.cmd.StdoutPipe()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Scanner) Scan() bool {
	if s.done {
		return false
	}
	if !s.scanCalled {
		if err := s.cmd.Start(); err != nil {
			s.err = err
			return false
		}
		s.scanner = bufio.NewScanner(s.stdout)
		s.scanCalled = true
	}
	for s.scanner.Scan() {
		token, err := parseVcfLine(s.scanner.Text(), s.vcf.Header.Samples)
		if err != nil {
			s.err = err
			return false
		}
		s.token = token
		return true
	}
	return false
}

func (s *Scanner) Variant() Variant {
	return s.token
}

func (s *Scanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}
