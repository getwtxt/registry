# twtxt Registry Library [![Build Status](https://travis-ci.com/getwtxt/registry.svg?branch=master)](https://travis-ci.com/getwtxt/registry)

This is a library that implements a custom-rolled twtxt registry. The
specific needs of getwtxt are addressed here with respect to retrieval,
storage, and querying of user data. Third-party dependencies are kept to
a minimum. No need for `memcached`, etc.

- [x] Maintain over 90% test coverage during development process
- [x] Define types and create in-memory cache objects
- [x] Basic actions (Add, Delete users)
- [x] Queries (Users, Mentions, Tags) with time-sorted output
- [ ] twtxt.txt file scraping functions
- [ ] Bridge in-memory cache to persistent storage

### Notes

* twtxt: [twtxt.readthedocs.io/en/latest/](https://twtxt.readthedocs.io/en/latest/)
* registry documentation, including API specification: [twtxt.readthedocs.io/en/latest/user/registry.html](https://twtxt.readthedocs.io/en/latest/user/registry.html)
