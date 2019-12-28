package golr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)


type ConnectorBase struct {
	url string
}
type SolrConnector struct {
	ConnectorBase
}

type ESConnector struct {
	ConnectorBase
}

type SolrAddOption struct {
	Concurrency int
	addressTemplate string
}

func (base *ConnectorBase) Connect (url string){
	base.url = url
}

type Connector interface {
	Connect (url string)
	PostUpdate(payload []byte) ([]byte, error)
	AddDocuments(container interface{}, opt *SolrAddOption) <-chan []byte
}

func NewSolrConnect(host string,port int) *SolrConnector {
	c := &SolrConnector{}
	c.Connect(fmt.Sprintf("http://%s:%d/solr/update/json", host, port))
	return c;
}

func NewESConnect ( host string, port int, collection string) *ESConnector {
	e := &ESConnector{}
	e.Connect(fmt.Sprintf("http://%s:%d/%s/_bulk", host, port, collection))
	return e
}

// Assumes it'll get arrays of some data structure
func (sc *SolrConnector) AddDocuments(container interface{}, opt *SolrAddOption) <-chan []byte {
	recvChan := make(chan []byte)

	var err error
	// todo: size constrain should be placed here
	defer func() {
		if err != nil {
			log.Printf("Error occured, uploading document failed")
		}
	}()
	go func(rC chan []byte) {
		b, err := json.Marshal(container)
		if err != nil {
			log.Println("Failed at marshaling json structure, ", err)
		}

		respB, err := sc.PostUpdate(b)
		if err != nil {
			log.Println(err)
		}
		rC <- respB
	}(recvChan)
	return recvChan
}

func (sc *SolrConnector) PostUpdate(payload []byte) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", sc.url, bytes.NewReader(payload))
	req.Header.Add("Content-type", "application/json")

	//dump, _ := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%s", dump)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Recieved %d bytes.\n", len(body))
	log.Printf("%s\n", body)
	return body, nil
}

// Assumes it'll get arrays of some data structure
func (sc *ESConnector) AddDocuments(container interface{}, opt *SolrAddOption) <-chan []byte {
	recvChan := make(chan []byte)

	var err error
	// todo: size constrain should be placed here
	defer func() {
		if err != nil {
			log.Printf("Error occured, uploading document failed")
		}
	}()
	go func(rC chan []byte) {
		b, err := json.Marshal(container)
		if err != nil {
			log.Println("Failed at marshaling json structure, ", err)
		}

		respB, err := sc.PostUpdate(b)
		if err != nil {
			log.Println(err)
		}
		rC <- respB
	}(recvChan)
	return recvChan
}

func (sc *ESConnector) PostUpdate(payload []byte) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", sc.url, bytes.NewReader(payload))
	req.Header.Add("Content-type", "application/json")

	//dump, _ := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%s", dump)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//log.Printf("Recieved %d bytes.\n", len(body))
	//log.Printf("%s\n", body)
	return body, nil
}

//TODO: Fix argument type
type JSONProcessor interface {
	Process(inputChan chan interface{}, opt *SolrAddOption, jsonReader io.Reader)
}

func (sc *SolrConnector) UploadJSONFile(path string, processor JSONProcessor, opt *SolrAddOption) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	defer jsonFile.Close()
	wg := new(sync.WaitGroup)
	inputChan := make(chan interface{})
	for i := 0; i < opt.Concurrency; i++ {
		wg.Add(1)
		go golrworker(sc, inputChan, opt, wg)
	}
	processor.Process(inputChan, opt, jsonFile)

	close(inputChan)
	wg.Wait()
}

func golrworker(c Connector, inputChan chan interface{}, opt *SolrAddOption, wg *sync.WaitGroup) {
	defer wg.Done()
	for pages := range inputChan {
		msg := <-c.AddDocuments(pages, opt)
		fmt.Println(string(msg[:]))
	}
}
