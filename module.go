package grpc

import (
	"go.k6.io/k6/js/modules"
	"google.golang.org/grpc/codes"
	gresolver "google.golang.org/grpc/resolver"
)

func init() {

	modules.Register("k6/x/grpc", New())
	modules.Register("k6/x/grpc_builder", NewNacosBuilder())

	gresolver.Register(&realNacosBuilder)
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the GRPC module for every VU.
	ModuleInstance struct {
		vu      modules.VU
		exports map[string]interface{}
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	mi := &ModuleInstance{
		vu:      vu,
		exports: make(map[string]interface{}),
	}

	mi.exports["Client"] = mi.NewClient
	mi.exports["Util"] = mi.NewUtil
	mi.defineConstants()
	return mi
}

// defineConstants defines the constant variables of the module.
func (mi *ModuleInstance) defineConstants() {
	rt := mi.vu.Runtime()
	mustAddCode := func(name string, code codes.Code) {
		mi.exports[name] = rt.ToValue(code)
	}

	mustAddCode("StatusOK", codes.OK)
	mustAddCode("StatusCanceled", codes.Canceled)
	mustAddCode("StatusUnknown", codes.Unknown)
	mustAddCode("StatusInvalidArgument", codes.InvalidArgument)
	mustAddCode("StatusDeadlineExceeded", codes.DeadlineExceeded)
	mustAddCode("StatusNotFound", codes.NotFound)
	mustAddCode("StatusAlreadyExists", codes.AlreadyExists)
	mustAddCode("StatusPermissionDenied", codes.PermissionDenied)
	mustAddCode("StatusResourceExhausted", codes.ResourceExhausted)
	mustAddCode("StatusFailedPrecondition", codes.FailedPrecondition)
	mustAddCode("StatusAborted", codes.Aborted)
	mustAddCode("StatusOutOfRange", codes.OutOfRange)
	mustAddCode("StatusUnimplemented", codes.Unimplemented)
	mustAddCode("StatusInternal", codes.Internal)
	mustAddCode("StatusUnavailable", codes.Unavailable)
	mustAddCode("StatusDataLoss", codes.DataLoss)
	mustAddCode("StatusUnauthenticated", codes.Unauthenticated)
}

// Exports returns the exports of the grpc module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: mi.exports,
	}
}
