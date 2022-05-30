module github.com/iamgoroot/dbietool

go 1.18

replace (
	github.com/iamgoroot/dbie => /home/user/GolandProjects/dbie
	github.com/iamgoroot/merge => /home/user/GolandProjects/merge
)

require (
	github.com/iamgoroot/merge v0.0.0-00010101000000-000000000000
	github.com/iancoleman/strcase v0.2.0
	golang.org/x/tools v0.1.10
)

require (
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
