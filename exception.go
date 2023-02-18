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

func Try(f func(ex Scope)) Scope {
	ex := Scope{
		exceptionLine: make(chan exception),
		catch:         make(chan exception),
	}

	go func() {
		res := <-ex.catch
		ex.exceptionLine <- res
	}()

	go func() {
		f(ex)
		close(ex.catch)
	}()

	return ex
}

func (ex Scope) Catch(f func(err error)) {
	result := <-ex.exceptionLine

	if result.err != nil {
		f(result.err)
		close(ex.exceptionLine)
	}
}

func (ex Scope) Throw(err error) {
	ex.exceptionLine <- exception{err: err}
	runtime.Goexit()
}

func (ex Scope) CatchIs(target error, f func(err error)) Scope {
	result, ok := <-ex.exceptionLine

	if errors.Is(result.err, target) {
		go func() {
			f(result.err)
			close(ex.exceptionLine)
		}()
	} else {
		if ok {
			go ex.Throw(result.err)
		}
	}

	return ex
}
