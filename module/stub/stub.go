package stub

import (
	"errors"
	"sync/atomic"
	"zhangbing/webcrawler/module"
)

type ModuleInternal interface {
	module.Module
//	调用计数器增1
	IncrCalledCount()
//	接受计数增1
	IncrAcceptedCount()
//	 成功完成计数增1
	IncrcompetedCount()
//	把实时处理计数增1
	IncrHandlingNumber()
//实时处理减1
	DecrHandlingNumber()
//	用于清空计数
	Clear()
}

type myModule struct {
	mid module.MID
	addr string
	score uint64
	scoreCalculator module.CaculatorScore
	calledCount uint64
	acceptedCount uint64
	compeletedCount uint64
	handlingNumber uint64
}

func (mm *myModule) ID() module.MID {
	return mm.mid
	//panic("implement me")
}

func (mm *myModule) Addr() string {
	return mm.addr
	//panic("implement me")
}

func (mm *myModule) Score() uint64 {
	return atomic.LoadUint64(&mm.score)
}

func (mm *myModule) SetScore(score uint64) {
	old :=mm.Score()
	atomic.CompareAndSwapUint64(&mm.score,old,score)
}

func (mm *myModule) ScoreCaculator() module.CaculatorScore {
	//panic("implement me")
	return mm.scoreCalculator
}


func (mm *myModule) CalledCount() uint64 {
	//panic("implement me")
	return atomic.LoadUint64(&mm.calledCount)
}


func (mm *myModule) AcceptedCount() uint64 {
	//panic("implement me")
	return atomic.LoadUint64(&mm.acceptedCount)
}

func (mm *myModule) CompeletedCount() uint64 {
	//panic("implement me")
	return  atomic.LoadUint64(&mm.compeletedCount)
}

func (mm *myModule) HandlingNumber() uint64 {
	//panic("implement me")
	return atomic.LoadUint64(&mm.handlingNumber)
}

func (mm *myModule) Counts() module.Counts {
	return module.Counts{CalledCount:mm.CalledCount(),AcceptedCount:mm.AcceptedCount(),CompeletedCount:mm.CompeletedCount(),HandlingCount:mm.HandlingNumber()}
}

func (mm *myModule) Summary() module.SummaryStruct {
	panic("implement me")
}

func (mm *myModule) IncrCalledCount() {
	panic("implement me")
}

func (mm *myModule) IncrAcceptedCount() {
	panic("implement me")
}

func (mm *myModule) IncrcompetedCount() {
	panic("implement me")
}

func (mm *myModule) IncrHandlingNumber() {
	panic("implement me")
}

func (mm *myModule) DecrHandlingNumber() {
	panic("implement me")
}

func (mm *myModule) Clear() {
	panic("implement me")
}

func NewModuleInternal(mid module.MID,scoreCalculator module.CaculatorScore)(ModuleInternal,error){
	parts,err :=module.SplidMID(mid)
	if err !=nil{
		return nil,errors.New("")
	}

	return &myModule{mid:mid,addr:parts[2],scoreCalculator:scoreCalculator},nil
}