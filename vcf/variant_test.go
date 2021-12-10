package vcf

import "testing"

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
