package killerr_test

import (
	"errors"
	"testing"

	"github.com/g-s-m/killerr"

	"github.com/stretchr/testify/assert"
)

func emptyFunc() {}

func TestException(t *testing.T) {
	t.Run("raise a error", func(t *testing.T) {
		killerr.Try(func(h killerr.Scope) {
			h.Throw(errors.New("some error"))
			assert.Fail(t, "wrong execution")
		}).Catch(func(err error) {
			assert.ErrorContains(t, err, "some error")
		})
	})
	t.Run("not raise a error", func(t *testing.T) {
		killerr.Try(func(h killerr.Scope) {
			emptyFunc()
		}).Catch(func(err error) {
			assert.Fail(t, "function should not raise a error")
		})
	})
	t.Run("nested raise a error", func(t *testing.T) {
		killerr.Try(func(h killerr.Scope) {
			emptyFunc()
			killerr.Try(func(hh killerr.Scope) {
				hh.Throw(errors.New("another error"))
			}).Catch(func(err error) {
				assert.ErrorContains(t, err, "another error")
				h.Throw(errors.New("some error"))
			})
		}).Catch(func(err error) {
			assert.ErrorContains(t, err, "some error")
		})
	})
	t.Run("catch specific error", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		}).CatchIs(error1, func(err error) {
			assert.ErrorIs(t, err, error1)
		}).CatchIs(error2, func(err error) {
			assert.Fail(t, "wrong exception")
		})
	})
	t.Run("catch specific error, another order", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		killerr.Try(func(h killerr.Scope) {
			h.Throw(error2)
			assert.Fail(t, "wrong execution")
		}).CatchIs(error1, func(err error) {
			assert.Fail(t, "wrong exception")
		}).CatchIs(error2, func(err error) {
			assert.ErrorIs(t, error2, err)
		})
	})
	t.Run("rethrow a error", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		scope := killerr.Scope{}
		killerr.Try(func(h killerr.Scope) {
			scope = h
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		}).CatchIs(error1, func(err error) {
			assert.ErrorIs(t, error1, err)
			scope.Throw(error2)
		}).CatchIs(error2, func(err error) {
			assert.ErrorIs(t, error2, err)
		})
	})
	t.Run("rethrow a error, another catch order", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error2, func(err error) {
			assert.ErrorIs(t, error2, err)
		})
		ex.CatchIs(error1, func(err error) {
			assert.ErrorIs(t, error1, err)
			ex.Throw(error2)
		})
	})
	t.Run("rethrow a error, catch all", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")
		error3 := errors.New("error3")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error2, func(err error) {
			assert.ErrorIs(t, error2, err)
			ex.Throw(error3)
		})
		ex.CatchIs(error1, func(err error) {
			assert.ErrorIs(t, error1, err)
			ex.Throw(error2)
		})
		ex.Catch(func(err error) {
			assert.ErrorIs(t, error2, err)
		})
	})
	t.Run("throw-catchIs-catch", func(t *testing.T) {
		error1 := errors.New("error1")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error1, func(err error) {
			assert.ErrorIs(t, error1, err)
		})
		ex.Catch(func(err error) {
			assert.Fail(t, "wrong execution")
		})
	})
	t.Run("throw-catch-catchIs", func(t *testing.T) {
		error1 := errors.New("error1")
		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.Catch(func(err error) {
			assert.ErrorIs(t, error1, err)
		})
		ex.CatchIs(error1, func(err error) {
			assert.Fail(t, "wrong execution")
		})
	})
}
