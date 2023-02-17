package killerr

import (
	"errors"
	"runtime"
)

type exception struct {
	err error
}

type Scope struct {
	exceptionLine chan exception
	catch         chan exception
}

func Try(f func(h Scope)) Scope {
	handler := Scope{
		exceptionLine: make(chan exception),
		catch:         make(chan exception),
	}

	go func() {
		res := <-handler.catch
		handler.exceptionLine <- res
	}()

	go func() {
		f(handler)
		close(handler.catch)
	}()

	return handler
}

func (h Scope) Catch(f func(err error)) {
	result := <-h.exceptionLine

	if result.err != nil {
		f(result.err)
		close(h.exceptionLine)
	}
}

func (h Scope) Throw(err error) {
	h.exceptionLine <- exception{err: err}
	runtime.Goexit()
}

func (h Scope) CatchIs(target error, f func(err error)) Scope {
	result, ok := <-h.exceptionLine

	if errors.Is(result.err, target) {
		go func() {
			f(result.err)
			close(h.exceptionLine)
		}()
	} else {
		if ok {
			go h.Throw(result.err)
		}
	}

	return h
}
