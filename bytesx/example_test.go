package bytesx_test

import (
	"fmt"

	"github.com/hamba/pkg/bytesx"
)

func ExampleBuffer() {
	p := bytesx.NewPool(1024) // A Buffer is returned from a pool

	buf := p.Get() // Buffer is Reset when getting it

	buf.WriteString("Hello")
	buf.Write([]byte(" World!"))

	fmt.Println(buf.String())

	p.Put(buf) // Release the buffer back to the pool

	// Output: Hello World!
}
