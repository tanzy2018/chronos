# Chronos
  **一个函数执行计时器**
## 使用

```go
 go get "github.com/tanzy2018/chronos"
```
## 用例
- **基础用例**
```go 
package main

import (
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	chronos.Link("save-point-1", "save-point-0")
	dur, _ := chronos.Consume("save-point-1")
	fmt.Printf("exec consumed:%dns", dur)
}
// Output
// exec consumed:104943639ns

```

- **指定区间**
```go
package main

import (
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	exeFun()
	chronos.Add("save-point-2")
	dur20, _ := chronos.ConsumeFrom("save-point-2", "save-point-0")
	fmt.Printf("0~2 exec consumed:%dns\n", dur20)
	dur21, _ := chronos.ConsumeFrom("save-point-2", "save-point-1")
	fmt.Printf("1~1 exec consumed:%dns\n", dur21)
}
// Output eg:
// 0~2 exec consumed:207272342ns
// 1~1 exec consumed:102329378ns

```

- **写入媒介**
```go
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	exeFun()
	chronos.Add("save-point-2")
	exeFun()
	chronos.Add("save-point-3")
	chronos.Link("save-point-3", "save-point-2")
	chronos.Link("save-point-2", "save-point-1")

	buf := &bytes.Buffer{}
	chronos.WriteTo(buf)
	fmt.Println(buf.String())
}

// Output
// {"consume":0,"from":"chronos-init-2W45","to":"chronos-init-2W45","unit":"ns"}
// {"consume":100071556,"from":"save-point-1","to":"save-point-2","unit":"ns"}
// {"consume":100182479,"from":"save-point-2","to":"save-point-3","unit":"ns"}
// {"consume":101845298,"from":"chronos-init-2W45","to":"save-point-1","unit":"ns"}
// {"consume":1055,"from":"chronos-init-2W45","to":"save-point-0","unit":"ns"}


```

# Chronos
  **A timer used to analysis the duration the function consume**
## Usage

```go
 go get "github.com/tanzy2018/chronos"
```
## Example
- **Base**
```go 
package main

import (
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	chronos.Link("save-point-1", "save-point-0")
	dur, _ := chronos.Consume("save-point-1")
	fmt.Printf("exec consumed:%dns", dur)
}
// Output
// exec consumed:104943639ns

```

- **Specified savepoint**
```go
package main

import (
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	exeFun()
	chronos.Add("save-point-2")
	dur20, _ := chronos.ConsumeFrom("save-point-2", "save-point-0")
	fmt.Printf("0~2 exec consumed:%dns\n", dur20)
	dur21, _ := chronos.ConsumeFrom("save-point-2", "save-point-1")
	fmt.Printf("1~1 exec consumed:%dns\n", dur21)
}
// Output eg:
// 0~2 exec consumed:207272342ns
// 1~1 exec consumed:102329378ns

```

- **Write to any io.Writer**
```go
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/tanzy2018/chronos"
)

func main() {
	exeFun := func() {
		time.Sleep(time.Millisecond * 100)
	}

	chronos.Add("save-point-0")
	exeFun()
	chronos.Add("save-point-1")
	exeFun()
	chronos.Add("save-point-2")
	exeFun()
	chronos.Add("save-point-3")
	chronos.Link("save-point-3", "save-point-2")
	chronos.Link("save-point-2", "save-point-1")

	buf := &bytes.Buffer{}
	chronos.WriteTo(buf)
	fmt.Println(buf.String())
}

// Output
// {"consume":0,"from":"chronos-init-2W45","to":"chronos-init-2W45","unit":"ns"}
// {"consume":100071556,"from":"save-point-1","to":"save-point-2","unit":"ns"}
// {"consume":100182479,"from":"save-point-2","to":"save-point-3","unit":"ns"}
// {"consume":101845298,"from":"chronos-init-2W45","to":"save-point-1","unit":"ns"}
// {"consume":1055,"from":"chronos-init-2W45","to":"save-point-0","unit":"ns"}


```



