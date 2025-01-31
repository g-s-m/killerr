# killerr [![GoDoc](https://pkg.go.dev/badge/g-s-m/killerr)](https://pkg.go.dev/github.com/g-s-m/killerr) [![Build Status](https://github.com/g-s-m/killerr/actions/workflows/go.yml/badge.svg)](https://github.com/g-s-m/killerr/actions/workflows/go.yml) [![Coverage Status](https://codecov.io/gh/g-s-m/killerr/branch/main/graph/badge.svg)](https://codecov.io/gh/g-s-m/killerr)

Kill "if err"(killerr) package is a simple implementation of exceptions for go language.

## How to use it
Add "github.com/g-s-m/killerr" package to your project.
```bash
go get -v github.com/g-s-m/killerr
```

## Description
There is a method
```go
Try(f func(h Scope)) Scope
```
in the package. It is supposed to use in that place where you are expecting an exception. `Scope` is a structure contained 3 methods:
```go
Catch(f func(err error))
CatchIs(target error, f func(err error))
Throw(err error)
```
`Try` receives your function that implements some logic and returns a `Scope`. That function must receive a `Scope` as an argument to be able to raise an exception. Pass this `Scope` to all calles may to throw an exception(like context.Context). In all the pieces of codes where a concrete error is generated you are able to call `Throw` to raise an exception.
Call `Catch` to catch an exception. Call `CatchIs` to receive a concrete error.

Example:
```go

func MyRepo1(ctx context.Context, ex killerr.Scope) {
  record, err := db.Insert()
  if err != nil {
    ex.Throw(fmt.Errorf("repo2 err: %w", err))
  }
}

func MyRepo2(ctx context.Context, ex killerr.Scope) {
  record, err := db.Find()
  if err != nil {
    ex.Throw(fmt.Errorf("repo2 err: %w", err))
  }
}

func MyInternalPackageLogic1(ctx context.Context, ex killerr.Scope) {
  MyRepo2(ctx, ex)
}

func MyInternalPackageLogic2(ctx context.Context, ex killerr.Scope)  {
  MyRepo1(ctx, ex)
  MyRepo2(ctx, ex)
  if somethingBadHappen {
    ex.Throw(ErrInternalPackage)
  }
}

func MyAppLogic(ctx context.Context, ex killerr.Scope) {
  MyInternalPackage1(ctx, ex)
  MyInternalPackage2(ctx, ex)
}

func RunMyApp() {
  ctx := context.Background()

  killerr.Try(func(ex killerr.Scope) {
    MyAppLogic(ctx, ex)
  }).CatchIs(ErrInternalPackage, func(err error){
    log.Error("internal error: %s", err.Error())
  }).CatchIs(ErrRepo, func(err error){
    log.Error("repo error: %s", err.Error())
  }).Catch(func(err error) {
    log.Error("another error: %s", err.Error())
  })
}
```

## Restrictions
* If you doesn't place catch block, you will lose your error
* `Try`, `Catch` and `Throw` won't work in different goroutines, like try-catch in other languages works in the same thread, they would work properly only in the same goroutine.
The following code will not work correctly:
```go
// do not do this:
killerr.Try(func(ex killerr.Scope) {
  go func() {
    ex.Throw()
  }()
}).Catch(func(err error){
  ...
})
```
* Order of `Catch` and `CatchIs` matters. For example if function throws `ErrInternal1` you may call `CatchIs(ErrInternal, ...)` first and then `Catch(...)` to catch ErrInternal separately of other errors. Always place `CatchIs` before `Catch` block.
See example.
```go
killerr.Try(func(ex killerr.Scope) {
    ex.Throw(ErrInternalPackage)
  }).CatchIs(ErrInternalPackage, func(err error){
    //this block will catch the error
    log.Error("internal error: %s", err.Error())
  }).Catch(func(err error) {
    log.Error("another error: %s", err.Error())
  })

  killerr.Try(func(ex killerr.Scope) {
    ex.Throw(ErrInternalPackage)
  }).Catch(func(err error) {
    //this block will catch the error
    log.Error("another error: %s", err.Error())
  }).CatchIs(ErrInternalPackage, func(err error){
    log.Error("internal error: %s", err.Error())
  })
```

-------------------------------------------------------------------------------
Released under the [LICENSE.txt](https://github.com/g-s-m/killerr/blob/main/LICENSE.txt).
