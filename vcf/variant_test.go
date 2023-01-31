package vcf

import (
	"reflect"
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
		{"t1", fields{Filter: []string{}}, false},
		{"t2", fields{Filter: []string{"PASS"}}, false},
		{"t3", fields{Filter: []string{"."}}, false},
		{"t4", fields{Filter: []string{"LowQual"}}, true},
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

func Test_parseCsq(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			"t1",
			args{"Consequence annotations from Ensembl VEP. Format: Allele|Consequence|IMPACT|SYMBOL|Gene|Feature_type|Feature|BIOTYPE|EXON|INTRON|HGVSc|HGVSp|cDNA_position|CDS_position|Protein_position|Amino_acids|Codons|Existing_variation|DISTANCE|STRAND|FLAGS|PICK|VARIANT_CLASS|SYMBOL_SOURCE|HGNC_ID|CANONICAL|REFSEQ_MATCH|REFSEQ_OFFSET|GIVEN_REF|USED_REF|BAM_EDIT|HGVS_OFFSET|HGVSg"},
			[]string{
				"Allele", "Consequence", "IMPACT", "SYMBOL", "Gene", "Feature_type", "Feature",
				"BIOTYPE", "EXON", "INTRON", "HGVSc", "HGVSp", "cDNA_position",
				"CDS_position", "Protein_position", "Amino_acids", "Codons",
				"Existing_variation", "DISTANCE", "STRAND", "FLAGS", "PICK",
				"VARIANT_CLASS", "SYMBOL_SOURCE", "HGNC_ID", "CANONICAL",
				"REFSEQ_MATCH", "REFSEQ_OFFSET", "GIVEN_REF", "USED_REF",
				"BAM_EDIT", "HGVS_OFFSET", "HGVSg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCsq(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCsq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariant_CsqKeys(t *testing.T) {
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
	header := NewHeader()
	header.AddHeaderLines(
		NewComplexHeaderLine("INFO", map[string]string{
			"ID":          "CSQ",
			"Number":      ".",
			"Type":        "String",
			"Description": "Consequence annotations from Ensembl VEP. Format: Allele|Consequence|IMPACT|SYMBOL|Gene|Feature_type|Feature|BIOTYPE|EXON|INTRON|HGVSc|HGVSp|cDNA_position|CDS_position|Protein_position|Amino_acids|Codons|Existing_variation|DISTANCE|STRAND|FLAGS|PICK|VARIANT_CLASS|SYMBOL_SOURCE|HGNC_ID|CANONICAL|REFSEQ_MATCH|REFSEQ_OFFSET|GIVEN_REF|USED_REF|BAM_EDIT|HGVS_OFFSET|HGVSg",
		}),
	)
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"t1",
			fields{
				Chrom:     "1",
				Pos:       1000,
				ID:        "",
				Ref:       "A",
				Alt:       []string{"T"},
				Qual:      "99",
				Filter:    []string{},
				Info:      make(map[string]string),
				Format:    []string{},
				genotypes: []Genotype{},
				header:    &header,
			},
			[]string{
				"Allele", "Consequence", "IMPACT", "SYMBOL", "Gene", "Feature_type", "Feature",
				"BIOTYPE", "EXON", "INTRON", "HGVSc", "HGVSp", "cDNA_position",
				"CDS_position", "Protein_position", "Amino_acids", "Codons",
				"Existing_variation", "DISTANCE", "STRAND", "FLAGS", "PICK",
				"VARIANT_CLASS", "SYMBOL_SOURCE", "HGNC_ID", "CANONICAL",
				"REFSEQ_MATCH", "REFSEQ_OFFSET", "GIVEN_REF", "USED_REF",
				"BAM_EDIT", "HGVS_OFFSET", "HGVSg",
			},
			false,
		},
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
			got, err := v.CsqKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("Variant.CsqKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variant.CsqKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
