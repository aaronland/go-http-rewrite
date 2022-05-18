package rewrite

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func baseAppendHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(`<html><head><title>Test</title></head><body>Hello world</body><html>`))
	}

	h := http.HandlerFunc(fn)
	return h
}

func TestAppendHTMLHandler(t *testing.T) {

	append_opts := &AppendResourcesOptions{
		JavaScript:     []string{"test.js"},
		Stylesheets:    []string{"test.css"},
		DataAttributes: map[string]string{"example": "example"},
	}

	base_handler := baseAppendHandler()

	base_handler = AppendResourcesHandler(base_handler, append_opts)

	s := &http.Server{
		Addr:    ":8080",
		Handler: base_handler,
	}

	defer s.Close()

	go func() {

		err := s.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to start server WHY, %v", err)
		}
	}()

	rsp, err := http.Get("http://localhost:8080")

	if err != nil {
		t.Fatalf("Failed to GET response, %v", err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("Failed to read response, %v", err)
	}

	expected := `<html><head><title>Test</title><script type="text/javascript" src="test.js"></script><link type="text/css" rel="stylesheet" href="test.css"/></head><body data-example="example">Hello world</body></html>`

	if string(body) != expected {
		t.Fatalf("Invalid output: '%s'", string(body))
	}

}
