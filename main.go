package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type receivers []string

func (r *receivers) String() string {
	return strings.Join(*r, ",")
}

func (r *receivers) Set(id string) error {
	*r = append(*r, id)

	return nil
}

func main() {
	appkey := flag.String("appkey", "", "appkey, the")
	secret := flag.String("secret", "", "secret, the")

	var to receivers
	flag.Var(&to, "to", "receivers' (user) id")
	flag.Parse()

	if *appkey == "" || *secret == "" {
		flag.Usage()
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	http.PostForm(
		"http://api.meiaoju.com/v1/official/send",
		url.Values{
			"key":     {*appkey},
			"sign":    {sign(*appkey, *secret, timestamp)},
			"time":    {timestamp},
			"text":    {string(body)},
			"fans_id": {strings.Join(to, ",")},
		},
	)
}

func sign(appkey, secret, nonce string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s-%s-%s", appkey, secret, nonce))))
}
