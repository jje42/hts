package seq

import (
	"fmt"
)

// revComp returns the reverse complement of the DNA string.
func ReverseComplement(s string) string {
	var r string
	for _, x := range s {
		switch x {
		case 'A':
			r = "T" + r
		case 'T':
			r = "A" + r
		case 'G':
			r = "C" + r
		case 'C':
			r = "G" + r
		default:
			panic("unexpected nucleotide")
		}
	}
	return r
}

// type DNASequence string

func Translate(s string) (string, error) {
	if len(s)%3 != 0 {
		return "", fmt.Errorf("length not a multiple of three (is %d)", len(s))
	}
	p := make([]byte, len(s)/3)
	j := 0
	for i := 0; i < len(s)-2; i += 3 {
		codon := s[i : i+3]
		aa := codon2aa(codon)
		p[j] = aa
		j++
	}
	return string(p), nil
}
func codon2aa(codon string) byte {
	switch codon {
	case "GCA", "GCC", "GCG", "GCT":
		return 'A'
	case "AGA", "AGG", "CGA", "CGC", "CGG", "CGT":
		return 'R'
	case "AAC", "AAT":
		return 'N'
	case "GAC", "GAT":
		return 'D'
	case "TGC", "TGT":
		return 'C'
	case "CAA", "CAG":
		return 'Q'
	case "GAA", "GAG":
		return 'E'
	case "GGA", "GGC", "GGG", "GGT":
		return 'G'
	case "CAC", "CAT":
		return 'H'
	case "ATA", "ATC", "ATT":
		return 'I'
	case "CTA", "CTC", "CTG", "CTT", "TTA", "TTG":
		return 'L'
	case "AAA", "AAG":
		return 'K'
	case "ATG":
		return 'M'
	case "TTC", "TTT":
		return 'F'
	case "CCA", "CCC", "CCG", "CCT":
		return 'P'
	case "AGC", "AGT", "TCA", "TCC", "TCG", "TCT":
		return 'S'
	case "ACA", "ACC", "ACG", "ACT":
		return 'T'
	case "TGG":
		return 'W'
	case "TAC", "TAT":
		return 'Y'
	case "GTA", "GTC", "GTG", "GTT":
		return 'V'
	case "TAA", "TAG", "TGA":
		return '*'
	default:
		panic(fmt.Sprintf("not a codon: %s", codon))
	}
}
