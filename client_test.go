package grpc

import (
	"context"
	"fmt"
	"go.k6.io/k6/js/modulestest"
	"runtime"
	"testing"

	xk6_nacos "github.com/shlsky/xk6-nacos"
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

func TestGrpcNacos(t *testing.T) {
	var nacos = xk6_nacos.New()
	err := nacos.Init("testNacos", "nacos.test.infra.ww5sawfyut0k.bitsvc.io", 8848, "nacos", "nacos", "efficiency-test")
	if err != nil {
		fmt.Println(err)
		return
	}

	var builder = NewNacosBuilder()
	builder.FillBuilder("testNacos", "")

	var grpcClient = Client{vu: &modulestest.VU{CtxField: context.Background()}}
	grpcClient.Load("CqoHCglmZWUucHJvdG8SA2ZlZSJlChNRdWVyeUZlZVJhdGVMaXN0UmVxEiYKBHBhZ2UYASABKAsyEi5mZWUuUXVlcnlMaXN0UGFnZVIEcGFnZRImCgVxdWVyeRgCIAEoCzIQLmZlZS5GZWVSYXRlSXRlbVIFcXVlcnkiUQoTUXVlcnlGZWVSYXRlTGlzdFJlcxIkCgRsaXN0GAEgAygLMhAuZmVlLkZlZVJhdGVJdGVtUgRsaXN0EhQKBWNvdW50GAIgASgFUgVjb3VudCLnAwoLRmVlUmF0ZUl0ZW0SDgoCaWQYASABKAVSAmlkEhcKB3VzZXJfaWQYAiABKAVSBnVzZXJJZBIUCgVlbWFpbBgDIAEoCVIFZW1haWwSFQoGYWZmX2lkGAQgASgFUgVhZmZJZBI2Cg1jb250cmFjdF90eXBlGAUgASgOMhEuZmVlLkNvbnRyYWN0VHlwZVIMY29udHJhY3RUeXBlEioKCWV4ZWNfdHlwZRgGIAEoDjINLmZlZS5FeGVjVHlwZVIIZXhlY1R5cGUSHwoLZGlzY291bnRfZTQYByABKAVSCmRpc2NvdW50RTQSJAoOcmF0ZV9hY3R1YWxfZTgYCCABKAVSDHJhdGVBY3R1YWxFOBIWCgZzdGF0dXMYCSABKAVSBnN0YXR1cxIlCg9taW5fZW5kX3RpbWVfZTAYCiABKAVSDG1pbkVuZFRpbWVFMBIpChFtYXhfc3RhcnRfdGltZV9lMBgLIAEoBVIObWF4U3RhcnRUaW1lRTASJwoPZGlzY291bnRfcmVhc29uGAwgASgJUg5kaXNjb3VudFJlYXNvbhISCgRub3RlGA0gASgJUgRub3RlEhcKB29wX3VzZXIYDiABKAVSBm9wVXNlchIXCgdvcF90aW1lGA8gASgFUgZvcFRpbWUiPgoNUXVlcnlMaXN0UGFnZRIXCgdsaW5lX2lkGAEgASgFUgZsaW5lSWQSFAoFbGltaXQYAiABKAVSBWxpbWl0Km4KDENvbnRyYWN0VHlwZRIOCgpVbmtvd25UeXBlEAASFAoQSW52ZXJzZVBlcnBldHVhbBABEhMKD0xpbmVhclBlcnBldHVhbBACEhIKDkludmVyc2VGdXR1cmVzEAMSDwoLQ29udHJhY3RBbGwQQiotCghFeGVjVHlwZRIKCgZVUERBVEUQABIHCgNSVU4QARIMCghST0xMQkFDSxACQg1aC2NoYXJsZXMvZmVlYgZwcm90bzMK7QIKC2hlbGxvLnByb3RvEgVoZWxsbxoJZmVlLnByb3RvIiIKDEhlbGxvUmVxdWVzdBISCgROYW1lGAEgASgJUgROYW1lIkwKDUhlbGxvUmVzcG9uc2USEgoEQ29kZRgBIAEoCVIEQ29kZRInCgREYXRhGAIgASgLMhMuaGVsbG8uUmVzcG9uc2VEYXRhUgREYXRhIjQKDFJlc3BvbnNlRGF0YRISCgROYW1lGAEgASgJUgROYW1lEhAKA0FHRRgCIAEoCVIDQUdFMooBCgVIZWxsbxI3CghTYXlIZWxsbxITLmhlbGxvLkhlbGxvUmVxdWVzdBoULmhlbGxvLkhlbGxvUmVzcG9uc2UiABJIChBRdWVyeUZlZVJhdGVMaXN0EhguZmVlLlF1ZXJ5RmVlUmF0ZUxpc3RSZXEaGC5mZWUuUXVlcnlGZWVSYXRlTGlzdFJlcyIAQg9aDWNoYXJsZXMvaGVsbG9iBnByb3RvMw==")
	_, err = grpcClient.ConnectV1("nacos:///eff-pts-agent", nil)
	if err != nil {
		fmt.Println(err)
	}
	//grpcClient.conn.Invoke(context.Background(), lib.Options{DiscardResponseBodies: null.BoolFrom(false)}, "hello.Hello/SayHello", nil, nil, nil)
}
