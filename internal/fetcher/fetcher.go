package fetcher

// Posts is the fetcher can return many results.
type Posts map[string]*ThePost

type ThePost struct {
	Site   string
	Domain string
	URL    string
	Title  string
	Body   string
	Date   string
}

type Paragraph struct {
	Type    string
	Content string
}
