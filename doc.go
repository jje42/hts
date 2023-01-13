/*
Example of filtering a VCF:

v, err := vcf.New("test.vcf")
if err != nil {
	panic(err)
}
w, err := vcf.NewWriter("out.bcf")
if err != nil {
	panic(err)
}
defer w.Close()
h := v.Header
h.AddHeaderLines(
	vcf.NewComplexHeaderLine(
		"FILTER",
		map[string]string{
			"ID":          "RANDOM_FILTER",
			"Description": "Randomly filter variants",
		}),
)
w.WriteHeader(h)
scanner, err := vcf.NewScanner(v)
if err != nil {
	panic(err)
}
for scanner.Scan() {
	vc := scanner.Variant()
	vc.Filter = append(vc.Filter, "RANDOM_FILTER")
	err := w.WriteVariant(vc)
	if err != nil {
		panic(err)
	}
}
if err := scanner.Err(); err != nil {
	panic(err)
}
err = vcf.CreateIndex("out.bcf")
if err != nil {
	panic(err)
}
*/
