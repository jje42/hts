package vcf

import (
	"fmt"
	"strconv"
	"strings"
)

type Variant struct {
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

func (v Variant) Sample(name string) (Genotype, error) {
	for _, g := range v.genotypes {
		if g.Name == name {
			return g, nil
		}
	}
	return Genotype{}, fmt.Errorf("no genotype for %s", name)
}

func (v Variant) Genotypes() []Genotype {
	return v.genotypes
}

func (v Variant) Alleles() []string {
	xs := []string{v.Ref}
	xs = append(xs, v.Alt...)
	return xs
}

func (v Variant) HasAttribute(key string) bool {
	_, ok := v.Info[key]
	return ok
}

func (v Variant) Attribute(key string) (string, error) {
	value, ok := v.Info[key]
	if !ok {
		return "", fmt.Errorf("no such info: %s", key)
	}
	return value, nil
}

func (v Variant) AttributeAsInt(key string) (int, error) {
	value, ok := v.Info[key]
	if !ok {
		return 0, fmt.Errorf("no such info: %s", key)
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %s to int: %w", key, err)
	}
	return i, nil
}

func (v Variant) AttributeAsFloat64(key string) (float64, error) {
	value, ok := v.Info[key]
	if !ok {
		return 0.0, fmt.Errorf("no such info: %s", key)
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("unable to convert %s to float64: %w", key, err)
	}
	return f, nil
}

// func (v Variant) AttributeAsStringSlice(key string) ([]string, error) {

// }

// func (v Variant) AttributeAsIntSlice(key string) ([]int, error) {

// }

// func (v Variant) AttributeAsFloat64Slice(key string) ([]float64, error) {

// }

func (v Variant) IsFiltered() bool {
	switch len(v.Filter) {
	case 0:
		return false
	case 1:
		if v.Filter[0] == "PASS" || v.Filter[0] == "." {
			return false
		}
		return true
	default:
		return true
	}
}

func (v Variant) IsNotFiltered() bool {
	return !v.IsFiltered()
}

func (v Variant) Type() Type {
	if len(v.Alt) == 0 {
		return NO_VARIATION
	}
	t := biallelicType(v.Ref, v.Alt[0])
	for _, alt := range v.Alt[1:] {
		bt := biallelicType(v.Ref, alt)
		if bt != t {
			return MIXED
		}
	}
	return t
}

// Type is a type of variation.
type Type int

const (
	// NO_VARIATION
	NO_VARIATION Type = iota
	// SNP ...
	SNP
	// MNP ...
	MNP
	// INDEL ...
	INDEL
	// SYMBOLIC ...
	SYMBOLIC
	// MIXED ...
	MIXED
)

func biallelicType(ref, alt string) Type {
	if containsAny(alt, []string{"*", "<", "[", "]", "."}) {
		return SYMBOLIC
	}
	if len(ref) == len(alt) {
		if len(alt) == 1 {
			return SNP
		}
		return MNP
	}
	return INDEL
}

func containsAny(s string, set []string) bool {
	for _, i := range set {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}

func (v Variant) IsSNP() bool {
	return v.Type() == SNP
}

func (v Variant) IsINDEL() bool {
	return v.Type() == INDEL
}

// Type()
// Start(), End()
// IsNotFiltered(), IsIndel(), IsComplexIndel(), IsBiallelic(), IsSimpleDeletion(), IsSymbolicOrSV(), IsVariant(), IsSNP(), IsSimpleIndel()

// v := vcf.NewVariant(&header)

func (v *Variant) AddGenotype(g Genotype) error {
	// We could addtionally validate the values match the expectation of the
	// header, but we do not have access to the header here.
	g.v = v // set the correct reference
	if len(g.values) > len(v.Format) {
		return fmt.Errorf("genotypes contains tags not listed in variant FORMAT field")
	}
	for _, format := range v.Format {
		_, ok := g.values[format]
		if !ok {
			return fmt.Errorf("genotype is missing %s key", format)
		}
	}
	v.genotypes = append(v.genotypes, g)
	return nil
}

func (v Variant) AsVCFLine() string {
	info := []string{}
	for k, v := range v.Info {
		info = append(info, strings.Join([]string{k, v}, "="))
	}
	qual := v.Qual
	if qual == "" {
		qual = "."
	}
	filter := "."
	if len(v.Filter) > 0 {
		filter = strings.Join(v.Filter, ";")
	}
	cols := []string{
		v.Chrom,
		fmt.Sprint(v.Pos),
		v.ID,
		v.Ref,
		strings.Join(v.Alt, ","),
		qual,
		filter,
		strings.Join(info, ";"),
	}
	if len(v.genotypes) > 0 {
		cols = append(cols, strings.Join(v.Format, ":"))
		for _, g := range v.genotypes {
			cols = append(cols, g.AsVCFString())
		}
	}
	return strings.Join(cols, "\t")
}

func parseVcfLine(line string, samples []string) (Variant, error) {
	bits := strings.Split(line, "\t")
	if len(bits) < 8 {
		return Variant{}, fmt.Errorf("less than 8 columns found in VCF line")
	}
	pos, err := strconv.Atoi(bits[1])
	if err != nil {
		return Variant{}, fmt.Errorf("unable to convert position: %w", err)
	}
	info := make(map[string]string)
	// A '.' in the INFO column indicates that there are no fields, do not
	// add this to the map!
	if bits[7] != "." {
		for _, i := range strings.Split(bits[7], ";") {
			bits := strings.SplitN(i, "=", 2)
			if len(bits) == 2 {
				info[bits[0]] = bits[1]
			} else {
				info[bits[0]] = "1"
			}
		}
	}
	filter := []string{}
	for _, i := range strings.Split(bits[6], ";") {
		if i != "PASS" && i != "." {
			filter = append(filter, i)
		}
	}
	vc := Variant{
		Chrom: bits[0],
		Pos:   pos,
		ID:    bits[2],
		Ref:   bits[3],
		Alt:   strings.Split(bits[4], ","),
		Qual:  bits[5],
		// Filter: strings.Split(bits[6], ";"),
		Filter: filter,
		Info:   info,
	}
	if len(bits) >= 9 {
		vc.Format = strings.Split(bits[8], ":")
	}
	for i, gt := range bits[9:] {
		xs := strings.Split(gt, ":")
		vs := make(map[string]string)
		for j, x := range xs {
			vs[vc.Format[j]] = x
		}
		g, err := NewGenotype(samples[i], vs)
		if err != nil {
			return Variant{}, fmt.Errorf("unable to create genotype: %w", err)
		}
		vc.AddGenotype(g)
	}

	return vc, nil
}

func stringSliceContains(xs []string, key string) bool {
	for _, x := range xs {
		if x == key {
			return true
		}
	}
	return false
}
