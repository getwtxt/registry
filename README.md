# getwtxt/registry [![Build Status](https://travis-ci.com/getwtxt/registry.svg?branch=master)](https://travis-ci.com/getwtxt/registry)
### twtxt Registry Library for Go

`getwtxt/registry` helps you implement twtxt registries in Go.
It uses no third-party dependencies whatsoever, only the standard library.

This library is still in early development. The readme will no longer have
this warning or the below task list once it's at a release stage.

- [x] Define types and create objects
- [x] Basic actions (Add, Delete users)
- [x] Queries (Users, Mentions, Tags) with time-sorted output
- [x] twtxt.txt file scraping functions
- [x] Find and squash bugs    `<---HERE`
- [ ] Refactor and optimize
- [ ] Triple-check concurrency safety
- [ ] Achieve 90% test coverage

### Notes

* twtxt: [twtxt.readthedocs.io/en/latest/](https://twtxt.readthedocs.io/en/latest/)
* registry documentation, including API specification: [twtxt.readthedocs.io/en/latest/user/registry.html](https://twtxt.readthedocs.io/en/latest/user/registry.html)
