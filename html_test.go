package rewrite

import (
	_ "fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"testing"
)

func baseRewriteHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(`<html><head><title>Test</title></head><body><p>hello world</p></body><html>`))
	}

	h := http.HandlerFunc(fn)
	return h
}

func TestRewriteHTMLHandler(t *testing.T) {

	base_handler := baseRewriteHandler()

	var rewrite_func RewriteHTMLFunc

	rewrite_func = func(n *html.Node, wr io.Writer) {

		if n.Type == html.ElementNode && n.Data == "p" {
			n.FirstChild.Data = "HELLO WORLD"
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rewrite_func(c, wr)
		}

	}

	base_handler = RewriteHTMLHandler(base_handler, rewrite_func)

	s := &http.Server{
		Addr:    ":8181",
		Handler: base_handler,
	}

	defer s.Close()

	go func(s *http.Server) {

		err := s.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to start server, %v", err)
		}

	}(s)

	rsp, err := http.Get("http://localhost:8181")

	if err != nil {
		t.Fatalf("Failed to GET response, %v", err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("Failed to read response, %v", err)
	}

	expected := `<html><head><title>Test</title></head><body><p>HELLO WORLD</p></body></html>`

	if string(body) != expected {
		t.Fatalf("Invalid output: '%s'", string(body))
	}

}
