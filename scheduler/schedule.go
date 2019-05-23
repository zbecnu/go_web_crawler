package scheduler

import (
	"context"
	"net/http"
	"sync"
	"zhangbing/webcrawler/buffer"
	"zhangbing/webcrawler/module"
)

type Status uint8

const   (
	SCHED_STATUS_UNINITIALIZED Status = iota
	SCHED_STATUS_INITALIZING
	SCHED_STATUS_INITALIZED
	SCHED_STATUS_STARTING
	SCHED_STATUS_STARTED
	SCHED_STATUS_STOPPING
	SCHED_STATUS_STOPPED
)

//调度器
type Scheduler interface {
	Init(requestArgs RequestArgs,dataArgs DataArgs,moduleArgs ModuleArgs)(err error)
	Start(firstHTTPReq *http.Request)(err error)
	Stop() (err error)
	Status() Status
	//获取错误
	ErrorChan() <-chan error
	//判定是否空闲
	Idle() bool
	//返回调度器内部状态摘要
	Summary() SchedSummary
}

type RequestArgs struct {
	AcceptedDomain []string `json:"accepted_primary_domains"`
	MaxDepth uint `json:"max_depth"`
}

//参数容器
type DataArgs struct {
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_cap"`
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
	
}

type ModuleArgs struct {
	Downloaders	[]module.Downloader
	Analyzers []module.Analyzer
	Pipelines []module.Pipeline
}

// 参数容器接口
type Args interface {
	Check() error
}

type SummaryStruct struct {
	RequestArgs RequestArgs `json:"request_args"`
	DataArgs DataArgs `json:"data_args"`
	ModuleArgs ModuleArgs `json:"module_args"`
	Status string `json:"status"`
	Downloaders []module.Downloader `json:"downloaders"`
	Analyzer  []module.Analyzer `json:"analyzer"`
	Pipelines []module.Pipeline `json:"pipelines"`
	ReqBufferPool BufferPoolSummaryStruct `json:"req_buffer_pool"`
	RespBufferPool BufferPoolSummaryStruct `json:"resp_buffer_pool"`
	ItemBufferPool BufferPoolSummaryStruct `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct `json:"error_buffer_pool"`
	NumURL uint64 `json:"num_url"`
}

type SchedSummary interface {
	Struct() SummaryStruct
	String() string
}


type BufferPoolSummaryStruct struct {
//	???
}

func newSchedSummary() SchedSummary{
	return nil
}

type myScheduler struct {
	maxDepth uint32
	acceptedDomainMap cmap.Concurrency
	registrar module.Registrar
	reqBufferPool buffer.Pool
	respBufferPool buffer.Pool
	itemBufferPool buffer.Pool
	errorBufferPool buffer.Pool
	urlMap cmap.ConcurrentMap
	ctx context.Context
	concanFunc context.CancelFunc
	status Status
	statusLock sync.RWMutex
	SummaryStruct SchedSummary
}

