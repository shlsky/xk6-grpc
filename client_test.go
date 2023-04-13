package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/grpc_testing"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/fsext"
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
func TestClient(t *testing.T) {
	t.Parallel()

	type testState struct {
		*modulestest.Runtime
		httpBin *httpmultibin.HTTPMultiBin
		samples chan metrics.SampleContainer
	}
	setup := func(t *testing.T) testState {
		t.Helper()

		tb := httpmultibin.NewHTTPMultiBin(t)
		samples := make(chan metrics.SampleContainer, 1000)
		testRuntime := modulestest.NewRuntime(t)

		cwd, err := os.Getwd()
		require.NoError(t, err)
		fs := afero.NewOsFs()
		if isWindows {
			fs = fsext.NewTrimFilePathSeparatorFs(fs)
		}
		testRuntime.VU.InitEnvField.CWD = &url.URL{Path: cwd}
		testRuntime.VU.InitEnvField.FileSystems = map[string]afero.Fs{"file": fs}

		return testState{
			Runtime: testRuntime,
			httpBin: tb,
			samples: samples,
		}
	}

	assertMetricEmitted := func(
		t *testing.T,
		metricName string,
		sampleContainers []metrics.SampleContainer,
		url string,
	) {
		seenMetric := false

		for _, sampleContainer := range sampleContainers {
			for _, sample := range sampleContainer.GetSamples() {
				surl, ok := sample.Tags.Get("url")
				assert.True(t, ok)
				if surl == url {
					if sample.Metric.Name == metricName {
						seenMetric = true
					}
				}
			}
		}
		assert.True(t, seenMetric, "url %s didn't emit %s", url, metricName)
	}

	tests := []testcase{
		{
			name: "ResponseMessage",
			initString: codeBlock{
				code: `
				var client = new grpc.Client();
				client.load('CpYCCgtoZWxsby5wcm90bxIFaGVsbG8iIgoMSGVsbG9SZXF1ZXN0EhIKBE5hbWUYASABKAlSBE5hbWUiTAoNSGVsbG9SZXNwb25zZRISCgRDb2RlGAEgASgJUgRDb2RlEicKBERhdGEYAiABKAsyEy5oZWxsby5SZXNwb25zZURhdGFSBERhdGEiNAoMUmVzcG9uc2VEYXRhEhIKBE5hbWUYASABKAlSBE5hbWUSEAoDQUdFGAIgASgJUgNBR0UyQAoFSGVsbG8SNwoIU2F5SGVsbG8SEy5oZWxsby5IZWxsb1JlcXVlc3QaFC5oZWxsby5IZWxsb1Jlc3BvbnNlIgBCDloMY2hhcmxlcy5zb25nYgZwcm90bzM=');`,
			},
			setup: func(tb *httpmultibin.HTTPMultiBin) {
				tb.GRPCStub.UnaryCallFunc = func(context.Context, *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) {
					return &grpc_testing.SimpleResponse{
						OauthScope: "水",
					}, nil
				}
			},
			vuString: codeBlock{
				code: `
				client.connect("10.110.97.4:9008");
				var resp = client.invoke("hello.Hello/SayHello", {"Name": "aaa"})
				if (!resp.message || resp.message.username !== "" || resp.message.oauthScope !== "水") {
					throw new Error("unexpected response message: " + JSON.stringify(resp.message))
				}`,
				asserts: func(t *testing.T, rb *httpmultibin.HTTPMultiBin, samples chan metrics.SampleContainer, _ error) {
					samplesBuf := metrics.GetBufferedSamples(samples)
					assertMetricEmitted(t, metrics.GRPCReqDurationName, samplesBuf, rb.Replacer.Replace("GRPCBIN_ADDR/grpc.testing.TestService/UnaryCall"))
				},
			},
		},
	}

	assertResponse := func(t *testing.T, cb codeBlock, err error, val goja.Value, ts testState) {
		if isWindows && cb.windowsErr != "" && err != nil {
			err = errors.New(strings.ReplaceAll(err.Error(), cb.windowsErr, cb.err))
		}
		if cb.err == "" {
			assert.NoError(t, err)
		} else {
			require.Error(t, err)
			assert.Contains(t, err.Error(), cb.err)
		}
		if cb.val != nil {
			require.NotNil(t, val)
			assert.Equal(t, cb.val, val.Export())
		}
		if cb.asserts != nil {
			cb.asserts(t, ts.httpBin, ts.samples, err)
		}
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := setup(t)

			m, ok := New().NewModuleInstance(ts.VU).(*ModuleInstance)
			require.True(t, ok)
			require.NoError(t, ts.VU.Runtime().Set("grpc", m.Exports().Named))

			// setup necessary environment if needed by a test
			if tt.setup != nil {
				tt.setup(ts.httpBin)
			}

			replace := func(code string) (goja.Value, error) {
				return ts.VU.Runtime().RunString(ts.httpBin.Replacer.Replace(code))
			}

			val, err := replace(tt.initString.code)
			assertResponse(t, tt.initString, err, val, ts)

			registry := metrics.NewRegistry()
			root, err := lib.NewGroup("", nil)
			require.NoError(t, err)

			state := &lib.State{
				Group:     root,
				Dialer:    ts.httpBin.Dialer,
				TLSConfig: ts.httpBin.TLSClientConfig,
				Samples:   ts.samples,
				Options: lib.Options{
					SystemTags: metrics.NewSystemTagSet(
						metrics.TagName,
						metrics.TagURL,
					),
					UserAgent: null.StringFrom("k6-test"),
				},
				BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
				Tags:           lib.NewVUStateTags(registry.RootTagSet()),
			}
			ts.MoveToVUContext(state)
			val, err = replace(tt.vuString.code)
			assertResponse(t, tt.vuString, err, val, ts)
		})
	}
}
