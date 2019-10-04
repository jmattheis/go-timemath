# jmattheis/go-timemath [![Build Status][travis-badge]][travis] [![codecov][codecov-badge]][codecov]

This package contains a parser for relative times like `now-1d` or `now/d` (start of day) similar to the
ranges in [Grafana](https://grafana.com/docs/reference/timerange/). It also has convenient helper functions for time math.

## Usage

Download the package with:
```bash
$ go get github.com/jmattheis/go-timemath
```

### Math

```go
package main

import (
	"fmt"
	"time"

	"github.com/jmattheis/go-timemath"
)

func main() {
	now, _ := time.ParseInLocation(time.RFC3339, "2019-05-12T15:55:23Z", time.UTC)
	now = timemath.Day.Add(now, 1)
	now = timemath.Hour.Subtract(now, 5)
	now = timemath.Minute.StartOf(now, time.Monday)
	fmt.Println(now.Format(time.RFC3339)) // 2019-05-13T10:55:00Z
}

```

### Parsing

```go
package main

import (
	"fmt"
	"time"

	"github.com/jmattheis/go-timemath"
)

func main() {
	now, _ := time.ParseInLocation(time.RFC3339, "2019-05-12T15:55:23Z", time.UTC)
	parsed, _ := timemath.Parse(now, "now+1d-5h/m", true, time.Monday)
	//                                         ^^- start of minute
	//                                      ^^^- subtract five hours
	//                                   ^^^- add one day
	//                                ^^^ now is the given parameter
	fmt.Println(parsed.Format(time.RFC3339)) // 2019-05-13T10:55:00Z
}
```

### Units

| Operator | full name |
| -------- | --------- |
| s        | second    |
| m        | minute    |
| h        | hour      |
| d        | day       |
| w        | week      |
| M        | month     |
| y        | year      |

 [travis-badge]: https://travis-ci.org/jmattheis/go-timemath.svg?branch=master
 [travis]: https://travis-ci.org/jmattheis/go-timemath
 [codecov-badge]: https://codecov.io/gh/jmattheis/go-timemath/branch/master/graph/badge.svg
 [codecov]: https://codecov.io/gh/jmattheis/go-timemath
 