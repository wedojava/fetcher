package htmldoc

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// GetRawAndDoc can get html raw bytes and html.Node by rawurl.
func GetRawAndDoc(url *url.URL, retryTimeout time.Duration) ([]byte, *html.Node, error) {
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(url.String())
		if err == nil { // success
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, nil, err
			}
			reader := bytes.NewBuffer(raw)
			doc, err := html.Parse(reader)
			return raw, doc, nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil, nil, nil
}

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	if len(name) == 0 {
		return nil
	}
	if doc.Type == html.ElementNode {
		for _, tag := range name {
			if tag == doc.Data {
				nodes = append(nodes, doc)
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagName(c, name...)...)
	}
	return nodes
}

func ElementsByTagAndClass(doc *html.Node, tag, class string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || class == "" {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "class" && a.Val == class {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndClass(c, tag, class)...)
	}
	return nodes
}

func ElementsByTagAndId(doc *html.Node, tag, id string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || id == "" {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "id" && a.Val == id {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndId(c, tag, id)...)
	}
	return nodes
}

func ElementsNextByTag(doc *html.Node, tag string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || tag == "" {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data && doc.NextSibling != nil {
			nodes = append(nodes, doc.NextSibling)
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsNextByTag(c, tag)...)
	}
	return nodes
}

func ForEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
