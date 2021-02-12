package vcf

import (
	"reflect"
	"testing"
)

func TestGenotype_Alleles(t *testing.T) {
	type fields struct {
		Name   string
		values map[string]string
		v      *Variant
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{"t1", fields{"", map[string]string{"GT": "0/0"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "A"}, false},
		{"t2", fields{"", map[string]string{"GT": "0/1"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "C"}, false},
		{"t3", fields{"", map[string]string{"GT": "1/1"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"C", "C"}, false},
		{"t4", fields{"", map[string]string{"GT": "0|0"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "A"}, false},
		{"t5", fields{"", map[string]string{"GT": "0|1"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"A", "C"}, false},
		{"t6", fields{"", map[string]string{"GT": "1|1"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{"C", "C"}, false},
		{"t7", fields{"", map[string]string{}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{}, true},
		{"t8", fields{"", map[string]string{"GT": "0/2"}, &Variant{Ref: "A", Alt: []string{"C"}}}, []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Genotype{
				Name:   tt.fields.Name,
				values: tt.fields.values,
				v:      tt.fields.v,
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
