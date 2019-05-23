package downloader

import (
	"errors"
	"net/http"
	"zhangbing/webcrawler/module"
	"zhangbing/webcrawler/module/stub"
)

type myDownloader struct {
	stub.ModuleInternal
	httpClient http.Client
}

func(downloader *myDownloader)Download(req *module.Request)(*module.Response,error){
	downloader.ModuleInternal.IncrHandlingNumber()
	defer downloader.ModuleInternal.DecrHandlingNumber()

	downloader.ModuleInternal.IncrCalledCount()

	if req ==nil{
		return nil,errors.New("")
	}


	httpReq :=req.HTTPReq()
	if httpReq ==nil{
		return nil,errors.New("")
	}

	downloader.ModuleInternal.IncrAcceptedCount()

	httpResp,err:= downloader.httpClient.Do(httpReq)
	if err !=nil{
		return nil,err
	}
	downloader.ModuleInternal.IncrcompetedCount()
	return module.NewResponse(httpResp,req.Depth()),nil

}

func New(mid module.MID,client *http.Client,scoreCalculator module.CaculatorScore)(module.Downloader,error){
	moduleBase,err:=stub.NewModuleInternal(mid,scoreCalculator)
	if err !=nil{
		return nil,err
	}
	if client ==nil{
		return nil,errors.New("")
	}

	return &myDownloader{ModuleInternal:moduleBase,httpClient:*client},nil
}