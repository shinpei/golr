Golr - Go client library for Apache Solr
====

[![Build Status](https://drone.io/github.com/shinpei/golr/status.png)](https://drone.io/github.com/shinpei/golr/latest)

Golr aims to provide you a fully accessibility to [Apache Solr](http://lucene.apache.org/solr) from Go.

## Install

```bash
go get github.com/shinpei/golr
```

## Usage

```go
import (
    "github.com/shinpei/golr"
    "runtime"
    "fmt"
)
func main () {
	con, _ := golr.Connect("localhost", 8983)
	title := "example"
	textBody := "this is an example"
	d := []Page{{
		Id:        "uniqueKey",
		Title:     title,
		Text:      textBody,
		TextCount: len(textBody),
	},
	}

	opt := &golr.SolrAddOption{
		Concurrency: runtime.NumCPU(),
	}
	msg := <-con.AddDocuments(d, opt)
	fmt.Println(string(msg[:]))
}
```


## Lisence, contact info, contribute
It's under [ASL2.0](http://www.apache.org/licenses/LICENSE-2.0). If you find bug or improvement request, please contact me through twitter, @shinpeintk. And always welcoming heartful pull request.

Cheers, :beer: :moyai:

