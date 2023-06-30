package grpc

import (
	"fmt"
	"runtime"
	"testing"

	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/metrics"
)

const isWindows = runtime.GOOS == "windows"

// codeBlock represents an execution of a k6 script.
type codeBlock struct {
	code       string
	val        interface{}
	err        string
	windowsErr string
	asserts    func(*testing.T, *httpmultibin.HTTPMultiBin, chan metrics.SampleContainer, error)
}

type testcase struct {
	name       string
	setup      func(*httpmultibin.HTTPMultiBin)
	initString codeBlock // runs in the init context
	vuString   codeBlock // runs in the vu context
}

func TestUtil(t *testing.T) {
	u := Util{}
	for i := 0; i < 100; i++ {
		fmt.Println(u.GetNanoStr())
	}

}
