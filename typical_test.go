package thf

import (
	"errors"
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

func TestApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) (*Bar, error) {
		t.Logf("thf api request：%s", in.Arg)
		if in.Ok {
			return &Bar{Content: testContent}, nil
		}

		return nil, errors.New(testError)
	}

	//http server
	ts := httptest.NewUnstartedServer(Wrap(api).Func())
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

func TestIEApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) error {
		t.Logf("thf api request：%s", in.Arg)
		if in.Ok {
			return nil
		}
		return errors.New(testError)
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapIE(api).Func())
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

func TestIApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) {
		t.Logf("thf api request：%s", in.Arg)
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapI(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	body := `{"arg":"` + testContent + `"}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)

	t.Logf("thf api response：%s", respBody)
}

func TestIOApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) *Bar {
		t.Logf("thf api request：%s", in.Arg)
		if in.Ok {
			return &Bar{Content: testContent}
		}
		return nil
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapIO(api).Func())
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

func TestOApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request) *Bar {
		return &Bar{Content: testContent}
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapO(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	body := `{"arg":"` + testContent + `"}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)

	t.Logf("thf api response：%s", respBody)
}

func TestOEApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request) (*Bar, error) {
		return &Bar{Content: testContent}, nil
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapOE(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	body := `{"arg":"` + testContent + `"}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)

	t.Logf("thf api response：%s", respBody)
}

func TestEApi(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request) error {
		return errors.New(testError)
	}

	//http server
	ts := httptest.NewUnstartedServer(WrapE(api).Func())
	ts.Start()
	defer ts.Close()

	//http client
	body := `{"arg":"` + testContent + `"}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)

	t.Logf("thf api response：%s", respBody)
}
