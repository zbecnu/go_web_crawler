package pipeline

import (
	"zhangbing/webcrawler/module"
	"zhangbing/webcrawler/module/stub"
)

type myPipeline struct {
	stub.ModuleInternal
	itemProcessors []module.ProcessItem
	failFast bool

}

func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	panic("implement me")
	return pipeline.itemProcessors
}

func (pipeline *myPipeline) Send(item module.Item) []error {
	panic("implement me")
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
	//panic("implement me")
}

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	//panic("implement me")
	pipeline.failFast = failFast
}

func New(mid module.MID,scoreCalculator module.CaculatorScore,itemProcessors []module.ProcessItem,failFast bool)(module.Pipeline,error){
	moduleBase,err:=stub.NewModuleInternal(mid,scoreCalculator)

	if err != nil{
		return nil,err
	}

	return &myPipeline{ModuleInternal:moduleBase,itemProcessors:itemProcessors,failFast:failFast},nil
}