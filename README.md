# adstxt

Implements the ads.txt spec described at the link below:

  https://iabtechlab.com/ads-txt/

The spec is pretty simple.  This package simlpy parses the ads.txt file.

## Quick Start

```go
package main

import (
		"time"
		"github.com/ayang64/adstxt"
		"log"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background, time.Millisecond * 500)
		defer cancel()

	ads, err := adstxt.Fetch(ctx, "https://www.example.com/ads.txt", "https://www.example2.com/ads.txt")

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("ads: %#v\n", ads)
}
```

# TODO

Finish code to asynchronously fetch ads.txt files with cancellation.
