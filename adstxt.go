package adstxt

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type AdsTxt struct {
	Source   string // URL of ads.txt file this represents.
	Partner  []Buyer
	Variable map[string][]string
}

type Buyer struct {
	Domain                 string
	PublisherID            string
	AccountType            string
	CertificationAuthority string
}

//
// Attempts to split a comma separated buyer record.
//
func (a *AdsTxt) parseBuyerRecord(line string) error {
	// if we made it here, we should look for a comma separated list.
	col := strings.Split(line, ",")
	if len(col) != 4 {
		// Something is very wrong here.
		return fmt.Errorf("Could not extract buyer records.")
	}
	a.Partner = append(a.Partner,
		Buyer{
			Domain:                 strings.TrimSpace(col[0]),
			PublisherID:            strings.TrimSpace(col[1]),
			AccountType:            strings.TrimSpace(col[2]),
			CertificationAuthority: strings.TrimSpace(col[3]),
		})
	return nil
}

// only two variables are supported: CONTACT and SUBDOMAIN.  lets check for
// those.
//
// i'm not sure how to write a generic parser for the variabels because the
// spec is ambiguous when it comes to if '=' can appear within the comma
// separated buyer record.  it is possible to confuse a valid buyer with
// a variable.
//
// so instead i only look for "CONTACT=" and "SUBDOMAINS=".
//
// naughty spec!
//
func (a *AdsTxt) parseVariable(line string) error {
	for _, v := range []string{"CONTACT", "SUBDOMAIN"} {
		tok := v + "="
		if len(line) >= len(tok) && strings.ToUpper(line[:len(tok)]) == tok {
			val := line[len(tok):]
			a.Variable[v] = append(a.Variable[v], val)
			return nil
		}
	}
	return fmt.Errorf("Line does not conain a variable.")
}

func Parse(srcurl, txt string) (AdsTxt, error) {
	rc := AdsTxt{
		Source:   srcurl,
		Variable: make(map[string][]string),
	}

	if txt == "" {
		return rc, fmt.Errorf("Given an empty string.  Nothing to parse.")
	}

	// create a scanner that reads line by line.
	scanner := bufio.NewScanner(strings.NewReader(txt))

	for scanner.Scan() {
		fmt.Printf(">> %s", scanner.Text())
		line := strings.TrimSpace(scanner.Text())

		// record is blank.
		if line == "" {
			continue
		}

		// record is a comment.
		if line[0] == '#' {
			continue
		}

		if rc.parseVariable(line) == nil {
			continue
		}

		// at this point there is no need to check the result of
		// parseBuyerRecord().  if it doesn't work then we have an invalid record
		// and the loop will continue anyway.
		rc.parseBuyerRecord(line)
	}
	return rc, nil
}

type fetchresult struct {
	URL      string
	Contents string
}

func fetch(url string, rc chan<- AdsTxt, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	a, _ := Parse(url, buf.String())
	rc <- a
}

func FetchAdTxt(urls ...string) ([]AdsTxt, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("No URLs supplied; nothing to do.")
	}

	results := make(chan AdsTxt, len(urls))
	wg := sync.WaitGroup{}

	wg.Add(len(urls))
	for _, url := range urls {
		go fetch(url, results, &wg)
	}

	wg.Wait()
	close(results)

	var rc []AdsTxt
	for r := range results {
		rc = append(rc, r)
	}

	return rc, nil
}
