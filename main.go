package main

import (
	"github.com/sebcat/har"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("no HAR files supplied")
	}

	var cli http.Client
	for i := 1; i < len(os.Args); i++ {
		har, err := har.LoadFile(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < len(har.Log.Entries); j++ {
			r := &har.Log.Entries[j].Request
			httpreq, err := r.Request()
			if err != nil {
				log.Fatal(err)
			}

			start := time.Now()
			resp, err := cli.Do(httpreq)
			dur := time.Since(start)
			if err != nil {
				log.Fatal(err)
			} else {
				io.Copy(ioutil.Discard, resp.Body)
				log.Println(r.Method, r.URL, len(r.PostData.Text), resp.Status, dur)
			}
		}

	}
}
