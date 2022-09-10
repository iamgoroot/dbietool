# dbietool

[![codecov](https://codecov.io/gh/iamgoroot/dbietool/branch/main/graph/badge.svg?token=HDGXEOT8BA)](https://codecov.io/gh/iamgoroot/dbietool)

Tool for [Dbie](https://github.com/iamgoroot/dbie) generating repositories out of interface signature

Instruct interface
```//go:generate dbietool -core=Bun,Gorm -consr=func
type User interface {
    dbie.Repo[model.User]
    // ...
}
```
Run in directory containing interface:
```
dbietool -core=Bun -consr=func
```
or
```
go run github.com/iamgoroot/dbietool -core=Bun -consr=factory
```