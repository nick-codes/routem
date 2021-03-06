# Routem [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/nick-codes/routem) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/nick-codes/routem/master/LICENSE) [![Build Status](http://img.shields.io/travis/nick-codes/routem.svg?style=flat-square)](https://travis-ci.org/nick-codes/routem) [![Coverage Status](https://coveralls.io/repos/nick-codes/routem/badge.svg?branch=master&service=github)](https://coveralls.io/github/nick-codes/routem?branch=master)

Yet another router for Go.

Routem focuses on integrating context with the routing stack, while
providing an easy path from those using legacy non-context aware
routes and net/http compatible handlers to routem handlers which take
a context.

It is also designed to be a completely abstract API that is bound to
an implementation via a HandlerFactory. This allows for multiple
backend implementations and experimentation with the actual
implementation.

Routem is currently a work in progress.
