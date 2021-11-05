package vcf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Genotype represents an individual genotyped sample in a VCF file.
type Genotype struct {
	Name          string
	values        map[string]string
	alleleIndexes []int
	phased        bool
	v             *Variant
}

func NewGenotype(name string, attributes map[string]string) (Genotype, error) {
	// The name should be present in the header, but we don't have access to the
	// header here to check. The check is performed if/when we try to add a
	// variant to a Writer (after adding the genotype to the variant).
	g := Genotype{
		Name:   name,
		values: attributes,
	}
	// It is not an error if there is no GT field. The specification only says
	// that it must be the first field if present.
	gt, ok := attributes["GT"]
	if ok {
		// The TSO500 Local App puts these non-standard genotypes in its
		// output. There appear when all reads in the AD count support
		// the alt allele but the DP is higher. Presumably some reads
		// are filtered.
		if gt == "1/." {
			gt = "1/1"
		}
		// Only attempt to convert the indexes if they are not no-calls. What about non-diplody organisms.
		if gt != "./." && gt != ".|." && gt != "." {
			sep := "/"
			if strings.Contains(gt, "|") {
				sep = "|"
				g.phased = true
			}
			for _, istr := range strings.Split(gt, sep) {
				i, err := strconv.Atoi(istr)
				if err != nil {
					return Genotype{}, fmt.Errorf("unable to convert %s to int: %w", istr, err)
				}
				g.alleleIndexes = append(g.alleleIndexes, i)
			}
		}
	}
	return g, nil
}

// Allele(i int) what should this return?

// Alleles returns a slice containing the alleles in this sample. May be an
// empty slice if the sample was not called.
func (g Genotype) Alleles() ([]string, error) {
	xs := []string{}
	alleles := g.v.Alleles()
	if len(g.alleleIndexes) == 0 {
		return []string{}, errors.New("genotype has no alleles")
	}
	for _, i := range g.alleleIndexes {
		if len(alleles) < i+1 {
			return []string{}, fmt.Errorf("GT has index %d, but the variant only has %d alleles", i, len(alleles))
		}
		xs = append(xs, alleles[i])

	}
	return xs, nil
}

// IsPhased returns true if the alleles are phased.
func (g Genotype) IsPhased() bool {
	return g.phased
}

// Ploidy returns the ploidy of this genotype; 0 if the site is no-called.
func (g Genotype) Ploidy() int {
	return len(g.alleleIndexes)
}

// GenotypeString() string

// IsCalled returns true if this genotype is comprised of alleles that are all
// called.
func (g Genotype) IsCalled() bool {
	sep := "/"
	if g.IsPhased() {
		sep = "|"
	}
	gt, err := g.Attribute("GT")
	if err != nil {
		return false
	}
	for _, i := range strings.Split(gt, sep) {
		if i == "." {
			return false
		}
	}
	return true
}
func intSliceAllEqual(xs []int) bool {
	for i := 1; i < len(xs); i++ {
		if xs[i] != xs[0] {
			return false
		}
	}
	return true
}

func (g Genotype) IsHom() bool {
	return intSliceAllEqual(g.alleleIndexes)
}

func (g Genotype) IsHomRef() bool {
	return g.alleleIndexes[0] == 0 && intSliceAllEqual(g.alleleIndexes)
}

func (g Genotype) IsHomVar() bool {
	return g.alleleIndexes[0] != 0 && intSliceAllEqual(g.alleleIndexes)
}

func (g Genotype) IsHet() bool {
	return !g.IsHom()
}

func (g Genotype) IsHetNonRef() bool {
	if g.IsHom() {
		return false
	}
	for _, i := range g.alleleIndexes {
		if i == 0 {
			return false
		}
	}
	return true
}

// IsHet(), IsHetNonRef(), IsHom(), IsHomRef(), IsHomVar(),

func (g Genotype) IsNoCall() bool {
	return !g.IsCalled()
}

// Attribute returns any genotype attribute as a string. It returns a non-nil
// error if key is not an attribute.
func (g Genotype) Attribute(key string) (string, error) {
	v, ok := g.values[key]
	if !ok {
		return "", fmt.Errorf("no such attribute: %s", key)
	}
	return v, nil
}

// AttributeAsInt returns any genotype attribute as a int. It returns a non-nil
// error if key is not an attribute or it can not be converted to an int.
func (g Genotype) AttributeAsInt(key string) (int, error) {
	v, err := g.Attribute(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("unable to parse %s as int: %w", key, err)
	}
	return i, nil
}

// AttributeAsFloat64 returns any genotype attribute as a float64. It returns a
// non-nil error if key is not an attribute or it can not be converted to an
// float64.
func (g Genotype) AttributeAsFloat64(key string) (float64, error) {
	v, ok := g.values[key]
	if !ok {
		return 0.0, fmt.Errorf("no such key: %s", key)
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0.0, fmt.Errorf("unable to parse %s as float: %w", key, err)
	}
	return f, nil
}

func (g Genotype) AsVCFString() string {
	xs := []string{}
	for _, format := range g.v.Format {
		// v, ok := g.values[format]
		xs = append(xs, g.values[format])
	}
	return strings.Join(xs, ":")
}
