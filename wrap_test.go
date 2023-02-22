package thf

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FooDecode struct {
}

func (f FooDecode) Decode(request *http.Request) (*Foo, error) {
	return &Foo{}, nil
}

type BarEncode struct {
}

func (b BarEncode) Out(writer http.ResponseWriter, request *http.Request, out *Bar) {
	log.Printf("out %v", out)
}
func (b BarEncode) Error(writer http.ResponseWriter, request *http.Request, err error) {
	log.Printf("error %s", err.Error())
}

var fooDecoder = &FooDecode{}
var barEncoder = &BarEncode{}

func TestFuncBy(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) (*Bar, error) {
		return nil, nil
	}

	//http server and custom decode and encode
	ts := httptest.NewUnstartedServer(Wrap(api).FuncBy(fooDecoder, barEncoder))
	ts.Start()
	defer ts.Close()

	body := `{"arg":"` + testContent + `","ok":true}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)
	t.Logf("thf api response：%d %s", resp.StatusCode, respBody)
}

func TestFuncByDecode(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request, in *Foo) {
	}

	//http server and custom decode and encode
	ts := httptest.NewUnstartedServer(WrapI(api).FuncBy(fooDecoder, nil))
	ts.Start()
	defer ts.Close()

	body := `{"arg":"` + testContent + `","ok":true}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)
	t.Logf("thf api response：%d %s", resp.StatusCode, respBody)
}

func TestFuncByEncode(t *testing.T) {

	//http.Handler upgrade
	api := func(writer http.ResponseWriter, request *http.Request) (*Bar, error) {
		return nil, errors.New("error")
	}

	//http server and custom decode and encode
	ts := httptest.NewUnstartedServer(WrapOE(api).FuncBy(nil, barEncoder))
	ts.Start()
	defer ts.Close()

	body := `{"arg":"` + testContent + `","ok":true}`
	resp, _ := ts.Client().Post(ts.URL, "application/json", strings.NewReader(body))
	respBody, _ := ioutil.ReadAll(resp.Body)
	t.Logf("thf api response：%d %s", resp.StatusCode, respBody)
}
