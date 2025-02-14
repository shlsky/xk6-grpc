package grpc

import (
	"github.com/grafana/sobek"
	"strconv"
	"time"
)

type Util struct {
}

// NewUtil is the JS constructor for the grpc Util.
func (mi *ModuleInstance) NewUtil(_ sobek.ConstructorCall) *sobek.Object {
	rt := mi.vu.Runtime()
	return rt.ToValue(&Util{}).ToObject(rt)
}

// GetNano 获取纳秒
func (c *Util) GetNano() int64 {
	return time.Now().UnixNano()
}

// GetMicro 获取微妙
func (c *Util) GetMicro() int64 {
	return time.Now().UnixMicro()
}

// GetNanoStr 获取纳秒字符串
func (c *Util) GetNanoStr() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// GetMicroStr 获取微妙字符串
func (c *Util) GetMicroStr() string {
	return strconv.FormatInt(time.Now().UnixMicro(), 10)
}
