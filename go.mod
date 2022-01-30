module github.com/hedzr/go-bumper

go 1.15

//replace github.com/hedzr/log => ../src/github.com/hedzr/log

//replace github.com/hedzr/logex => ../src/github.com/hedzr/logex

//replace github.com/hedzr/cmdr => ../src/github.com/hedzr/cmdr

//replace github.com/hedzr/cmdr-addons => ../src/github.com/hedzr/cmdr-addons

//replace github.com/hedzr/pools => ../src/github.com/hedzr/pools

//replace gopkg.in/hedzr/errors.v2 => ../src/github.com/hedzr/errors

require (
	github.com/go-git/go-git/v5 v5.4.2
	github.com/hedzr/cmdr v1.10.7
	github.com/hedzr/log v1.5.9
	github.com/hedzr/logex v1.5.9
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	gopkg.in/hedzr/errors.v2 v2.1.5
)
