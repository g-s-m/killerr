package killerr_test

import (
	"errors"
	"testing"

	"github.com/g-s-m/killerr"

	"github.com/stretchr/testify/assert"
)

func fooError(ex killerr.Scope) {
	ex.Throw(errors.New("some error"))
}

func fooAnotherError(ex killerr.Scope) {
	ex.Throw(errors.New("another error"))
}

func fooLogicNoErr() {
}

func fooLogic(ex killerr.Scope) {
	fooLogicNoErr()
	fooError(ex)
}

func TestException(t *testing.T) {
	t.Run("raise a error", func(t *testing.T) {
		ex := killerr.Try(func(h killerr.Scope) {
			fooLogic(h)
			assert.Fail(t, "wrong execution")
		})

		ex.Catch(func(err error) {
			assert.ErrorContains(t, err, "some error")
		})
	})
	t.Run("not raise a error", func(t *testing.T) {
		ex := killerr.Try(func(h killerr.Scope) {
			fooLogicNoErr()
		})

		ex.Catch(func(err error) {
			assert.Fail(t, "function should not raise a error")
		})
	})
	t.Run("nested raise a error", func(t *testing.T) {
		ex := killerr.Try(func(h killerr.Scope) {
			fooLogicNoErr()
			ex2 := killerr.Try(func(h killerr.Scope) {
				fooAnotherError(h)
			})

			ex2.Catch(func(err error) {
				assert.ErrorContains(t, err, "another error")
				fooError(h)
			})
		})

		ex.Catch(func(err error) {
			assert.ErrorContains(t, err, "some error")
		})
	})
	t.Run("raise a error in a goroutine", func(t *testing.T) {
		ex := killerr.Try(func(h killerr.Scope) {
			fooLogicNoErr()
			ex2 := killerr.Try(func(h killerr.Scope) {
				go fooAnotherError(h)
			})

			ex2.Catch(func(err error) {
				assert.ErrorContains(t, err, "another error")
				go fooError(h)
			})
		})

		ex.Catch(func(err error) {
			assert.ErrorContains(t, err, "some error")
		})
	})
	t.Run("catch specific error", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error1, func(err error) {
			assert.ErrorIs(t, err, error1)
		})
		ex.CatchIs(error2, func(err error) {
			assert.Fail(t, "wrong exception")
		})
	})
	t.Run("catch specific error, another order", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error2)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error1, func(err error) {
			assert.Fail(t, "wrong exception")
		})
		ex.CatchIs(error2, func(err error) {
			assert.ErrorIs(t, error2, err)
		})
	})
	t.Run("rethrow a error", func(t *testing.T) {
		error1 := errors.New("error1")
		error2 := errors.New("error2")

		ex := killerr.Try(func(h killerr.Scope) {
			h.Throw(error1)
			assert.Fail(t, "wrong execution")
		})
		ex.CatchIs(error1, func(err error) {
			assert.ErrorIs(t, error1, err)
			ex.Throw(error2)
		})
		ex.CatchIs(error2, func(err error) {
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
	t.Run("throw a error, catchIs, catch all", func(t *testing.T) {
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
	// to do: find a way to handle catch and catchIs together
	// now test is dead locked because of sync waiting in Catch
	// t.Run("throw a error, catch all, catchIs", func(t *testing.T) {
	// 	error1 := errors.New("error1")

	// 	ex := killerr.Try(func(h killerr.Scope) {
	// 		h.Throw(error1)
	// 		assert.Fail(t, "wrong execution")
	// 	})
	// 	ex.Catch(func(err error) {
	// 		assert.Fail(t, "wrong execution")
	// 	})
	// 	ex.CatchIs(error1, func(err error) {
	// 		assert.ErrorIs(t, error1, err)
	// 	})
	// })
}
