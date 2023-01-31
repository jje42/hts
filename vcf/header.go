package vcf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const metadataIndicator = "##"
const headerIndicator = "#"

// Header is a VCF header.
type Header struct {
	version float64
	lines   []HeaderLine
	Samples []string
}

// NewHeader creates a new VCF header with the default file format version.
func NewHeader() Header {
	return Header{
		version: 4.2,
	}
}

// Version returns the VCF file format version.
func (h Header) Version() float64 {
	return h.version
}

// AddHeaderLines add the header lines to the header.
func (h *Header) AddHeaderLines(lines ...HeaderLine) {
	h.lines = append(h.lines, lines...)
}

// HeaderLines returns all of the lines in the header.
func (h Header) HeaderLines() []HeaderLine {
	return h.lines
}

func (h Header) allHeaderLines(key string) []HeaderLine {
	xs := []HeaderLine{}
	for _, l := range h.lines {
		if l.Key == key {
			xs = append(xs, l)
		}
	}
	return xs
}

// Filters returns all FILTER header lines.
func (h Header) Filters() []HeaderLine {
	return h.allHeaderLines("FILTER")
}

// Infos returns all INFO header lines.
func (h Header) Infos() []HeaderLine {
	return h.allHeaderLines("INFO")
}

// Formats returns all FORMAT header lines.
func (h Header) Formats() []HeaderLine {
	return h.allHeaderLines("FORMAT")
}

// Contigs returns all contig header lines
func (h Header) Contigs() []HeaderLine {
	return h.allHeaderLines("contig")
}

// Others returns all non- FILTER, INFO, FORMAT and contig header lines.
func (h Header) Others() []HeaderLine {
	xs := []HeaderLine{}
	for _, l := range h.lines {
		if l.Key != "FILTER" && l.Key != "INFO" && l.Key != "FORMAT" && l.Key != "contig" {
			xs = append(xs, l)
		}
	}
	return xs
}

// func parseHeader(r io.Reader) (Header, error)
// func parseHeaderFromStringSlice(headerLines []string) (Header, error)
func readHeaderFromFile(path string) (Header, error) {
	if _, err := os.Stat(path); err != nil {
		return Header{}, fmt.Errorf("can not stat file: %w", err)
	}
	exe, err := findBcftools()
	if err != nil {
		return Header{}, err
	}
	cmd := exec.Command(exe, "view", "--no-version", "-h", path)
	bs, err := cmd.Output()
	if err != nil {
		return Header{}, fmt.Errorf("process failed: %v", err)
	}
	headerLines := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(bs))
	for scanner.Scan() {
		line := scanner.Text()
		headerLines = append(headerLines, line)
	}
	if err := scanner.Err(); err != nil {
		return Header{}, fmt.Errorf("scanning header failed: %v", err)
	}
	return parseHeader(headerLines)
}

// would this be better to accept io.Reader instead of []string?
func parseHeader(headerLines []string) (Header, error) {
	h := Header{}
	for _, line := range headerLines {
		if strings.HasPrefix(line, metadataIndicator) {
			var headerLine HeaderLine
			var err error
			if strings.Contains(line, "<") {
				headerLine, err = parseHeaderLine(line)
				if err != nil {
					return Header{}, err
				}
			} else {
				bits := strings.SplitN(strings.TrimPrefix(line, "##"), "=", 2)
				headerLine = HeaderLine{Key: bits[0], Value: bits[1]}
				if headerLine.Key == "fileformat" {
					v := strings.TrimPrefix(headerLine.Value, "VCFv")
					f, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return Header{}, fmt.Errorf("unable to parse version: %w", err)
					}
					h.version = f
				}
			}
			h.AddHeaderLines(headerLine)
		}
		if strings.HasPrefix(line, "#CHROM") {
			bits := strings.Split(line, "\t")
			if len(bits) > 9 {
				h.Samples = bits[9:]
			}
		}
	}
	if h.version == 0 {
		return Header{}, errors.New("VCF has no version number")
	}
	return h, nil
}

// HeaderLine is a line from a VCF header.
type HeaderLine struct {
	Key     string
	Value   string
	mapping map[string]string
}

// ID returns the ID tag from the header line. If the line has no ID tag an
// empty string will be returned.
func (h HeaderLine) ID() string {
	return h.mapping["ID"]
}

func (h HeaderLine) Get(id string) string {
	return h.mapping[id]
}

// AsVCFString returns the header line in the format expected in a VCF header.
func (h HeaderLine) AsVCFString() string {
	// The VCF specification is vague on whether the tags in a header line
	// should appear in any particular order or what that order might be.
	// This is problematic for testing as the tags are stored in a map and
	// may appear in any order. Is there a single correct way to write every
	// header or is it always ambiguous?
	if h.Value == "" {
		// if h.Key == "SAMPLE" {
		// 	s, _ := sampleHeaderAsVCFString(h)
		// 	return s
		// }
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("##%s=<", h.Key))
		pairs := []string{}
		id, ok := h.mapping["ID"]
		if ok {
			pairs = append(pairs, fmt.Sprintf("ID=%s", id))
		}
		number, ok := h.mapping["Number"]
		if ok {
			pairs = append(pairs, fmt.Sprintf("Number=%s", number))
		}
		_type, ok := h.mapping["Type"]
		if ok {
			pairs = append(pairs, fmt.Sprintf("Type=%s", _type))
		}
		desc, ok := h.mapping["Description"]
		if ok {
			pairs = append(pairs, fmt.Sprintf(`Description="%s"`, desc))
		}
		// It's not specified what order, if any, extra tags should be in.
		for k, v := range h.mapping {
			if k != "ID" && k != "Number" && k != "Type" && k != "Description" {
				// It is not specified in the spec whether
				// contig field values should be quoted or not,
				// but all files I have access to do not quote
				// the values.
				if h.Key == "contig" {
					// Should not be quoted
					pairs = append(pairs, fmt.Sprintf(`%s=%s`, k, v))
				} else {
					pairs = append(pairs, fmt.Sprintf(`%s="%s"`, k, v))
				}
			}
		}
		builder.WriteString(strings.Join(pairs, ","))
		builder.WriteString(">")
		return builder.String()
	}
	return fmt.Sprintf("##%s=%s", h.Key, h.Value)
}

// the examples in v4.3 specs use completely different tags (Assay, Ethnicity
// and Disease). I can only conclude they are completely generic and a user can
// put whatever they want in a SAMPLE header line.
// func sampleHeaderAsVCFString(h HeaderLine) (string, error) {
// 	ids, ok := h.mapping["ID"]
// 	if !ok {
// 		return "", fmt.Errorf("SAMPLE header line has no ID tag")
// 	}
// 	genomes, ok := h.mapping["Genomes"]
// 	if !ok {
// 		return "", fmt.Errorf("SAMPLE header line has no Genomes tag")
// 	}
// 	mixture, ok := h.mapping["Mixture"]
// 	if !ok {
// 		return "", fmt.Errorf("SAMPLE header line has no Mixture tag")
// 	}
// 	description, ok := h.mapping["Description"]
// 	if !ok {
// 		return "", fmt.Errorf("SAMPLE header line has no Description tag")
// 	}
// 	return fmt.Sprintf(`##SAMPLE=<ID=%s,Genomes=%s,Mixture=%s,Description="%s">`, ids, genomes, mixture, description), nil
// }

// NewSimpleHeaderLine creates a "simple" header line: one with just a key and
// value.
func NewSimpleHeaderLine(key, value string) HeaderLine {
	return HeaderLine{Key: key, Value: value}
}

// NewComplexHeaderLine creates a "complex" header line: one with a key and one
// or more key/value pairs of attributes. For example, INFO or FORMAT header
// lines. This is a crap name.
func NewComplexHeaderLine(key string, mapping map[string]string) HeaderLine {
	return HeaderLine{Key: key, mapping: mapping}
}

func parseHeaderLine(s string) (HeaderLine, error) {
	match := regexp.MustCompile(`##(.+?)=`).FindStringSubmatch(s)
	if match == nil {
		return HeaderLine{}, fmt.Errorf("malformed header line")
	}
	headerKey := match[1]
	value := ""
	var mapping map[string]string
	if strings.Contains(s, "<") {
		var err error
		if !strings.Contains(s, ">") {
			return HeaderLine{}, fmt.Errorf("header line is malformed: %s", s)
		}
		mapping, err = parseValues(s)
		if err != nil {
			return HeaderLine{}, fmt.Errorf("failed to parse header line: %w", err)
		}
		switch headerKey {
		case "contig":
			_, ok := mapping["ID"]
			if !ok {
				return HeaderLine{}, fmt.Errorf("%s header line must contain an ID tag", headerKey)
			}
		case "FILTER", "ALT":
			err := checkTags(mapping, []string{"ID", "Description"})
			if err != nil {
				return HeaderLine{}, fmt.Errorf("failed to parse header line: %w", err)
			}
		case "FORMAT", "INFO":
			err := checkTags(mapping, []string{"ID", "Number", "Type", "Description"})
			if err != nil {
				return HeaderLine{}, fmt.Errorf("failed to parse header line: %w", err)
			}
		}
	} else {
		value = strings.SplitN(s, "=", 2)[1]
	}
	return HeaderLine{headerKey, value, mapping}, nil
}

func checkTags(mapping map[string]string, requiredTags []string) error {
	for _, tag := range requiredTags {
		_, ok := mapping[tag]
		if !ok {
			return fmt.Errorf("required tag %s is missing", tag)
		}
	}
	return nil
}

// A little switch machine to parse out the tags: thanks htsjdk team!
func parseValues(s string) (map[string]string, error) {
	var builder strings.Builder
	ret := make(map[string]string)
	key := ""
	index := 0
	inQuote := false
	escape := false

	for _, c := range s {
		if c == '"' {
			if escape {
				_, _ = builder.WriteRune(c)
				escape = false
			} else {
				inQuote = !inQuote
			}
		} else if inQuote {
			if escape {
				if c == '\\' {
					builder.WriteRune(c)
				} else {
					builder.WriteRune('\\')
					builder.WriteRune(c)
				}
				escape = false
			} else if c != '\\' {
				builder.WriteRune(c)
			} else {
				escape = true
			}
		} else {
			escape = false
			switch c {
			case '<':
				if index == 0 {
					break
				}
			case '>':
				if index == len(s)-1 {
					ret[key] = builder.String()
					break
				}
			case '=':
				key = builder.String()
				builder = strings.Builder{}
			case ',':
				ret[key] = builder.String()
				builder = strings.Builder{}
			default:
				builder.WriteRune(c)
			}
		}
		index++
	}
	if inQuote {
		return make(map[string]string), errors.New("unclosed quote in header line")
	}
	return ret, nil
}

// StandardHeaderLines returns a slice of standard VCF headers.
func StandardHeaderLines() []HeaderLine {
	return []HeaderLine{
		NewComplexHeaderLine("FILTER", map[string]string{
			"ID":          "PASS",
			"Description": "All filters passed",
		}),
		NewComplexHeaderLine("FORMAT", map[string]string{
			"ID":          "GT",
			"Number":      "1",
			"Type":        "String",
			"Description": "Genotype",
		}),
		NewComplexHeaderLine("INFO", map[string]string{
			"ID":          "AC",
			"Number":      "A",
			"Type":        "Integer",
			"Description": "Allele count in genotypes, for each ALT allele, in the same order as listed",
		}),
		NewComplexHeaderLine("INFO", map[string]string{
			"ID":          "AF",
			"Number":      "A",
			"Type":        "Float",
			"Description": "Allele Frequency, for each ALT allele, in the same order as listed",
		}),
		NewComplexHeaderLine("INFO", map[string]string{
			"ID":          "AN",
			"Number":      "1",
			"Type":        "Integer",
			"Description": "Total number of alleles in called genotypes",
		}),
	}
}
