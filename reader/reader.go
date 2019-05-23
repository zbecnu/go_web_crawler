package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"zhangbing/webcrawler/toolkit"
)
//多重读取器实现
type myMutipleReader struct {
	data []byte
}

func (rr *myMutipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}

func NewMutipleReader(reader io.Reader)(toolkit.MutipleReader,error){
	var data []byte
	var err error
	if reader !=nil{
		data,err = ioutil.ReadAll(reader)
		if err !=nil{
			return nil,fmt.Errorf("%s",err)
		}
	}else{
		data =[]byte{}
	}
	//myMutipleReader
	//return &myM,nil
	return &myMutipleReader{data:data},nil
}
