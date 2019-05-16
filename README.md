# getwtxt/registry [![Build Status](https://travis-ci.com/getwtxt/registry.svg?branch=master)](https://travis-ci.com/getwtxt/registry)
### twtxt Registry Library for Go

`getwtxt/registry` helps you implement a custom-rolled twtxt registry in Go.
It uses no third-party dependencies whatsoever, only the standard library.

- [x] Maintain over 90% test coverage during development process
- [x] Define types and create in-memory cache objects
- [x] Basic actions (Add, Delete users)
- [x] Queries (Users, Mentions, Tags) with time-sorted output
- [ ] twtxt.txt file scraping functions
- [ ] Refactor and optimize
- [ ] Triple-check concurrency safety

### Notes

* twtxt: [twtxt.readthedocs.io/en/latest/](https://twtxt.readthedocs.io/en/latest/)
* registry documentation, including API specification: [twtxt.readthedocs.io/en/latest/user/registry.html](https://twtxt.readthedocs.io/en/latest/user/registry.html)
