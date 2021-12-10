package vcf

import (
	"reflect"
	"testing"
)

func TestGenotype_Alleles(t *testing.T) {
	type fields struct {
		Name          string
		values        map[string]string
		alleleIndexes []int
		phased        bool
		v             *Variant
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{"t1", fields{"", map[string]string{"GT": "0/0"}, []int{0, 0}, false, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "A"}, false},
		{"t2", fields{"", map[string]string{"GT": "0/1"}, []int{0, 1}, false, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "C"}, false},
		{"t3", fields{"", map[string]string{"GT": "1/1"}, []int{1, 1}, false, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"C", "C"}, false},
		{"t4", fields{"", map[string]string{"GT": "0|0"}, []int{0, 0}, true, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "A"}, false},
		{"t5", fields{"", map[string]string{"GT": "0|1"}, []int{0, 1}, true, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "C"}, false},
		{"t6", fields{"", map[string]string{"GT": "1|1"}, []int{1, 1}, true, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"C", "C"}, false},
		{"t7", fields{"", map[string]string{}, []int{}, false, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{}, true},
		{"t8", fields{"", map[string]string{"GT": "0/2"}, []int{0, 2}, false, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Genotype{
				Name:          tt.fields.Name,
				values:        tt.fields.values,
				alleleIndexes: tt.fields.alleleIndexes,
				phased:        tt.fields.phased,
				v:             tt.fields.v,
			}
			got, err := g.Alleles()
			if (err != nil) != tt.wantErr {
				t.Errorf("Genotype.Alleles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Genotype.Alleles() = %v, want %v", got, tt.want)
			}
		})
	}
}
