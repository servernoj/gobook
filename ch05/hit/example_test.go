package hit_test

import (
	"context"
	"fmt"
	"log"

	"github.com/servernoj/gobook/ch05/hit"
)

func ExampleDo() {
	stat, err := hit.Do(
		context.Background(),
		"https://google.com",
		5,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stat.Count)
	// output:
	// 5
}
