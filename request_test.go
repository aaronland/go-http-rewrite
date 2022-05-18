package rewrite

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func baseRequestHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		msg := req.Header.Get("X-Message")
		body := fmt.Sprintf("<html><head><title>Test</title></head><body><p>%s</p></body></html>", msg)

		rsp.Header().Set("Content-type", "text/html")
		rsp.Write([]byte(body))
	}

	h := http.HandlerFunc(fn)
	return h
}

func TestRewriteRequestHandler(t *testing.T) {

	rewrite_func := func(req *http.Request) (*http.Request, error) {

		req.Header.Set("X-Message", "hello world")
		return req, nil
	}

	request_handler := baseRequestHandler()

	request_handler = RewriteRequestHandler(request_handler, rewrite_func)

	s := &http.Server{
		Addr:    ":9664",
		Handler: request_handler,
	}

	// defer s.Close()

	go func(s *http.Server) {

		err := s.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to start server, %v", err)
		}

	}(s)

	rsp, err := http.Get("http://localhost:9664")

	if err != nil {
		t.Fatalf("Failed to GET response, %v", err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("Failed to read response, %v", err)
	}

	expected := `<html><head><title>Test</title></head><body><p>hello world</p></body></html>`

	if string(body) != expected {
		t.Fatalf("Invalid output: '%s'", string(body))
	}

}
