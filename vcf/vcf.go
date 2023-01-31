package vcf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func findBcftools() (string, error) {
	exe, err := exec.LookPath("bcftools")
	if err != nil {
		// PBS Pro doesn't set the PATH to the same as a login shell.
		// Even if bcftools is in the same directory as tso and it's on
		// PATH, we may still not find it with LookPath.
		this, err := os.Executable()
		if err != nil {
			return "", errors.New("cannot get executable path")
		}
		exe = filepath.Join(filepath.Dir(this), "bcftools")
		if _, err := os.Stat(exe); errors.Is(err, os.ErrNotExist) {
			return "", errors.New("unable to find bcftools binary")
		}
	}
	return exe, nil
}

func NewScanner(v VCF, loc ...string) (*Scanner, error) {
	var err error
	s := &Scanner{vcf: v}
	exe, err := findBcftools()
	if err != nil {
		return nil, err
	}
	s.cmd = exec.Command(exe, "view", "-H", v.file)
	s.stdout, err = s.cmd.StdoutPipe()
	if err != nil {
		return s, err
	}
	return s, nil
}

func NewScannerFromCommand(cmd *exec.Cmd, loc ...string) (*Scanner, error) {
	var err error
	s := &Scanner{}
	s.cmd = cmd
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
		buf := make([]byte, 0, 100000)
		s.scanner.Buffer(buf, 100000)
		s.scanCalled = true
	}
	for s.scanner.Scan() {
		token, err := parseVcfLine(s.scanner.Text(), s.vcf.Header.Samples)
		if err != nil {
			s.err = err
			return false
		}
		token.header = &s.vcf.Header
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

func CreateIndex(f string) error {
	if strings.HasSuffix(f, ".vcf") {
		return errors.New("cannot index uncompressed VCF file")
	}
	exe, err := findBcftools()
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(f, ".vcf.gz") {
		cmd = exec.Command(exe, "index", "-t", f)
	} else {
		cmd = exec.Command(exe, "index", f)
	}
	return cmd.Run()
}
