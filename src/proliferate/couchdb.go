package proliferate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CouchURL parses config.json and returns couch http URL
func (node *Node) CouchURL() string {
	n := *node
	c := n.Config.Couch
	return (fmt.Sprintf("%s://%s:%s", c.Protocol, c.Host, c.Port))
}

// CouchTest Verifies CouchDB host is responsive
func (node *Node) CouchTest() bool {
	n := *node
	data := n.CouchGet("")

	if data["couchdb"] == "Welcome" {
		return true
	}

	return false
}

//func (node *Node) InitialzeDatabase() {
//	n.CouchPut("proliferate")
//}

// CouchGet returns map[string]interface of couch http.Get
func (node *Node) CouchGet(body string) map[string]interface{} {
	n := *node
	request := fmt.Sprintf("%v/%v", n.CouchURL(), body)

	response, err := http.Get(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	data := make(map[string]interface{})

	_ = json.Unmarshal([]byte(responseData), &data)

	return data
}

// CouchPut http Put body to couchdb
func (node *Node) CouchPut(body string) {
	n := *node
	client := &http.Client{}
	request, err := http.NewRequest("PUT", n.CouchURL(), strings.NewReader(body))
	request.ContentLength = int64(len(body))

	response, err := client.Do(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			n.Log(Message{
				Level: 2,
				Text:  err.Error(),
			})
		}
		//fmt.Println("The calculated length is:", len(string(contents)), "for the url:", url)
		//fmt.Println("   ", response.StatusCode)

		hdr := response.Header
		for key, value := range hdr {
			fmt.Println("   ", key, ":", value)
		}
		fmt.Println(contents)
	}
}
