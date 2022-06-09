package main

import (
	"context"
	"fmt"

	"github.com/domainr/whois"
)

func main() {
	req, err := whois.NewRequest("bcube.co.uk")
	if err != nil {
		panic(err)
	}
	w := whois.NewClient(0)
	res, err := w.FetchContext(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())
}
