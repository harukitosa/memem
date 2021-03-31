<h1 class="center">Memem</h1>

Memem is in-memory cache package.

The key is of type string and the value can be any Object.

# Get Start


```
go get github.com/harukitosa/memem
```

# sample

```go
package main

import (
	"log"

	"github.com/harukitosa/memem"
)

func main() {
	c := memem.NewCache()
	c.Set("key", "valued")
	log.Println(c.Get("key"))
}

```

# Contributing

We are waiting for the pullrequest.
