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
			Parse(test.Domain, test.Contents)
		})

	}

}

func TestFetchURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	ads, err := Fetch(ctx, "https://ayan.net/ads/ads1.txt", "https://ayan.net/ads/ads2.txt")

	if err != nil {
		t.Logf("error: %s", err)
	}

	t.Logf("ads: %#v", ads)
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
