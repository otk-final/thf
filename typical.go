package thf

import "net/http"

// ioeFunc declare in,out,error
type ioeFunc[T any, R any] func(http.ResponseWriter, *http.Request, T) (R, error)

// ioFunc declare in,out
type ioFunc[T any, R any] func(http.ResponseWriter, *http.Request, T) R

// ieFunc declare in,error
type ieFunc[T any] func(http.ResponseWriter, *http.Request, T) error

// iFunc declare in
type iFunc[T any] func(http.ResponseWriter, *http.Request, T)

// oeFunc declare out,error
type oeFunc[R any] func(http.ResponseWriter, *http.Request) (R, error)

// oFunc declare out
type oFunc[R any] func(http.ResponseWriter, *http.Request) R

// eFunc declare error
type eFunc func(http.ResponseWriter, *http.Request) error

type void struct{}

func wrapBuild[In any, Out any](target ioeFunc[In, Out], hasIn bool, hasOut bool, hasErr bool) Wrapper[In, Out] {
	return Wrapper[In, Out]{fn: target, hasIn: hasIn, hasOut: hasOut, hasErr: hasErr}
}

func Wrap[In any, Out any](target ioeFunc[In, Out]) Wrapper[In, Out] {
	return wrapBuild[In, Out](target, true, true, true)
}

type ioProxy[In any, Out any] struct {
	fun ioFunc[In, Out]
}

func WrapIO[In any, Out any](target ioFunc[In, Out]) Wrapper[In, Out] {
	caller := ioProxy[In, Out]{fun: target}
	return wrapBuild[In, Out](func(w http.ResponseWriter, r *http.Request, in In) (Out, error) {
		return caller.fun(w, r, in), nil
	}, true, true, false)
}

type ieProxy[In any] struct {
	fun ieFunc[In]
}

func WrapIE[In any](target ieFunc[In]) Wrapper[In, void] {
	caller := ieProxy[In]{fun: target}
	return wrapBuild[In, void](func(w http.ResponseWriter, r *http.Request, in In) (void, error) {
		return void{}, caller.fun(w, r, in)
	}, true, false, true)
}

type iProxy[In any] struct {
	fun iFunc[In]
}

func WrapI[In any](target iFunc[In]) Wrapper[In, void] {
	caller := iProxy[In]{fun: target}
	return wrapBuild[In, void](func(w http.ResponseWriter, r *http.Request, in In) (void, error) {
		caller.fun(w, r, in)
		return void{}, nil
	}, true, false, false)
}

type oeProxy[Out any] struct {
	fun oeFunc[Out]
}

func WrapOE[Out any](target oeFunc[Out]) Wrapper[void, Out] {
	caller := oeProxy[Out]{fun: target}
	return wrapBuild[void, Out](func(w http.ResponseWriter, r *http.Request, in void) (Out, error) {
		return caller.fun(w, r)
	}, false, true, true)
}

type oProxy[Out any] struct {
	fun oFunc[Out]
}

func WrapO[Out any](target oFunc[Out]) Wrapper[void, Out] {
	caller := oProxy[Out]{fun: target}
	return wrapBuild[void, Out](func(w http.ResponseWriter, r *http.Request, in void) (Out, error) {
		return caller.fun(w, r), nil
	}, false, true, false)
}

type eProxy struct {
	fun eFunc
}

func WrapE(target eFunc) Wrapper[void, void] {
	caller := eProxy{fun: target}
	return wrapBuild[void, void](func(w http.ResponseWriter, r *http.Request, in void) (void, error) {
		return void{}, caller.fun(w, r)
	}, false, false, true)
}
