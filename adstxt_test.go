package adstxt

import (
	"context"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

	tests := []struct {
		Name     string
		Domain   string
		Contents string
	}{
		{
			Name:   "AV Club",
			Domain: "avclub.com",
			Contents: `sonobi.com, 93ee4e8333, DIRECT
amazon-adsystem.com, 3076, DIRECT
facebook.com, 743604429120938, DIRECT
rubiconproject.com, 12156, DIRECT
indexexchange.com, 183957, DIRECT
indexexchange.com, 184856, DIRECT
adtech.com, 10434, DIRECT
aolcloud.net, 10434, DIRECT
advertising.com, 10809, DIRECT
google.com, pub-9268440883448925, DIRECT
yieldmo.com, Fusion%20Media%20Group, DIRECT
yieldmo.com, 1701426062972061316, DIRECT`,
		}, {
			Name:   "Full Example #1",
			Domain: "https://www.ayan.net/",
			Contents: `# Ads.txt file for example.com:
greenadexchange.com, 12345, DIRECT, d75815a79
silverssp.com, 9675, RESELLER, f496211
blueadexchange.com, XF436, DIRECT
orangeexchange.com, 45678, RESELLER
silverssp.com, ABE679, RESELLER`,
		},
		{
			Name:   "Full Example #2",
			Domain: "https://www.ayan.net/",
			Contents: `# Ads.txt file for example.com:
greenadexchange.com, 12345, DIRECT, d75815a79
blueadexchange.com, XF436, DIRECT
contact=adops@example.com
contact=http://example.com/contact-us`,
		},

		{
			Name:   "Full Example #3",
			Domain: "https://www.ayan.net/",
			Contents: `# Ads.txt file for example.com:
greenadexchange.com, 12345, DIRECT, d75815a79
blueadexchange.com, XF436, DIRECT
subdomain=divisionone.example.com`,
		},

		{
			Name:     "Contact Variable",
			Domain:   "https://www.ayan.net/ads.txt",
			Contents: "CONTACT=ayan@ayan.net",
		}, {
			Name:   "Multiple Contacts",
			Domain: "https://www.ayan.net/ads.txt",
			Contents: `CONTACT=ayan@ayan.net
CONTACT=http://www.ayan.net`,
		}, {
			Name:   "Multiple Variables",
			Domain: "https://www.ayan.net/ads.txt",
			Contents: `CONTACT=ayan@ayan.net
CONTACT=http://www.ayan.net
SUBDOMAIN=goosgoarch.com
CONTACT=ayan@goosgoarch.com
`,
		}, {
			Name:     "First Test",
			Domain:   "https://www.ayan.net/ads.txt",
			Contents: "greenadexchange.com, XF7342, DIRECT, 5jyxf8k54",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			a, err := Parse(test.Domain, test.Contents)

			if err != nil {
				t.Fatalf("error: %v", err)
			}

			t.Logf("%#v", a)
		})

	}

}

func TestFetchURL(t *testing.T) {
	// some reandomly selected test domains
	urls := []string{
		"http://todaysgolfer.co.uk/ads.txt",
		"http://rtl.be/ads.txt",
		"http://cookingwithnonna.com/ads.txt",
		"http://dangthatsdelicious.com/ads.txt",
		"http://bebrainfit.com/ads.txt",
		"http://lifemanagerka.pl/ads.txt",
		"http://dotgolf.it/ads.txt",
		"http://tapisdedouche.com/ads.txt",
		"http://nestoria.in/ads.txt",
		"http://inallyoudo.net/ads.txt",
		"http://polishexpress.co.uk/ads.txt",
		"http://abcya.com/ads.txt",
		"http://rpgsite.net/ads.txt",
		"http://notuxedo.com/ads.txt",
		"http://kanonierzy.com/ads.txt",
		"http://sun.com/ads.txt",
		"http://avclub.com/ads.txt",
		"http://marry-xoxo.com/ads.txt",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	srcmap := make(map[string]struct{})

	for i := range urls {
		srcmap[urls[i]] = struct{}{}
	}

	ads, err := Fetch(ctx, urls...)

	if err != nil {
		t.Logf("error: %s", err)
	}

	for i := range ads {
		t.Logf("ads: %#v", ads[i])
	}

	for i := range ads {
		delete(srcmap, ads[i].Source)
	}


	for i := range a

	if len(srcmap) > 0 {
		t.Logf("We did not get a response from the following sources:")
		for k := range srcmap {
			t.Logf("> %s", k)
		}
	}
}

/*
func TestFetchMultiplex(t *testing.T) {
	url := []string{
		"https://www.google.com",
		"https://www.digg.com",
		"https://www.reddit.com",
	}

	result := make(chan string, len(url))
	wg := sync.WaitGroup{}

}
*/
