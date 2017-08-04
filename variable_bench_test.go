package variable_test

import (
	"net/url"
	"testing"
	"net/http"
	"variable"
)

// TODO: Add benchmarks

func Benchmark_populating_undefined_variables(b *testing.B) {
	input := variableReplaceTest	{
		Description: "Should handle undefined variables",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			//Host:  "www.techcrunch.com",
			Form: map[string][]string{},
		},
	}

	for i := 0; i < b.N; i++ {
		variable.PopulateRequestTemplate(input.Req, input.Variables)
	}

}

func Benchmark_variable_replacement(b *testing.B) {
	input := variableReplaceTest{
		Description: "Simple: If variable doesn't exist don't modify request.",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			//Host:  "www.techcrunch.com",
			Form: map[string][]string{},
		},
		Variables: `{"host2": "www.techcrunch.com"}`,
	}

	for i := 0; i < b.N; i++ {
		variable.PopulateRequestTemplate(input.Req, input.Variables)
	}
}

func Benchmark_double_action_replacement(b *testing.B) {
	input := variableReplaceTest{
	Description: "Bonus: Replace request body with nested variable syntax.",
		Req: &http.Request{
			Method: "POST",
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor:       1,
			ProtoMinor:       1,
			Header:           http.Header{},
			Close:            true,
			TransferEncoding: []string{"chunked"},
		},
			Body: []byte("{{key}}"),
			Variables: `{
	    "key": "body",
	    "body": "abcdef"
	  }`,
			Expected: "POST /search HTTP/1.1\r\n" +
			"Host: www.google.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Connection: close\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"Accept-Encoding: gzip\r\n\r\n" +
			chunk("abcdef") + chunk(""),
	}

	for i := 0; i < b.N; i++ {
		variable.PopulateRequestTemplate(input.Req, input.Variables)
	}
}

