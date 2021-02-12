# HTS: High Throughput Sequencing file parsing

Currently, only some basic VCF reading/writing is supported. It uses
`bcftools` (which must be on the PATH) to decode/encode files, so it can
handle `.vcf`, `.vcf.gz` and `.bcf` files. It is still a working in progress
but functional.

## Example

Reading VCF files should be familiar if you've used `bufio.Scanner` before:

```go
import (
        "log"
        "fmt"
        "github.com/jje42/hts/vcf"
)

func main() {
	v, err := vcf.New("test.vcf.gz")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("VCF with %d samples", len(v.Header.Samples))
	scanner, err := vcf.NewScanner(v)
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		vc := scanner.Variant()
		fmt.Printf("%+v\n", vc)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
```

You can write VCFs as well:

```go
func copyVCF() {
	r, err := vcf.New("input.vcf.gz")
	if err != nil {
		log.Fatal(err)
	}
	w, err := vcf.NewWriter("output.vcf.gz")
	if err != nil {
		log.Fatalf("failed to create writer: %v", err)
	}
	w.WriteHeader(r.Header)
	scanner, err := vcf.NewScanner(r)
	if err != nil {
		log.Fatalf("failed to create scanner: %v", err)
	}
	for scanner.Scan() {
		v := scanner.Variant()
		err := w.WriteVariant(v)
		if err != nil {
			log.Fatalf("failed to write variant: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanning failed: %v", err)
	}
	w.Close()
}
```

You can also create VCFs from scratch:

```go
func createVCF() {
	h := vcf.NewHeader()
	h.AddHeaderLines(vcf.StandardHeaderLines()...)
	h.AddHeaderLines(vcf.NewComplexHeaderLine("contig", map[string]string{
		"ID":     "1",
		"length": "249250621",
	}))
	h.Samples = []string{"HG001"}
	w, err := vcf.NewWriter("my.vcf.gz")
	if err != nil {
		log.Fatalf("failed to create writer: %v", err)
	}
	defer w.Close()

	err = w.WriteHeader(h)
	if err != nil {
		log.Fatalf("failed to write header: %v", err)
	}
	v := vcf.Variant{
		Chrom:  "1",
		Pos:    10,
		ID:     "rs123",
		Ref:    "A",
		Alt:    []string{"C"},
		Qual:   "",
		Filter: []string{"PASS"},
		Info: map[string]string{
			"AC": "100",
		},
		Format: []string{"GT"},
	}
	g, err := vcf.NewGenotype("HG001", map[string]string{"GT": "0/1"})
	if err != nil {
		log.Fatalf("failed to create genotype: %v", err)
	}
	v.AddGenotype(g)
	err = w.WriteVariant(v)
	if err != nil {
		log.Fatal(err)
	}
}
```