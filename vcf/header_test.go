package vcf

import (
	"reflect"
	"testing"
)

func Test_parseHeaderLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    HeaderLine
		wantErr bool
	}{
		{"t1", args{"##bcftools_annotateVersion=1.9+htslib-1.9"}, HeaderLine{Key: "bcftools_annotateVersion", Value: "1.9+htslib-1.9"}, false},
		{"t2", args{"##filedate=20151210"}, HeaderLine{Key: "filedate", Value: "20151210"}, false},
		{"t3", args{`##source="simplfy-vcf (r1211)"`}, HeaderLine{Key: "source", Value: `"simplfy-vcf (r1211)"`}, false},
		{"t4", args{"##foobar"}, HeaderLine{}, true},
		{"t1", args{"##contig=<ID=1,length=249250621,assembly=b37>"}, HeaderLine{"contig", "", map[string]string{"ID": "1", "length": "249250621", "assembly": "b37"}}, false},
		{"t2", args{"##contig=<ID=GL000207.1,length=4262,assembly=b37>"}, HeaderLine{"contig", "", map[string]string{"ID": "GL000207.1", "length": "4262", "assembly": "b37"}}, false},
		{"t3", args{"##contig=<ID=1,length=249250621>"}, HeaderLine{"contig", "", map[string]string{"ID": "1", "length": "249250621"}}, false},
		{"t4", args{"##contig=<ID=1>"}, HeaderLine{"contig", "", map[string]string{"ID": "1"}}, false},
		{"t5", args{"##contig=<length=249250621>"}, HeaderLine{}, true},
		{"t6", args{`##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">`}, HeaderLine{"FORMAT", "", map[string]string{"ID": "GT", "Number": "1", "Type": "String", "Description": "Genotype"}}, false},
		{"t7", args{`##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Minimum GenCall score, encoded as a phred quality integer.",Source="description",Version="128">`}, HeaderLine{"FORMAT", "", map[string]string{"ID": "GQ", "Number": "1", "Type": "Integer", "Description": "Minimum GenCall score, encoded as a phred quality integer.", "Source": "description", "Version": "128"}}, false},
		// Missing ID
		{"t8", args{`##FORMAT=<Number=1,Type=String,Description="Genotype">`}, HeaderLine{}, true},
		// Missing Number
		{"t9", args{`##FORMAT=<ID=GT,Type=String,Description="Genotype">`}, HeaderLine{}, true},
		// Missing Type
		{"t10", args{`##FORMAT=<ID=GT,Number=1,Description="Genotype">`}, HeaderLine{}, true},
		// Missing Description
		{"t11", args{`##FORMAT=<ID=GT,Number=1,Type=String,>`}, HeaderLine{}, true},
		// Source not quoted
		// {"t12", args{`##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype",Source=description>`}, HeaderLine{}, true},
		{"t13", args{`##INFO=<ID=AC,Number=A,Type=Integer,Description="Allele count in genotypes">`}, HeaderLine{"INFO", "", map[string]string{"ID": "AC", "Number": "A", "Type": "Integer", "Description": "Allele count in genotypes"}}, false},
		{"t13", args{`##INFO=<ID=AC,Number=A,Type=Integer,Description="Allele count in genotypes",Source="description",Version="128">`}, HeaderLine{"INFO", "", map[string]string{"ID": "AC", "Number": "A", "Type": "Integer", "Description": "Allele count in genotypes", "Source": "description", "Version": "128"}}, false},
		{"t14", args{`##INFO=<Number=A,Type=Integer,Description="Allele count in genotypes">`}, HeaderLine{}, true},
		{"t15", args{`##INFO=<ID=AC,Type=Integer,Description="Allele count in genotypes">`}, HeaderLine{}, true},
		{"t16", args{`##INFO=<ID=AC,Number=A,Description="Allele count in genotypes">`}, HeaderLine{}, true},
		{"t17", args{`##INFO=<ID=AC,Number=A,Type=Integer>`}, HeaderLine{}, true},
		{"t18", args{`##FILTER=<ID=LowQual,Description="Low quality">`}, HeaderLine{"FILTER", "", map[string]string{"ID": "LowQual", "Description": "Low quality"}}, false},
		{"t19", args{`##FILTER=<ID=LowQual,Description="Low quality",Source="description",Version="128">`}, HeaderLine{"FILTER", "", map[string]string{"ID": "LowQual", "Description": "Low quality", "Source": "description", "Version": "128"}}, false},
		{"t20", args{`##FILTER=<Description="Low quality">`}, HeaderLine{}, true},
		{"t21", args{`##FILTER=<ID=LowQual>`}, HeaderLine{}, true},
		// ALT, SAMPLE, PEDIGREE
		// listed in specs DEL INS DUP INV CNV DUP:TANDEM DEL:ME INS:ME (exclusive?)
		{"t22", args{`##ALT=<ID=DEL,Description="description">`}, HeaderLine{"ALT", "", map[string]string{"ID": "DEL", "Description": "description"}}, false},
		{"t23", args{`##ALT=<ID=DEL,Description="description",Source="description",Version="128">`}, HeaderLine{"ALT", "", map[string]string{"ID": "DEL", "Description": "description", "Source": "description", "Version": "128"}}, false},
		{"t24", args{`##ALT=<Description="description">`}, HeaderLine{}, true},
		{"t25", args{`##ALT=<ID=DEL>`}, HeaderLine{}, true},
		{"t26", args{`##SAMPLE=<ID=Blood,Genomes=Germline,Mixture=1.,Description="Patient germline genome">`}, HeaderLine{"SAMPLE", "", map[string]string{"ID": "Blood", "Genomes": "Germline", "Mixture": "1.", "Description": "Patient germline genome"}}, false},
		{
			"t27",
			args{`##SAMPLE=<ID=TissueSample,Genomes=Germline;Tumor,Mixture=.3;.7,Description="Patient germline genome;Patient tumor genome">`},
			HeaderLine{
				"SAMPLE",
				"",
				map[string]string{
					"ID":          "TissueSample",
					"Genomes":     "Germline;Tumor",
					"Mixture":     ".3;.7",
					"Description": "Patient germline genome;Patient tumor genome",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaderLine(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeaderLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeaderLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeaderLine_AsVCFString(t *testing.T) {
	type fields struct {
		Key     string
		Value   string
		mapping map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"t1", fields{Key: "bcftools_annotateVersion", Value: "1.9+htslib-1.9"}, "##bcftools_annotateVersion=1.9+htslib-1.9"},
		{"t2", fields{Key: "filedate", Value: "20151210"}, "##filedate=20151210"},
		{"t3", fields{Key: "source", Value: `"simplfy-vcf (r1211)"`}, `##source="simplfy-vcf (r1211)"`},
		// For contig length and assembly are stored in map it will be
		// purely coincidental if these tests pass.
		{"t4", fields{"contig", "", map[string]string{"ID": "1", "length": "249250621", "assembly": "b37"}}, "##contig=<ID=1,length=249250621,assembly=b37>"},
		{"t5", fields{"contig", "", map[string]string{"ID": "GL000207.1", "length": "4262", "assembly": "b37"}}, "##contig=<ID=GL000207.1,length=4262,assembly=b37>"},
		{"t6", fields{"contig", "", map[string]string{"ID": "1", "length": "249250621"}}, "##contig=<ID=1,length=249250621>"},
		{"t7", fields{"contig", "", map[string]string{"ID": "1"}}, "##contig=<ID=1>"},
		{"t8", fields{"FORMAT", "", map[string]string{"ID": "GT", "Number": "1", "Type": "String", "Description": "Genotype"}}, `##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">`},
		// {"t9", fields{"FORMAT", "", map[string]string{"ID": "GQ", "Number": "1", "Type": "Integer", "Description": "Minimum GenCall score, encoded as a phred quality integer.", "Source": "description", "Version": "128"}}, `##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Minimum GenCall score, encoded as a phred quality integer.",Source="description",Version="128">`},
		{"t10", fields{"INFO", "", map[string]string{"ID": "AC", "Number": "A", "Type": "Integer", "Description": "Allele count in genotypes"}}, `##INFO=<ID=AC,Number=A,Type=Integer,Description="Allele count in genotypes">`},
		{"t11", fields{"INFO", "", map[string]string{"ID": "AC", "Number": "A", "Type": "Integer", "Description": "Allele count in genotypes", "Source": "description", "Version": "128"}}, `##INFO=<ID=AC,Number=A,Type=Integer,Description="Allele count in genotypes",Source="description",Version="128">`},
		{"t12", fields{"FILTER", "", map[string]string{"ID": "LowQual", "Description": "Low quality"}}, `##FILTER=<ID=LowQual,Description="Low quality">`},
		// {"t13", fields{"FILTER", "", map[string]string{"ID": "LowQual", "Description": "Low quality", "Source": "description", "Version": "128"}}, `##FILTER=<ID=LowQual,Description="Low quality",Source="description",Version="128">`},
		{"t14", fields{"ALT", "", map[string]string{"ID": "DEL", "Description": "description"}}, `##ALT=<ID=DEL,Description="description">`},
		// {"t15", fields{"ALT", "", map[string]string{"ID": "DEL", "Description": "description", "Source": "description", "Version": "128"}}, `##ALT=<ID=DEL,Description="description",Source="description",Version="128">`},
		// SAMPLE the order of keys is not defined is specs
		// {
		// 	"t16",
		// 	fields{
		// 		"SAMPLE",
		// 		"",
		// 		map[string]string{
		// 			"ID":          "Blood",
		// 			"Genomes":     "Germline",
		// 			"Mixture":     "1.",
		// 			"Description": "Patient germline genome",
		// 		},
		// 	},
		// 	`##SAMPLE=<ID=Blood,Genomes=Germline,Mixture=1.,Description="Patient germline genome">`,
		// },
		// {
		// 	"t17",
		// 	fields{
		// 		"SAMPLE",
		// 		"",
		// 		map[string]string{
		// 			"ID":          "TissueSample",
		// 			"Genomes":     "Germline;Tumor",
		// 			"Mixture":     ".3;.7",
		// 			"Description": "Patient germline genome;Patient tumor genome",
		// 		},
		// 	},
		// 	`##SAMPLE=<ID=TissueSample,Genomes=Germline;Tumor,Mixture=.3;.7,Description="Patient germline genome;Patient tumor genome">`,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HeaderLine{
				Key:     tt.fields.Key,
				Value:   tt.fields.Value,
				mapping: tt.fields.mapping,
			}
			if got := h.AsVCFString(); got != tt.want {
				t.Errorf("HeaderLine.AsVCFString() = %v, want %v", got, tt.want)
			}
		})
	}
}
