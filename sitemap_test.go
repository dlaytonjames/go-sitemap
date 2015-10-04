package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	server := server()
	defer server.Close()

	data, err := Get(server.URL + "/sitemap.xml")

	if len(data.URL) != 13 {
		t.Error("Get() should return Sitemap.Url(13 length)")
	}

	if err != nil {
		t.Error("Get() should not has error")
	}
}

func TestGetRecivedInvalidSitemapURL(t *testing.T) {
	server := server()
	defer server.Close()

	_, err := Get(server.URL + "/emptymap.xml")

	if err == nil {
		t.Error("Get() should return error")
	}
}

func TestGetRecivedSitemapIndexURL(t *testing.T) {
	server := server()
	defer server.Close()

	SetInterval(time.Nanosecond)
	data, err := Get(server.URL + "/sitemapindex.xml")

	if len(data.URL) != 39 {
		t.Error("Get() should return Sitemap.Url(39 length)")
	}

	if err != nil {
		t.Error("Get() should not has error")
	}
}

func TestParse(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/sitemap.xml")
	sitemap, _ := Parse(data)

	if len(sitemap.URL) != 13 {
		t.Error("Parse() should return Sitemap.URL(13 length)")
	}
}

func TestParseIndex(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")
	index, _ := ParseIndex(data)

	if len(index.Sitemap) != 3 {
		t.Error("ParseIndex() should return Index.Sitemap(3 length)")
	}
}

func TestSetInterval(t *testing.T) {
	newInterval := 3 * time.Second
	SetInterval(newInterval)

	if interval != newInterval {
		t.Error("interval should be time.Minute")
	}

	if interval == time.Second {
		t.Error("interval should not be Default(time.Second)")
	}
}

func TestSetFetch(t *testing.T) {
	f := func(url string) ([]byte, error) {
		var err error
		return []byte(url), err
	}

	SetFetch(f)

	url := "http://example.com"
	data, _ := fetch(url)

	if string(data) != url {
		t.Error("fetch(url) should return " + url)
	}
}

func BenchmarkGetSitemap(b *testing.B) {
	server := server()
	defer server.Close()
	Get(server.URL + "/sitemap.xml")
}

func BenchmarkGetSitemapIndex(b *testing.B) {
	server := server()
	defer server.Close()
	Get(server.URL + "/sitemapindex.xml")
}

func BenchmarkParseSitemap(b *testing.B) {
	data, _ := ioutil.ReadFile("./testdata/sitemap.xml")
	Parse(data)
}

func BenchmarkParseSitemapIndex(b *testing.B) {
	data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")
	ParseIndex(data)
}

func server() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "" {
			// index page is always not found
			http.NotFound(w, r)
		}

		res, err := ioutil.ReadFile("./testdata" + r.RequestURI)
		if err != nil {
			http.NotFound(w, r)
		}
		str := strings.Replace(string(res), "HOST", r.Host, -1)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, str)
	}))

	return server
}
