package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Measure struct {
	Extension string `json:"extension"`
	Count     uint64 `json:"count"`
	Lines     uint64 `json:"lines"`
	Bytes     uint64 `json:"bytes"`
}

// TODO: how to get the lines and bytes from the file
// regexLinesAndBytes = regexp.MustCompile(`<div class="Box-header js-blob-header py-2 pr-2 d-flex flex-shrink-0 flex-md-row flex-items-center">(.*?)</div>`)

var regexLinks = regexp.MustCompile(`<a class="js-navigation-open Link--primary" title="(.*?)" data-turbo-frame="repo-content-turbo-frame" href="(.*?)">(.*?)</a>`)

func main() {
	flag.Parse()

	repos := flag.Args()
	for _, r := range repos {
		html, _ := extractHtml(r)
		measure(string(html)) // TODO: ignore errors for now
	}
}

func measure(html string) (Measure, error) {
	_, err := extractValues(html)
	if err != nil {
		return Measure{}, err
	}

	// for _, d := range data {
	// 	fmt.Printf("d = %+v\n", d)
	// }

	return Measure{}, nil
}

type Data struct {
	Path     string
	Filename string
	IsDir    bool
}

func extractValues(html string) ([]Data, error) {
	var result []Data

	s := regexLinks.FindAllStringSubmatch(html, -1)
	for _, e := range s {
		v := e[len(e)-2:]

		data := Data{
			Path:     v[0],
			Filename: v[1],
			IsDir:    strings.Contains(v[0], "/tree/"),
		}

		result = append(result, data)

		url := fmt.Sprintf("https://github.com%s", data.Path)
		if data.IsDir {
			h, err := extractHtml(url)
			if err != nil {
				return []Data{}, err
			}

			data, err := extractValues(h)
			if err != nil {
				return []Data{}, err
			}

			result = append(result, data...)
		} else {
			s, err := extractHtml(url)
			if err != nil {
				return []Data{}, err
			}

			b := regexLinesAndBytes.FindAllStringSubmatch(s, -1)
			fmt.Printf("b = %+v\n", b)
		}
	}

	return result, nil
}

func extractHtml(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("no such repo: %q\n", url)
	}
	defer res.Body.Close()

	// TODO(Peter): Handle memory leak, we should do Defer Close here, right?

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("err trying to read response: %q\n", url)
	}

	return string(b), nil
}
