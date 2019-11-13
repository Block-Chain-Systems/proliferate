package proliferate

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (node *Node) CouchURL() string {
	n := *node
	c := n.Config.Couch
	return (fmt.Sprintf("%s://%s:%s", c.Protocol, c.Host, c.Port))
}

func (node *Node) CouchTest() {
	n := *node

	response, err := http.Get(n.CouchURL())

	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(responseData))
}
