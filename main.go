package main

import (
	fmt "fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Measure struct {
	Extension string `json:"extension"`
	Count     uint64 `json:"count"`
	Lines     uint64 `json:"lines"`
	Bytes     uint64 `json:"bytes"`
}

var regexLinks = regexp.MustCompile(`<a class="js-navigation-open Link--primary" title="compiler" data-turbo-frame="repo-content-turbo-frame" href=".?">compiler</a>`)

func main() {
	http.HandleFunc("/gm", measure)
	log.Fatal(http.ListenAndServe(":6969", nil))
}

func measure(w http.ResponseWriter, r *http.Request) {
	repo, err := repo(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // TODO(Peter): Duplicated Code
		fmt.Fprintf(w, err.Error())
		return
	}

	b, err := html(repo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // TODO(Peter): Duplicated Code
		fmt.Fprintf(w, err.Error())
		return
	}

	// TODO(Peter): do the regex thing
	fmt.Println(regexLinks.FindAll(b, -1))
	fmt.Fprintf(w, string(b))
}

func html(repo string) ([]byte, error) {
	res, err := http.Get(repo)
	if err != nil {
		return nil, fmt.Errorf("no such repo: %q\n", repo)
	}

	// TODO(Peter): Handle memory leak, we should do Defer Close here, right?

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("err trying to read response: %q\n", repo)
	}

	return b, nil
}

func repo(r *http.Request) (string, error) {
	repo := r.URL.Query().Get("repo")
	if len(repo) <= 0 {
		return "", fmt.Errorf("repo is empty")
	}
	return repo, nil
}
