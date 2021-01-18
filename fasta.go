package hts

import (
	"fmt"
	"os"

	"github.com/biogo/hts/fai"
)

// Fasta is a FASTA file.
type Fasta struct {
	r  *os.File
	ri *os.File
	fa *fai.File
}

// NewFasta creates a new Fasta type with index.
func NewFasta(path string) (Fasta, error) {
	var ret Fasta
	fi, err := os.Open(fmt.Sprintf("%s.fai", path))
	if err != nil {
		return ret, err
	}
	idx, err := fai.ReadFrom(fi)
	if err != nil {
		return ret, err
	}
	f, err := os.Open(path)
	if err != nil {
		return ret, err
	}
	fa := fai.NewFile(f, idx)
	return Fasta{r: f, ri: fi, fa: fa}, nil
}

// Query returns the sequence found in contig between start and end.
func (f Fasta) Query(contig string, start, end int) (string, error) {
	seq, err := f.fa.SeqRange(contig, start, end)
	if err != nil {
		return "", err
	}
	b := make([]byte, end-start)
	_, err = seq.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
