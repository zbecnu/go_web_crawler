package analyzer

import (
	"errors"
	"fmt"
	"zhangbing/webcrawler/module"
	"zhangbing/webcrawler/module/stub"
	"zhangbing/webcrawler/reader"
)

type myAnalyzer struct {
	stub.ModuleInternal
	respParsers []module.ParseResponse
}

func (analyzer *myAnalyzer) RespParses() []module.ParseResponse {
	panic("implement me")
}

func (analyzer *myAnalyzer) Analyze(resp *module.Response) (dataList []module.Data, errorList []error) {
	panic("implement me")
	analyzer.ModuleInternal.IncrHandlingNumber()
	defer analyzer.ModuleInternal.DecrHandlingNumber()

	analyzer.ModuleInternal.IncrCalledCount()

	if resp ==nil{
		errorList =append(errorList,errors.New(""))
	}

	httpResp := resp.HTTPResp()

	if httpResp ==nil{
		errorList =append(errorList,errors.New(""))
		return
	}

	httpReq := httpResp.Request
	if httpReq ==nil{
		errorList =append(errorList,errors.New(""))
		return
	}

	var  reqURL = httpReq.URL

	if reqURL == nil{
		errorList =append(errorList,errors.New(""))
		return
	}

	analyzer.ModuleInternal.IncrAcceptedCount()
	respDepth := resp.Depth()
	fmt.Println(respDepth)
	if httpResp.Body != nil{
		defer  httpResp.Body.Close()
	}

	multipleReader,err := reader.NewMutipleReader(httpResp.Body)
	if err !=nil{
		errorList =append(errorList,errors.New(""))
		return
	}
	//初始化切片
	dataList = []module.Data{}
	//dataList = make([]module.Data,10)
	for _,respParser := range analyzer.respParsers{
		httpResp.Body = multipleReader.Reader()
		pDataList,pErrorList := respParser(httpResp,respDepth)
		if pDataList !=nil{
			for _,PData :=range pDataList{
				if PData==nil{
					continue
				}
				dataList=appendDataList(dataList,PData,respDepth)

			}
		}

		if pErrorList !=nil{
			for _,pError :=range pErrorList{
				if pError ==nil{
					continue
				}
				errorList = append(errorList,pError)

			}
		}
	}
	if len(errorList) ==0{
		analyzer.ModuleInternal.IncrcompetedCount()

	}
	return dataList,errorList

}

func appendDataList(dataList []module.Data,PData module.Data,depth uint32 ) []module.Data{
	// 怎么实现
	return dataList
}

func new(mid module.MID,respParsers []module.ParseResponse,scoreCalculator module.CaculatorScore)(module.Analyzer,error) {
	moduleBase,err:=stub.NewModuleInternal(mid,scoreCalculator)

	if err != nil{
		return nil,err
	}
	if respParsers == nil{
		return nil,errors.New("")
	}
	if len(respParsers) == 0 {
		return nil,errors.New("")

	}
	var innerParsers []module.ParseResponse
	for _,parser :=range respParsers{
		if parser ==nil{
			return nil,errors.New("")
		}
		innerParsers =append(innerParsers,parser)
	}
	return &myAnalyzer{ModuleInternal:moduleBase,respParsers:innerParsers},nil

}