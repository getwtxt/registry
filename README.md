# getwtxt/registry [![Build Status](https://travis-ci.com/getwtxt/registry.svg?branch=master)](https://travis-ci.com/getwtxt/registry) [![GoDoc](https://godoc.org/github.com/getwtxt/registry?status.svg)](https://godoc.org/github.com/getwtxt/registry)
### `twtxt` Registry Library for Go

`getwtxt/registry` helps you implement `twtxt` registries in Go.
It uses no third-party dependencies whatsoever, only the standard library.

## Using the Library

You can grab a copy by issuing:

```
$ go get -u github.com/getwtxt/registry
```

If you're using Go Modules, `go get` is smart enough to snag the most recent 
tagged version. Subsequent runs of `go get -u` will only update the local copy 
of the library by minor and patch versions (0.x.x), not major versions (x.0.0).

Then, in the appropriate source file of your project, include this in your
`import` statement:

```go
import(
  "github.com/getwtxt/registry"
)
```

## Documentation

The code is commented, so feel free to browse the files themselves. 
Alternatively, the generated documentation can be found at:

[godoc.org/github.com/getwtxt/registry](https://godoc.org/github.com/getwtxt/registry)

## Contributions

All contributions are very welcome! Please feel free to submit a `PR` if you find something
that needs improvement.

## Notes

* getwtxt - parent project:
  * [github.com/getwtxt/getwtxt](https://github.com/getwtxt/getwtxt) 

* twtxt documentation: 
  * [twtxt.readthedocs.io/en/latest/](https://twtxt.readthedocs.io/en/latest/)
* twtxt registry documentation:
  * [twtxt.readthedocs.io/en/latest/user/registry.html](https://twtxt.readthedocs.io/en/latest/user/registry.html)
