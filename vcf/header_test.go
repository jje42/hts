package vcf

import (
	"reflect"
	"testing"
)

func Test_parseOtherHeaderLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"t1", args{"##bcftools_annotateVersion=1.9+htslib-1.9"}, "bcftools_annotateVersion", "1.9+htslib-1.9"},
		{"t2", args{"##filedate=20151210"}, "filedate", "20151210"},
		{"t3", args{`##source="simplfy-vcf (r1211)"`}, "source", `"simplfy-vcf (r1211)"`},
		{"t4", args{"##foobar"}, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseOtherHeaderLine(tt.args.s)
			if got != tt.want {
				t.Errorf("parseOtherHeaderLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseOtherHeaderLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func TestOtherHeaderLine_String(t *testing.T) {
	type fields struct {
		ID    string
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"t1", fields{"bcftools_annotateVersion", "1.9+htslib-1.9"}, "##bcftools_annotateVersion=1.9+htslib-1.9"},
		{"t2", fields{"filedate", "20151210"}, "##filedate=20151210"},
		{"t3", fields{"source", `"simplfy-vcf (r1211)"`}, `##source="simplfy-vcf (r1211)"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := OtherHeaderLine{
				ID:    tt.fields.ID,
				Value: tt.fields.Value,
			}
			if got := h.String(); got != tt.want {
				t.Errorf("OtherHeaderLine.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseContigHeaderLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Contig
		wantErr bool
	}{
		{"t1", args{"##contig=<ID=1,length=249250621,assembly=b37>"}, Contig{"1", map[string]string{"ID": "1", "length": "249250621", "assembly": "b37"}}, false},
		{"t2", args{"##contig=<ID=GL000207.1,length=4262,assembly=b37>"}, Contig{"GL000207.1", map[string]string{"ID": "GL000207.1", "length": "4262", "assembly": "b37"}}, false},
		{"t3", args{"##contig=<ID=1,length=249250621>"}, Contig{"1", map[string]string{"ID": "1", "length": "249250621"}}, false},
		{"t4", args{"##contig=<ID=1>"}, Contig{"1", map[string]string{"ID": "1"}}, false},
		{"t5", args{"##contig=<length=249250621>"}, Contig{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContigHeaderLine(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContigHeaderLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContigHeaderLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
