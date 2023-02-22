package resp

import (
	"errors"
	"github.com/otk-final/thf"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testContent = "this is test message"

var testError = "this is test error"

type Foo struct {
	Ok  bool   `json:"ok"`
	Arg string `json:"arg"`
}

type Bar struct {
	Content string `json:"content"`
}

type Baz struct {
	Content string `json:"content"`
}

func TestEntryApi(t *testing.T) {
	InitRespEntryCode("1000", "2000")

	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) *Entry[*Bar] {
		if in.Ok {
			return NewEntry(&Bar{Content: testContent})
		}
		return NewError[*Bar](errors.New(testError))
	}

	//http server
	ts := httptest.NewUnstartedServer(thf.WrapIO(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	bodyContents := []string{
		`{"arg":"` + testContent + `","ok":true}`,
		`{"arg":"` + testContent + `","ok":false}`,
	}

	for _, body := range bodyContents {
		resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
		respBody, _ := ioutil.ReadAll(resp.Body)
		t.Logf("thf api response：%d %s", resp.StatusCode, respBody)
	}
}

func TestAnyApi(t *testing.T) {
	InitRespEntryCode("1000", "2000")

	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) *Entry[any] {
		if !in.Ok {
			return NewError[any](errors.New(testError))
		}
		if in.Arg == "bar" {
			return NewAny(&Bar{Content: "bar"})
		} else {
			return NewAny(&Baz{Content: "baz"})
		}
	}

	//http server
	ts := httptest.NewUnstartedServer(thf.WrapIO(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	bodyContents := []string{
		`{"arg":"bar","ok":true}`,
		`{"arg":"bar","ok":false}`,
		`{"arg":"baz","ok":true}`,
		`{"arg":"baz","ok":false}`,
	}

	for _, body := range bodyContents {
		resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
		respBody, _ := ioutil.ReadAll(resp.Body)
		t.Logf("thf api response：%d %s", resp.StatusCode, respBody)
	}
}
