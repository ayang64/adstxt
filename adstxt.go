package adstxt

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// AdsTxt stores all of the data associated with an adx.txt file.
type AdsTxt struct {
	Source   string // URL of ads.txt file this represents.
	Partner  []Buyer
	Variable map[string][]string
}

// Buyer encodes data assocated with one of the partners found in an ads.txt
// file.
type Buyer struct {
	Domain                 string
	PublisherID            string
	AccountType            string
	CertificationAuthority string
}

// Attempt to split a comma separated buyer record.
func (a *AdsTxt) parseBuyerRecord(line string) error {
	// if we made it here, we should look for a comma separated list.
	col := strings.Split(line, ",")
	if len(col) != 3 && len(col) != 4 {
		// Something is very wrong here.
		return fmt.Errorf("could not extract buyer records")
	}

	if len(col) == 3 {
		// add an empty optional certificate authority field if it wasn't specified.
		col = append(col, "")
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
	return fmt.Errorf("line does not conain a variable")
}

// Parse parses the supplied adx.txt data and returns a packed AdsTxt
// structure.
func Parse(srcurl, txt string) (AdsTxt, error) {
	rc := AdsTxt{
		Source:   srcurl,
		Variable: make(map[string][]string),
	}

	if txt == "" {
		return rc, fmt.Errorf("given an empty string; nothing to parse")
	}

	// create a scanner that reads line by line.

	for scanner := bufio.NewScanner(strings.NewReader(txt)); scanner.Scan(); {
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

func fetch(url string, rc chan<- AdsTxt) {
	resp, err := http.Get(url)

	if err != nil {
		rc <- AdsTxt{Source: url}
		return
	}

	defer resp.Body.Close()

	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	a, _ := Parse(url, buf.String())
	rc <- a
}

// Fetch downloads and parses ads.txt files at each of the supplied URLs.
func Fetch(ctx context.Context, urls ...string) ([]AdsTxt, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs supplied; nothing to do")
	}

	ads := make(chan AdsTxt)

	go func() {
		// make sure we close we our ads channel on return.
		defer close(ads)

		results := make(chan AdsTxt)

		// dispatch web queries.
		go func() {
			for _, url := range urls {
				go fetch(url, results)
			}
		}()

		// wait for respnses from each of the servers.
		for i := 0; i < len(urls); i++ {
			select {
			case r := <-results:
				ads <- r
			case <-ctx.Done():
				// deadline reached.
				return
			}
		}
	}()

	var rc []AdsTxt
	for a := range ads {
		rc = append(rc, a)
	}

	return rc, nil
}
