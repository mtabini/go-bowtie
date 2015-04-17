# Bowtie - Web middleware for Go

Bowtie is an HTTP middleware for Go. It makes heavy use of interfaces to provide a programming model that is both idiomatic and easily extensible (with, hopefully, relatively minimal overhead).

## Getting started

Standing up a basic Bowtie server only requires a few lines of code:

```
package main

import (
    "github.com/mtabini/bowtie"
    "net/http"
)

func main() {
    s := bowtie.NewServer()

    http.ListenAndServe(":8000", s)
}
```

Like all comparable systems, however, Bowtie becomes powerful once you start adding middlewares to it.

## Middlewares and contexts

Bowtie works by listening for an HTTP(S) connection and then creating an execution context that encapsulates its primary elements (that is, a request object and a response writer). It then executes a series of zero or more middlewares against this context until either some data is written to the output or an error occurs.

The context is 