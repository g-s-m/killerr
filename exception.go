package killerr

import (
	"errors"
	"runtime"
)

type exception struct {
	err error
}

// Scope allows to raise and catch errors
type Scope struct {
	exceptionLine chan exception
	catch         chan exception
}

// Try receives a function which can raise a error and returns a Scope
//
//	scope := killerr.Try(func(ex killerr.Scope) {
//		ex.Throw(errors.New("new error"))
//	})
// Function should use incoming argument to raise a error
//
// Returned Scope should be used to catch a error
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

// Catch receives a function that will be called to handle catched error
//
//	killerr.Try(func(ex killerr.Scope){
//		ex.Throw(errors.New("new error"))
//	}).Catch(f func(err error){
//		log.Error(err.Error())
//	})
func (ex Scope) Catch(f func(err error)) {
	result := <-ex.exceptionLine

	if result.err != nil {
		f(result.err)
		close(ex.exceptionLine)
	}
}

// Throw raises a error that could be catched with Catch
//
//	killerr.Try(func(ex killerr.Scope){
//		ex.Throw(errors.New("new error"))
//	}).Catch(f func(err error){
//		log.Error(err.Error())
//	})
func (ex Scope) Throw(err error) {
	ex.exceptionLine <- exception{err: err}
	runtime.Goexit()
}

// CatchIs receives a target error and a function that will be called to handle a specific error
//
//	killerr.Try(func(ex killerr.Scope){
//		ex.Throw(ErrMyError)
//	}).CatchIs(ErrMyError, f func(err error){
//		log.Error(err.Error())
//	})
//
// It is possible to use several CatchIs to handle specific error by a specific function
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
