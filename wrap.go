package thf

import (
	"net/http"
)

type Wrapper[In any, Out any] struct {
	fn      ioeFunc[In, Out]
	hasIn   bool
	hasOut  bool
	hasErr  bool
	decoder Decoder[In]
	encoder Encoder[Out]
}

func (w Wrapper[T, R]) Func() http.HandlerFunc {
	//default decoder
	de := w.decoder
	if de == nil {
		de = &jsonDecoder[T]{}
	}

	//default encoder
	en := w.encoder
	if en == nil {
		en = &jsonEncoder[R]{}
	}
	return w.FuncBy(de, en)
}

func (w Wrapper[In, Out]) FuncBy(decode Decoder[In], encode Encoder[Out]) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// decode
		var in In
		var out Out
		var err error
		if w.hasIn {
			in, err = decode.Decode(request)
		}

		// decode err
		if err != nil {
			//don't check hasErr always respond
			encode.Error(writer, request, err)
			return
		}

		// execute func
		out, err = w.fn(writer, request, in)
		// error check has priority
		if err != nil {
			//log.Printf("thf api [%s] call failed %s", request.RequestURI, err.Error())
			if w.hasErr {
				encode.Error(writer, request, err)
			}
			//do nothing
			return
		}

		//do nothing
		if !w.hasOut {
			return
		}
		encode.Out(writer, request, out)
	}
}
