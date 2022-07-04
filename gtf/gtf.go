package gtf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/brentp/faidx"
	"github.com/jje42/hts/seq"
)

type Feature struct {
	Seqname   string
	Source    string
	Feature   string
	Start     int
	End       int
	Score     float64 // can be '.' for "no score"
	Strand    string
	Frame     string
	Attribute map[string]string
}

type GTF struct {
	Header   []string
	Features []*Feature
	txs      map[string]*Feature
	exons    map[string][]*Feature
}

func New(r io.Reader) (*GTF, error) {
	g := &GTF{}
	g.txs = make(map[string]*Feature)
	g.exons = make(map[string][]*Feature)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			g.Header = append(g.Header, line)
		} else {
			bits := strings.Split(line, "\t")
			if len(bits) != 9 {
				return nil, errors.New("line contains unexpected number of fields")
			}
			start, err := strconv.Atoi(bits[3])
			if err != nil {
				return nil, fmt.Errorf("unable to convert start: %w", err)
			}
			end, err := strconv.Atoi(bits[4])
			if err != nil {
				return nil, fmt.Errorf("unable to convert end: %w", err)
			}
			score := -1.0
			if bits[5] != "." {
				score, err = strconv.ParseFloat(bits[5], 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse score: %w", err)
				}
			}
			m := make(map[string]string)
			attrs := strings.Split(strings.TrimSuffix(bits[8], ";"), ";")
			for _, at := range attrs {
				foo := strings.SplitN(strings.TrimSpace(at), " ", 2)
				key := foo[0]
				s, err := strconv.Unquote(foo[1])
				if err != nil {
					return nil, fmt.Errorf("failed to unquote attribute: %w", err)
				}
				// Ensembl's GTF can have multiple "tag" keys
				// in feature attributes.
				if _, ok := m[key]; ok {
					m[key] = m[key] + ";" + s
				} else {
					m[key] = s
				}
			}
			f := &Feature{
				Seqname:   bits[0],
				Source:    bits[1],
				Feature:   bits[2],
				Start:     start,
				End:       end,
				Score:     score,
				Strand:    bits[6],
				Frame:     bits[7],
				Attribute: m,
			}
			g.Features = append(g.Features, f)
			if f.Feature == "transcript" {
				txid, ok := f.Attribute["transcript_id"]
				if !ok {
					return nil, fmt.Errorf("transcript without transcript_id")
				}
				g.txs[txid] = f
			}
			if f.Feature == "exon" {
				txid, ok := f.Attribute["transcript_id"]
				if !ok {
					return nil, fmt.Errorf("exon without transcript_id")
				}
				g.exons[txid] = append(g.exons[txid], f)
			}
		}
		// if len(g.Features) >= 5 {
		// 	break
		// }
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return g, nil
}

type Transcript struct {
	Gene       string
	Biotype    string
	Tag        string
	Seqname    string
	Start, End int
	Strand     string
	Feat       *Feature
	Exons      []*Feature
}

func (t *Transcript) Sequence(fai *faidx.Faidx) (string, error) {
	var s string

	sort.SliceStable(t.Exons, func(i, j int) bool {
		if t.Exons[i].Start == t.Exons[j].Start {
			return t.Exons[i].End < t.Exons[j].End
		}
		return t.Exons[i].Start < t.Exons[j].Start
	})

	for _, f := range t.Exons {
		ss, err := fai.Get(f.Seqname, f.Start-1, f.End)
		if err != nil {
			return "", err
		}
		s += ss
	}
	if t.Strand == "-" {
		return seq.ReverseComplement(s), nil
	}
	return s, nil
}

func (g *GTF) Transcript(id string) (*Transcript, error) {
	t, ok := g.txs[id]
	if !ok {
		return nil, errors.New("no such transcript")
	}
	e, ok := g.exons[id]
	if !ok {
		return nil, errors.New("no such transcript")
	}
	return &Transcript{
		Seqname: t.Seqname,
		Start:   t.Start,
		End:     t.End,
		Strand:  t.Strand,
		Feat:    t,
		Exons:   e,
	}, nil
}
