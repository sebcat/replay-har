// replay-har
// ==========
//
// Application demonstrating the use of [github.com/sebcat/har](https://github.com/sebcat/har)
//
// Loads a [HAR](http://en.wikipedia.org/wiki/.har) from file and retransmits
// the HTTP requests from it sequentially, printing status to stderr. Dies if
// a request results in an error.
//
// Please note that the requests in the HAR may contain stateful
// information (cookies, session tokens in the URL/request body, &c) that may
// yield different responses when retransmitted.
//
// Example use from shell, including [GOPATH configuration](https://golang.org/doc/code.html#GOPATH):
//
// ```
// mkdir $HOME/go-code
// export GOPATH=$HOME/go-code PATH=$GOPATH/bin:$PATH
// go get github.com/sebcat/replay-har
// replay-har <har-file>
//
// ```
//
// You can save a HAR file from chrome:
// Developer Tools (Ctrl-Shift-J, Command-J, &c)-> Network -> (visit a page) -> Right click on view -> Copy all as HAR or Save as HAR with Content

//
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
				resp.Body.Close()
				log.Println(r.Method, r.URL, len(r.PostData.Text), resp.Status, dur)
			}
		}

	}
}
