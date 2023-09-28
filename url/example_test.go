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
	fmt.Println(u)
	// output:
	// https://google.com:80/?q=hello
}

func ExampleURL_fields() {
	u, _ := url.Parse("http://google.com:80/?q=hello")
	fmt.Println(u.Port())
	fmt.Println(u.Hostname())
	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	fmt.Println(u.Path)
	// output:
	// 80
	// google.com
	// http
	// google.com:80
	// /?q=hello
}
