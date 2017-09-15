package main

/*
Parse Nmap NSE scripts.
*/

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"path/filepath"
)

var re_desc = regexp.MustCompile(`(?m)^description = "(.*)"`)
var re_desc_m = regexp.MustCompile(`(?sm)^description = \[\[\n(.*?)\]\]`)
var re_emph = regexp.MustCompile(`\*([a-zA-Z ]+)\*`)
var re_list = regexp.MustCompile(` *[*-] +`)

func clean(data string) string {
	data = re_emph.ReplaceAllString(data, "$1")
	data = strings.Replace(data, "<code>", "", -1)
	data = strings.Replace(data, "</code>", "", -1)

	return data
}

func list(data string) string {
	items := re_list.Split(data, -1)
	items[0] = wrap(items[0], false)

	for i, _ := range items[1:] {
		items[i+1] = fmt.Sprintf("\n * %s", wrap(items[i+1], true))
	}

	return strings.Join(items, "")
}

func format(data string) string {
	var paragraphs []string

	for _, p := range strings.Split(data, "\n\n") {
		paragraphs = append(paragraphs, list(p))
	}

	return strings.Join(paragraphs, "\n\n")
}

func parseNSE(data []byte) string {
	var description string

	m := re_desc.FindSubmatch(data)
	if m != nil {
		description = clean(string(m[1]))
	} else {
		m = re_desc_m.FindSubmatch(data)
		if m != nil {
			description = format(clean(string(m[1])))
		}
	}

	return description
}

func loadNSE(filename string) (scanner, error) {
	var nse scanner

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nse, err
	}

	nse.SetName(filepath.Base(filename))
	nse.SetPath(filename)
	nse.SetDescription(parseNSE(data))

	return nse, nil
}
