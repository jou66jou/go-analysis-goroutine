//隨手測試用
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

var (
	lastNum = 1000
)

func c(i int) {
	fmt.Printf("%v", i)
}

func main() {
	runtime.GOMAXPROCS(1)
	trace.Start(os.Stderr)
	defer trace.Stop()

	for i := 0; i < lastNum; i++ {
		go c(i)
	}

	fmt.Scanln()
}
