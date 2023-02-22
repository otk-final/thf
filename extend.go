package thf

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"strings"
)

var defaultRender = render.New()

func InitRender(rd *render.Render) {
	defaultRender = rd
}

type Decoder[T any] interface {
	Decode(*http.Request) (T, error)
}

type Encoder[R any] interface {
	Out(http.ResponseWriter, *http.Request, R)
	Error(http.ResponseWriter, *http.Request, error)
}

type jsonDecoder[T any] struct {
}

func (rev *jsonDecoder[T]) Decode(request *http.Request) (T, error) {
	var t T
	if strings.EqualFold(request.Method, http.MethodGet) {
		//take parameters from
		err := mapstructure.Decode(request.Form, &t)
		if err != nil {
			return t, err
		}
		return t, err
	}

	//take parameters from  body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return t, err
	}
	//json decode
	err = json.Unmarshal(body, &t)
	return t, err
}

type jsonEncoder[R any] struct {
}

func (rev *jsonEncoder[R]) Error(writer http.ResponseWriter, request *http.Request, err error) {
	_ = defaultRender.Text(writer, http.StatusInternalServerError, err.Error())
}

func (rev *jsonEncoder[R]) Out(writer http.ResponseWriter, request *http.Request, r R) {
	_ = defaultRender.JSON(writer, http.StatusOK, r)
}
