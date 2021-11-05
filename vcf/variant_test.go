package vcf

import (
	"testing"
)

func TestVariant_Type(t *testing.T) {
	type fields struct {
		Chrom     string
		Pos       int
		ID        string
		Ref       string
		Alt       []string
		Qual      string
		Filter    []string
		Info      map[string]string
		Format    []string
		genotypes []Genotype
		header    *Header
	}
	tests := []struct {
		name   string
		fields fields
		want   Type
	}{
		{"t1", fields{Ref: "A", Alt: []string{"C"}}, SNP},
		{"t2", fields{Ref: "ATGC", Alt: []string{"CGTA"}}, MNP},
		{"t3", fields{Ref: "A", Alt: []string{"ATGC"}}, INDEL},
		{"t4", fields{Ref: "ATGC", Alt: []string{"A"}}, INDEL},
		{"t5", fields{Ref: "A", Alt: []string{"C", "G"}}, SNP},
		{"t6", fields{Ref: "ATGC", Alt: []string{"AAAA", "TTTT"}}, MNP},
		{"t7", fields{Ref: "ATGC", Alt: []string{"AT", "AGTGTA"}}, INDEL},
		{"t8", fields{Ref: "A", Alt: []string{"C", "ATGC"}}, MIXED},
		{"t8", fields{Ref: "A", Alt: []string{}}, NO_VARIATION},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Variant{
				Chrom:     tt.fields.Chrom,
				Pos:       tt.fields.Pos,
				ID:        tt.fields.ID,
				Ref:       tt.fields.Ref,
				Alt:       tt.fields.Alt,
				Qual:      tt.fields.Qual,
				Filter:    tt.fields.Filter,
				Info:      tt.fields.Info,
				Format:    tt.fields.Format,
				genotypes: tt.fields.genotypes,
				header:    tt.fields.header,
			}
			if got := v.Type(); got != tt.want {
				t.Errorf("Variant.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariant_IsFiltered(t *testing.T) {
	type fields struct {
		Chrom     string
		Pos       int
		ID        string
		Ref       string
		Alt       []string
		Qual      string
		Filter    []string
		Info      map[string]string
		Format    []string
		genotypes []Genotype
		header    *Header
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{"t1", fields{Filter: []string{"PASS"}}, false},
		{"t2", fields{Filter: []string{"."}}, false},
		{"t3", fields{Filter: []string{"LowQual"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Variant{
				Chrom:     tt.fields.Chrom,
				Pos:       tt.fields.Pos,
				ID:        tt.fields.ID,
				Ref:       tt.fields.Ref,
				Alt:       tt.fields.Alt,
				Qual:      tt.fields.Qual,
				Filter:    tt.fields.Filter,
				Info:      tt.fields.Info,
				Format:    tt.fields.Format,
				genotypes: tt.fields.genotypes,
				header:    tt.fields.header,
			}
			if got := v.IsFiltered(); got != tt.want {
				t.Errorf("Variant.IsFiltered() = %v, want %v", got, tt.want)
			}
		})
	}
}
