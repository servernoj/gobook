package url_test

import (
	"fmt"
	"log"

	"github.com/servernoj/gobook/url"
)

func ExampleURL() {
	u, err := url.Parse("http://google.com:80/?q=hello")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	fmt.Println(u.Port())
	fmt.Println(u.Hostname())
	fmt.Println(u)
	// output:
	// 80
	// google.com
	// https://google.com:80/?q=hello
}
