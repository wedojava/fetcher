package htmldoc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// ExtractLinks makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func ExtractLinks(weburl string) ([]string, error) {
	u, err := url.Parse(weburl)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(weburl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", weburl, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", weburl, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		// TODO: compress layers
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				// append only the target website
				if strings.HasPrefix(a.Val, "http") && strings.Contains(a.Val, u.Hostname()) {
					links = append(links, link.String())
				} else if strings.HasPrefix(a.Val, "/") {
					links = append(links, link.String())
				}

			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	return links, nil
}

func ElementsNext(doc *html.Node) []*html.Node {
	nodes := []*html.Node{}
	if doc == nil {
		return nil
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil {
			nodes = append(nodes, n)
		}
	}
	ForEachNode(doc, visitNode, nil)
	return nodes
}

func ElementsRmByTag(doc *html.Node, name ...string) {
	if len(name) == 0 || doc == nil {
		return
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode {
			for _, tag := range name {
				if tag == n.NextSibling.Data {
					n.Parent.RemoveChild(n.NextSibling)
				}
			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	rmFirstTag := func(n *html.Node) {
		if n.FirstChild != nil && n.FirstChild.Type == html.ElementNode {
			for _, tag := range name {
				if tag == n.FirstChild.Data {
					n.RemoveChild(n.FirstChild)
				}
			}
		}
	}
	ForEachNode(doc, rmFirstTag, nil)
}

func ElementsByTag(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	if len(name) == 0 || doc == nil {
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
		nodes = append(nodes, ElementsByTag(c, name...)...)
	}
	return nodes
}

func ElementsByTagAndClass(doc *html.Node, tag, class string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || class == "" || doc == nil {
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

func ElementsByTagAndClass2(raw []byte, tag, class string) []byte {
	if raw == nil || tag == "" || class == "" {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	var b bytes.Buffer
	for {
		tt := z.Next()
		t := z.Token()
		if err := z.Err(); err == io.EOF {
			break
		}
		switch tt {
		case html.StartTagToken:
			if tag == t.Data {
				for _, a := range t.Attr {
					if a.Key == "class" && a.Val == class {
						b.Write(z.Raw())
					}
				}
			}
		}
	}
	return b.Bytes()
}

func ElementsByTagAndId(doc *html.Node, tag, id string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || tag == "" || id == "" {
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

func ElementsByTagAndId2(raw []byte, tag, id string) []byte {
	if raw == nil || tag == "" || id == "" {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	for {
		tt := z.Next()
		t := z.Token()
		if err := z.Err(); err != nil && err == io.EOF {
			break
		}
		switch tt {
		case html.StartTagToken:
			if tag == t.Data {
				for _, a := range t.Attr {
					if a.Key == "id" && a.Val == id {
						return z.Buffered()
					}
				}
			}
		}
	}
	return nil
}
func ElementsByTagAndType(doc *html.Node, tag, attrType string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || attrType == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "type" && a.Val == attrType {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndType(c, tag, attrType)...)
	}
	return nodes
}

func ElementsNextByTag(doc *html.Node, tag string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || doc == nil {
		return nil
	}
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

func MetasByName(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || name == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if doc.Data == "meta" {
			for _, a := range doc.Attr {
				if a.Key == "name" {
					for _, v := range name {
						if v == a.Val {
							nodes = append(nodes, doc)
						}
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, MetasByName(c, name...)...)
	}
	return nodes
}

func MetasByProperty(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || name == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if doc.Data == "meta" {
			for _, a := range doc.Attr {
				if a.Key == "property" {
					for _, v := range name {
						if v == a.Val {
							nodes = append(nodes, doc)
						}
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, MetasByProperty(c, name...)...)
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
