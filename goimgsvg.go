package goimgsvg

import (
	"encoding/base64"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// LoadMapCurrencyTld ...
func loadMapCurrencyTld() (map[string][]string, error) {

	f, err := os.Open("./assets/goimgsvg/tld_cur.tsv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	currencyTld := make(map[string][]string)

	// https://golang.org/pkg/encoding/csv/
	csvr := csv.NewReader(f)
	csvr.Comma = '\t'

	lines, err := csvr.ReadAll()
	if err != nil {
		return currencyTld, err
	}

	for _, line := range lines {
		tld := strings.TrimPrefix(line[0], ".")
		cur := strings.Split(strings.ToLower(line[1]), ",")
		for _, val := range cur {
			currencyTld[val] = append(currencyTld[val], tld)
		}

	}

	return currencyTld, nil

}

func loadMapFilenameSVGBase64() (map[string]string, error) {

	svg := make(map[string]string)
	svgPath := "./assets/goimgsvg/svg"

	files, err := ioutil.ReadDir(svgPath)
	if err != nil {
		// log.Fatal(err)
		return svg, err
	}

	for _, file := range files {
		svgFileName := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		svgContent, err := ioutil.ReadFile(svgPath + "/" + file.Name())
		if err != nil {
			// log.Fatal(err)
			return svg, err
		}
		svg[svgFileName] = string(base64.StdEncoding.EncodeToString(svgContent))
	}

	return svg, nil

}

// GoImgSVG ...
type GoImgSVG struct {
	Base64SVG   map[string]string // [iso3166]base64
	CurrencyTLD map[string][]string
}

// currencyToISO3166 ...
func (cf *GoImgSVG) currencyToISO3166(currency string) []string {

	return cf.CurrencyTLD[currency]

}

// printSVG ...
func (cf *GoImgSVG) printSVG(base64 string, alt string) string {

	img := "<div class='goimgsvg'><img src='data:image/svg+xml;base64," + base64 + "' alt='" + alt + "'/></div>"
	return img

}

// GetSVGByCurrency ...
func (cf *GoImgSVG) GetSVGByCurrency(currency string) string {

	tld := cf.currencyToISO3166(currency)
	var img string
	for _, v := range tld {
		img += cf.GetSVGByFilename(v)
	}
	return img

}

// GetSVGByFilename ... (country flags use ISO-3166 )
func (cf *GoImgSVG) GetSVGByFilename(id string) string {

	svg := cf.Base64SVG[id]
	return cf.printSVG(svg, id)

}

// NewGoImgSVG returns a CountrySVGs instance with a pre-cached flags-hashmap
// Only call once at the beginning because it's slow and therefore expensive
func NewGoImgSVG() *GoImgSVG {

	base64SVG, err := loadMapFilenameSVGBase64()
	if err != nil {
		log.Fatal(err)
	}

	currencyTld, err := loadMapCurrencyTld()
	if err != nil {
		log.Fatal(err)
	}

	cf := &GoImgSVG{
		Base64SVG:   base64SVG,
		CurrencyTLD: currencyTld,
	}

	return cf

}
