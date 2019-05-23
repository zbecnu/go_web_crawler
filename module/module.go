package module

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Type string
type MID string
type CaculatorScore func(counts Counts) uint64
const(
	TYPE_DOWNLOADER Type = "downloader"
	TYPE_ANALYZER Type = "analyzer"
	TYPE_PIPELINE Type ="pipeline"
)

//组件的基础接口
type Module interface {
	//组件id
	ID() MID
	//网络地址字符串
	Addr() string
	//当前组件的评分
	Score() uint64
	//设置当前组件的评分
	SetScore(score uint64)
	//获取评分计数器
	ScoreCaculator() CaculatorScore
	//当前组件呗调用的计数
	CalledCount() uint64
	//
	AcceptedCount() uint64
	//
	CompeletedCount() uint64
	//
	HandlingNumber() uint64
	//一次性获取所有的计数
	Counts() Counts
	Summary() SummaryStruct
}

//计数信息
type Counts struct {
	CalledCount uint64
	AcceptedCount uint64
	CompeletedCount uint64
	HandlingCount uint64
}

//合法的组件类型
var legalTypeLetterMap =map[Type]string{
	TYPE_DOWNLOADER:"D",
	TYPE_ANALYZER:"A",
	TYPE_PIPELINE:"p",
}

type SummaryStruct struct {
	ID MID			`json:"id"`
	called uint64	`json:"called"`
	Accepted uint64	`json:"accepted"`
	Compeleted uint64 `json:"compeleted"`
	Handing uint64	`json:"handing"`
	Extra interface{}	`json:"extra,omitempty"`
}

//序列号生成器接口
type SNGenertor interface {
	start() uint64
	Max() uint64
	Next() uint64
	CycleCount() uint64
	Get() uint64
}


//组件注册器接口
type Registrar interface {
	Register(module Module)(bool,error)
	UnRegister(mid MID)(bool ,error)
	Get(moduleType Type)(Module,error)
	GetAllByType(moduleType Type)(map[MID]Module,error)
	GetAll() map[MID]Module
	Clear()
}

//组件注册器的实现
type myRegistrar struct {
	moduleTypeMap map[Type]map[MID]Module
	rwlock sync.RWMutex
}

func (mr *myRegistrar) Register(module Module) (bool, error) {
	//panic("implement me")
	if module ==nil{
		return false,errors.New("")
	}
	mid:=module.ID()
	parts,err := SplidMID(mid)
	if err!=nil{
		return false,err
	}
	moduleType :=legalTypeLetterMap[parts[0]]
	if !CheckType(moduleType,module){
		errMsg:=fmt.Sprintf("incorrect module type:%s",moduleType)
		return false,errors.New(errMsg)
	}

	mr.moduleTypeMap[moduleType][module.ID()] = module
	return true,nil
//	///

}

func (mr *myRegistrar) UnRegister(mid MID) (bool, error) {
	//delete(mr.moduleTypeMap,mr.moduleTypeMap[])
	panic("implement me")
}

func (mr *myRegistrar) Get(moduleType Type) (Module, error) {
	//modules,err:= mr.GetAllByType(moduleType)
	//if err !=nil{
	//	return nil,err
	//}
	//
	//minScore :=uint64(0)
	//var selectedModule Module
	//for _,module :=range modules{
	//
	//}
	panic("implement me")
}

//返回所有的module
func (mr *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module, error) {

	return mr.moduleTypeMap[moduleType],nil
	panic("implement me")
}

func (mr *myRegistrar) GetAll() map[MID]Module {
	panic("implement me")
}

func (mr *myRegistrar) Clear() {
	panic("implement me")
}

func CheckType(p Type,module Module) bool{

	return false
}

//创建一个组建注册器
func NewRegistrar() Registrar  {

	return &myRegistrar{moduleTypeMap: map[Type]map[MID]Module{},}
}

//下载器接口
type Downloader interface {
	Module
	Download(req *Request)(*Response,error)
}

/*分析器
*/
type ParseResponse func(httpResp *http.Response,respDepth uint32)([]Data,[]error)
type Analyzer interface {
	Module
	RespParses() []ParseResponse
	Analyze(resp *Response)([]Data,[]error)
}

type ProcessItem func(item Item)(result Item,err error)
type Pipeline interface {
	Module
	ItemProcessors() []ProcessItem
	Send(item Item)[]error
	FailFast() bool
	SetFailFast(failFast bool)
}

func SplidMID(mid MID)([]string,error)  {
	return nil,nil
}