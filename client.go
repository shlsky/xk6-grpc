package grpc

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc/codes"
	"io"
	"strings"
	"time"

	"github.com/shlsky/xk6-grpc/xgrpc_conn"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/metrics"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Client represents a gRPC client that can be used to make RPC requests
type Client struct {
	mds  map[string]protoreflect.MethodDescriptor
	conn *xgrpc_conn.Conn
	vu   modules.VU
	addr string
}

type Util struct {
}

func init() {

	modules.Register("k6/x/grpc", New())
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

// NewClient is the JS constructor for the grpc Client.
func (mi *ModuleInstance) NewClient(call goja.ConstructorCall) *goja.Object {
	rt := mi.vu.Runtime()
	return rt.ToValue(&Client{vu: mi.vu}).ToObject(rt)
}

// NewUtil is the JS constructor for the grpc Util.
func (mi *ModuleInstance) NewUtil(call goja.ConstructorCall) *goja.Object {
	rt := mi.vu.Runtime()
	return rt.ToValue(&Util{}).ToObject(rt)
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

// Getnano 获取纳秒
func (c *Util) Getnano() int64 {
	return time.Now().UnixNano()
}

// Getmicro 获取微妙
func (c *Util) Getmicro() int64 {
	return time.Now().UnixMicro()
}

// Getnano 获取纳秒
func (c *Util) GetnanoStr() string {
	return string(time.Now().UnixNano())
}

// Getmicro 获取微妙
func (c *Util) GetmicroStr() int64 {
	return time.Now().UnixMicro()
}

// Load will parse the given proto files and make the file descriptors available to request.
func (c *Client) Load(descriptor string) ([]MethodInfo, error) {
	if c.vu.State() != nil {
		return nil, errors.New("load must be called in the init context")
	}

	initEnv := c.vu.InitEnv()
	if initEnv == nil {
		return nil, errors.New("missing init environment")
	}

	bb, err := base64.StdEncoding.DecodeString(descriptor)
	if err != nil {
		return nil, err
	}

	fds, err := LoadFileDescriptorSet(bb)
	if err != nil {
		return nil, err
	}

	fdset := &descriptorpb.FileDescriptorSet{}

	seen := make(map[string]struct{})
	for _, fd := range fds {
		fdset.File = append(fdset.File, walkFileDescriptors(seen, fd)...)
	}
	return c.convertToMethodInfo(fdset)
}

// LoadFileDescriptorSet 加载pb.bin文件二进制
func LoadFileDescriptorSet(b []byte) (map[string]*desc.FileDescriptor, error) {
	var fds descriptor.FileDescriptorSet

	if err := proto.Unmarshal(b, &fds); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorsFromSet(&fds)
}

// LoadProtoset will parse the given protoset file (serialized FileDescriptorSet) and make the file
// descriptors available to request.
func (c *Client) LoadProtoset(protosetPath string) ([]MethodInfo, error) {
	if c.vu.State() != nil {
		return nil, errors.New("load must be called in the init context")
	}

	initEnv := c.vu.InitEnv()
	if initEnv == nil {
		return nil, errors.New("missing init environment")
	}

	absFilePath := initEnv.GetAbsFilePath(protosetPath)
	fdsetFile, err := initEnv.FileSystems["file"].Open(absFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open protoset: %w", err)
	}

	defer func() { _ = fdsetFile.Close() }()
	fdsetBytes, err := io.ReadAll(fdsetFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read protoset: %w", err)
	}

	fdset := &descriptorpb.FileDescriptorSet{}
	if err = proto.Unmarshal(fdsetBytes, fdset); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal protoset file %s: %w", protosetPath, err)
	}

	return c.convertToMethodInfo(fdset)
}

// Connect is a block dial to the gRPC server at the given address (host:port)
func (c *Client) Connect(addr string, params map[string]interface{}) (bool, error) {
	state := c.vu.State()
	if state == nil {
		return false, common.NewInitContextError("connecting to a gRPC server in the init context is not supported")
	}

	p, err := c.parseConnectParams(params)
	if err != nil {
		return false, fmt.Errorf("invalid grpc.connect() parameters: %w", err)
	}

	opts := xgrpc_conn.DefaultOptions(c.vu)

	var tcred credentials.TransportCredentials
	if !p.IsPlaintext {
		tlsCfg := state.TLSConfig.Clone()
		tlsCfg.NextProtos = []string{"h2"}

		// TODO(rogchap): Would be good to add support for custom RootCAs (self signed)
		tcred = credentials.NewTLS(tlsCfg)
	} else {
		tcred = insecure.NewCredentials()
	}
	opts = append(opts, grpc.WithTransportCredentials(tcred))

	if ua := state.Options.UserAgent; ua.Valid {
		opts = append(opts, grpc.WithUserAgent(ua.ValueOrZero()))
	}

	ctx, cancel := context.WithTimeout(c.vu.Context(), p.Timeout)
	defer cancel()

	if p.MaxReceiveSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(p.MaxReceiveSize))))
	}

	if p.MaxSendSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(p.MaxSendSize))))
	}

	c.addr = addr
	c.conn, err = xgrpc_conn.Dial(ctx, addr, opts...)
	if err != nil {
		return false, err
	}

	if !p.UseReflectionProtocol {
		return true, nil
	}
	fdset, err := c.conn.Reflect(ctx)
	if err != nil {
		return false, err
	}
	_, err = c.convertToMethodInfo(fdset)
	if err != nil {
		return false, fmt.Errorf("can't convert method info: %w", err)
	}

	return true, err
}

// Invoke creates and calls a unary RPC by fully qualified method name
func (c *Client) Invoke(
	method string,
	req goja.Value,
	params goja.Value,
) (*xgrpc_conn.Response, error) {
	state := c.vu.State()
	if state == nil {
		return nil, common.NewInitContextError("invoking RPC methods in the init context is not supported")
	}
	if c.conn == nil {
		return nil, errors.New("no gRPC connection, you must call connect first")
	}
	if method == "" {
		return nil, errors.New("method to invoke cannot be empty")
	}
	if method[0] != '/' {
		method = "/" + method
	}
	methodDesc := c.mds[method]
	if methodDesc == nil {
		return nil, fmt.Errorf("method %q not found in file descriptors", method)
	}

	p, err := c.parseInvokeParams(params)
	if err != nil {
		return nil, fmt.Errorf("invalid grpc.invoke() parameters: %w", err)
	}

	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	b, err := req.ToObject(c.vu.Runtime()).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("unable to serialise request object: %w", err)
	}

	md := metadata.New(nil)
	for param, strval := range p.Metadata {
		md.Append(param, strval)
	}

	ctx, cancel := context.WithTimeout(c.vu.Context(), p.Timeout)
	defer cancel()

	if state.Options.SystemTags.Has(metrics.TagURL) {
		p.TagsAndMeta.SetSystemTagOrMeta(metrics.TagURL, fmt.Sprintf("%s%s", c.addr, method))
	}
	parts := strings.Split(method[1:], "/")
	p.TagsAndMeta.SetSystemTagOrMetaIfEnabled(state.Options.SystemTags, metrics.TagService, parts[0])
	p.TagsAndMeta.SetSystemTagOrMetaIfEnabled(state.Options.SystemTags, metrics.TagMethod, parts[1])

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := p.TagsAndMeta.Tags.Get("name"); !ok {
		p.TagsAndMeta.SetSystemTagOrMetaIfEnabled(state.Options.SystemTags, metrics.TagName, method)
	}

	reqmsg := xgrpc_conn.Request{
		MethodDescriptor: methodDesc,
		Message:          b,
		TagsAndMeta:      &p.TagsAndMeta,
	}

	r, e := c.conn.Invoke(ctx, method, md, reqmsg)

	return r, e

}

// Close will close the client gRPC connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}

// MethodInfo holds information on any parsed method descriptors that can be used by the goja VM
type MethodInfo struct {
	Package         string
	Service         string
	FullMethod      string
	grpc.MethodInfo `json:"-" js:"-"`
}

func (c *Client) convertToMethodInfo(fdset *descriptorpb.FileDescriptorSet) ([]MethodInfo, error) {
	files, err := protodesc.NewFiles(fdset)
	if err != nil {
		return nil, err
	}
	var rtn []MethodInfo
	if c.mds == nil {
		// This allows us to call load() multiple times, without overwriting the
		// previously loaded definitions.
		c.mds = make(map[string]protoreflect.MethodDescriptor)
	}
	appendMethodInfo := func(
		fd protoreflect.FileDescriptor,
		sd protoreflect.ServiceDescriptor,
		md protoreflect.MethodDescriptor,
	) {
		name := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())
		c.mds[name] = md
		rtn = append(rtn, MethodInfo{
			MethodInfo: grpc.MethodInfo{
				Name:           string(md.Name()),
				IsClientStream: md.IsStreamingClient(),
				IsServerStream: md.IsStreamingServer(),
			},
			Package:    string(fd.Package()),
			Service:    string(sd.Name()),
			FullMethod: name,
		})
	}
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		sds := fd.Services()
		for i := 0; i < sds.Len(); i++ {
			sd := sds.Get(i)
			mds := sd.Methods()
			for j := 0; j < mds.Len(); j++ {
				md := mds.Get(j)
				appendMethodInfo(fd, sd, md)
			}
		}
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			_, errFind := protoregistry.GlobalTypes.FindMessageByName(message.FullName())
			if errors.Is(errFind, protoregistry.NotFound) {
				err = protoregistry.GlobalTypes.RegisterMessage(dynamicpb.NewMessageType(message))
				if err != nil {
					return false
				}
			}
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

type invokeParams struct {
	Metadata    map[string]string
	TagsAndMeta metrics.TagsAndMeta
	Timeout     time.Duration
}

func (c *Client) parseInvokeParams(paramsVal goja.Value) (*invokeParams, error) {
	result := &invokeParams{
		Timeout:     1 * time.Minute,
		TagsAndMeta: c.vu.State().Tags.GetCurrentValues(),
	}
	if paramsVal == nil || goja.IsUndefined(paramsVal) || goja.IsNull(paramsVal) {
		return result, nil
	}
	rt := c.vu.Runtime()
	params := paramsVal.ToObject(rt)
	for _, k := range params.Keys() {
		switch k {
		case "headers":
			c.vu.State().Logger.Warn("The headers property is deprecated, replace it with the metadata property, please.")
			fallthrough
		case "metadata":
			result.Metadata = make(map[string]string)
			v := params.Get(k).Export()
			rawHeaders, ok := v.(map[string]interface{})
			if !ok {
				return result, errors.New("metadata must be an object with key-value pairs")
			}
			for hk, kv := range rawHeaders {
				// TODO(rogchap): Should we manage a string slice?
				strval, ok := kv.(string)
				if !ok {
					return result, fmt.Errorf("metadata %q value must be a string", hk)
				}
				result.Metadata[hk] = strval
			}
		case "tags":
			if err := common.ApplyCustomUserTags(rt, &result.TagsAndMeta, params.Get(k)); err != nil {
				return result, fmt.Errorf("metric tags: %w", err)
			}
		case "timeout":
			var err error
			v := params.Get(k).Export()
			result.Timeout, err = types.GetDurationValue(v)
			if err != nil {
				return result, fmt.Errorf("invalid timeout value: %w", err)
			}
		default:
			return result, fmt.Errorf("unknown param: %q", k)
		}
	}
	return result, nil
}

type connectParams struct {
	IsPlaintext           bool
	UseReflectionProtocol bool
	Timeout               time.Duration
	MaxReceiveSize        int64
	MaxSendSize           int64
}

func (c *Client) parseConnectParams(raw map[string]interface{}) (connectParams, error) {
	params := connectParams{
		IsPlaintext:           false,
		UseReflectionProtocol: false,
		Timeout:               time.Minute,
		MaxReceiveSize:        0,
		MaxSendSize:           0,
	}
	for k, v := range raw {
		switch k {
		case "plaintext":
			var ok bool
			params.IsPlaintext, ok = v.(bool)
			if !ok {
				return params, fmt.Errorf("invalid plaintext value: '%#v', it needs to be boolean", v)
			}
		case "timeout":
			var err error
			params.Timeout, err = types.GetDurationValue(v)
			if err != nil {
				return params, fmt.Errorf("invalid timeout value: %w", err)
			}
		case "reflect":
			var ok bool
			params.UseReflectionProtocol, ok = v.(bool)
			if !ok {
				return params, fmt.Errorf("invalid reflect value: '%#v', it needs to be boolean", v)
			}
		case "maxReceiveSize":
			var ok bool
			params.MaxReceiveSize, ok = v.(int64)
			if !ok {
				return params, fmt.Errorf("invalid maxReceiveSize value: '%#v', it needs to be an integer", v)
			}
			if params.MaxReceiveSize < 0 {
				return params, fmt.Errorf("invalid maxReceiveSize value: '%#v, it needs to be a positive integer", v)
			}
		case "maxSendSize":
			var ok bool
			params.MaxSendSize, ok = v.(int64)
			if !ok {
				return params, fmt.Errorf("invalid maxSendSize value: '%#v', it needs to be an integer", v)
			}
			if params.MaxSendSize < 0 {
				return params, fmt.Errorf("invalid maxSendSize value: '%#v, it needs to be a positive integer", v)
			}

		default:
			return params, fmt.Errorf("unknown connect param: %q", k)
		}
	}
	return params, nil
}

func walkFileDescriptors(seen map[string]struct{}, fd *desc.FileDescriptor) []*descriptorpb.FileDescriptorProto {
	fds := []*descriptorpb.FileDescriptorProto{}

	if _, ok := seen[fd.GetName()]; ok {
		return fds
	}
	seen[fd.GetName()] = struct{}{}
	fds = append(fds, fd.AsFileDescriptorProto())

	for _, dep := range fd.GetDependencies() {
		deps := walkFileDescriptors(seen, dep)
		fds = append(fds, deps...)
	}

	return fds
}
