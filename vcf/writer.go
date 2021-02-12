package vcf

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Writer ...
type Writer struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	header *Header
	err    error
}

// NewWriter ...
func NewWriter(f string) (*Writer, error) {
	format := "v"
	if strings.HasSuffix(f, ".vcf.gz") {
		format = "z"
	}
	if strings.HasSuffix(f, ".bcf") {
		format = "b"
	}
	cmd := exec.Command("bcftools", "view", "--no-version", "-O", format, "-o", f)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return &Writer{}, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return &Writer{}, fmt.Errorf("failed to start process: %w", err)
	}
	return &Writer{cmd: cmd, stdin: stdin}, nil
}

// Write ...
func (w Writer) Write(p []byte) (int, error) {
	return w.stdin.Write(p)
}

// WriteString ...
func (w Writer) WriteString(s string) (int, error) {
	return io.WriteString(w.stdin, s)
}

// Close ...
func (w Writer) Close() error {
	err := w.stdin.Close()
	if err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}
	err = w.cmd.Wait()
	return err
	// return w.stdin.Close()
}

// WriteHeader ...
func (w *Writer) WriteHeader(h Header) error {
	w.header = &h
	var line string
	line = fmt.Sprintf("##fileformat=VCFv%2.1f", h.version)
	io.WriteString(w, line+"\n")
	for _, filter := range h.Filters() {
		io.WriteString(w, filter.AsVCFString()+"\n")
	}
	for _, format := range h.Formats() {
		io.WriteString(w, format.AsVCFString()+"\n")
	}
	for _, info := range h.Infos() {
		io.WriteString(w, info.AsVCFString()+"\n")
	}
	// ALT
	// for _, alt := range h.Alts() {
	// 	line := fmt.Sprintf(`##ALT=<ID=%s,Description="%s">`, alt.ID, alt.Description)
	// 	io.WriteString(w, line+"\n")
	// }
	// SAMPLE
	// PEDIGREE
	for _, other := range h.Others() {
		io.WriteString(w, other.AsVCFString()+"\n")
	}
	for _, contig := range h.Contigs() {
		io.WriteString(w, contig.AsVCFString()+"\n")
	}
	columns := []string{
		"#CHROM", "POS", "ID", "REF", "ALT", "QUAL", "FILTER", "INFO",
	}
	if len(h.Samples) > 0 {
		columns = append(columns, "FORMAT")
		columns = append(columns, h.Samples...)
	}
	line = strings.Join(columns, "\t")
	io.WriteString(w, line+"\n")
	return nil
}

// WriteVariant adds the variant to the writer. Returns non-nil error if the
// variant can not be written or it is invalid, for example, if the writers
// header defines contigs and its Chrom is not defined in the header.
func (w Writer) WriteVariant(v Variant) error {
	// There should be more validation before adding the variant
	if w.header == nil {
		return fmt.Errorf("Writer has no header, unable to add variants %v", w.header)
	}
	contigs := w.header.Contigs()
	// Only check if there is a corresponding contig in the header if there are
	// any contigs in the header (contig header lines are not required).
	if len(contigs) > 0 {
		hasContig := false
		for _, c := range contigs {
			if c.ID() == v.Chrom {
				hasContig = true
			}
		}
		if !hasContig {
			return fmt.Errorf("header missing contig %s", v.Chrom)
		}
	}

	for _, f := range v.Filter {
		if !hasID(w.header.Filters(), f) {
			return fmt.Errorf("filter %s not found in header", f)
		}
	}
	// There is more validation of INFO and FORMAT fields that could be done.
	// For example, could validate that the mapping values correspond to what's
	// in the header.

	for k := range v.Info {
		if !hasID(w.header.Infos(), k) {
			return fmt.Errorf("info %s not found in header", k)
		}
	}

	for _, f := range v.Format {
		if !hasID(w.header.Formats(), f) {
			return fmt.Errorf("format %s not found in header", f)
		}
	}

	// If the genotype has samples they must match the header. If the header has
	// samples there must all be present.
	gn := []string{}
	for _, g := range v.genotypes {
		gn = append(gn, g.Name)
	}
	if !stringSliceEqual(gn, w.header.Samples) {
		return fmt.Errorf("the genotype samples do not match the samples in the header")
	}

	_, err := io.WriteString(w, v.AsVCFLine()+"\n")
	return err
}

// Do two string slices contain the same elements in the same order?
func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func hasID(lines []HeaderLine, id string) bool {
	for _, l := range lines {
		if l.ID() == id {
			return true
		}
	}
	return false
}
