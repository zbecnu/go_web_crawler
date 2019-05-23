package toolkit

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

//type Pool interface {
//	BufferCap() uint32
//	MaxBufferNumber() uint32
//	BufferNumber() uint32
//	Total() uint64
//	//本方法是阻塞的
//	Put(datum interface{})error
//	Get()(datum interface{},err error)
//	Close() bool
//	//返回缓冲池的关闭状态
//	Closed() bool
//}

type Buffer interface {
	Cap() uint32
	Len() uint32
	Put(datum interface{}) (bool,error)
	Get()(interface{},error)
	Close() bool
	Closed() bool
}

//多重读取器接口
type MutipleReader interface {
	Reader() io.ReadCloser
}

//func NewMutipleReader(reader io.Reader)(MutipleReader,error){
//	var data []byte
//	var err error
//	if reader !=nil{
//		data,err = ioutil.ReadAll(reader)
//		if err !=nil{
//			return nil,fmt.Errorf("%s",err)
//		}
//	}else{
//		data =[]byte{}
//	}
//	return &myMutipleReader{data:data},nil
//}

//缓冲器实现
type myBuffer struct {
	ch chan interface{}
	closed uint32
	closingLock sync.RWMutex
}

func (buff *myBuffer)Cap() uint32{
	return uint32((cap(buff.ch)))
}

func (buff *myBuffer)Len() uint32{
	return uint32((len(buff.ch)))
}

func (buff *myBuffer)Closed() bool{
	if atomic.LoadUint32(&buff.closed) ==0 {
		return false
	}
	return true
}


func (buff *myBuffer)Put(datnum interface{}) (ok bool,err error){
	buff.closingLock.RLock()
	defer buff.closingLock.RUnlock()
	if buff.Closed() {
		return  false,errors.New("close buffer")
	}
	select {
	case buff.ch<-datnum:
		ok =true
	default:
		ok=false
	}

	return
}

func (buff *myBuffer)Get() (interface{},error){
	select {
	case datum,ok :=<-buff.ch:
		if !ok{
			return nil,errors.New("close buffer")
		}
		return datum,nil
	default:
		return nil,nil
	}
}

func (buff *myBuffer)Close() bool{
	if atomic.CompareAndSwapUint32(&buff.closed,0,1) {
		buff.closingLock.Lock()
		close(buff.ch)
		buff.closingLock.Unlock()
		return  true

	}
	return false
}

func NewBuffer(size uint32)(Buffer,error)  {
	if size ==0{
		errmsg := fmt.Sprintf("illegal size for buffer:%d",size)
		return  nil,errors.New(errmsg)
	}
	return &myBuffer{ch:make(chan interface{},size)},nil

}



type myPool struct {
//	缓冲器的统一容量
	buffCap uint32
	maxBufferNumber uint32
	//缓冲器的实际数量
	bufferNumber uint32
	total uint64
	//存放缓冲器的通道
	bufCh chan Buffer
	//关闭状态
	closed uint32
	rwlock sync.RWMutex
}

func NewPool() *myPool{
	return nil
}

func (pool *myPool)BufferCap() uint32 {
	return pool.buffCap
}

func (pool *myPool)MaxBufferNumber() uint32{
	return pool.maxBufferNumber
}

func (pool *myPool)BufferNumber() uint32{
	return atomic.LoadUint32(&pool.bufferNumber)
}

func (pool *myPool)Total()uint64{
	return atomic.LoadUint64(&pool.total)
}


func(pool *myPool) Close() bool{
	if !atomic.CompareAndSwapUint32(&pool.closed,0,1){
		return false
	}
	pool.rwlock.Lock()


	defer pool.rwlock.Unlock()
	close(pool.bufCh)
	for buf:=range pool.bufCh{
		buf.Close()
	}
	return true
}



//??/
func(pool *myPool)Closed()bool{
	return false
}

func (pool *myPool)Put(datum interface{}) (err error){
	if pool.Closed() {
		return errors.New("closed")
	}
	var count uint32
	maxCount :=pool.BufferNumber()*5
	var ok bool
	for buf := range  pool.bufCh{
		ok,err = pool.putData(buf,datum,&count,maxCount)
		if ok||err !=nil{
			break

		}
	}
	return
}

func (pool *myPool)putData(buf Buffer,datum interface{},count *uint32,maxCount uint32)(ok bool,err error){
//
	if pool.Closed(){
		return  false,errors.New("pool closed")

	}
	defer func() {
		pool.rwlock.RLock()
		if pool.Closed(){
			atomic.AddUint32(&pool.bufferNumber,^uint32(0))
			err = errors.New("pool closed")
		}else{
			pool.bufCh<-buf

		}
		pool.rwlock.RUnlock()
	}()
	ok,err =buf.Put(datum)
	if ok{
		atomic.AddUint64(&pool.total,1)
		return
	}
	if err!=nil{
		return
	}
	(*count)++
	if *count >=maxCount && pool.BufferNumber() <pool.MaxBufferNumber(){
		pool.rwlock.Lock()
		if pool.BufferNumber() <pool.MaxBufferNumber(){
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}

			newBuf,_:=NewBuffer(pool.buffCap)
			newBuf.Put(datum)
			pool.bufCh<-newBuf
			atomic.AddUint32(&pool.bufferNumber,1)
			atomic.AddUint64(&pool.total,1)
			ok = true
		}

		pool.rwlock.Unlock()
		*count = 0
	}

	return false,nil
}

func (pool *myPool) Get()(datum interface{},err error){
	if pool.Closed() {
		return nil,errors.New("")

	}
	var count uint32
	maxCount :=pool.BufferNumber()*10

	for buf :=range pool.bufCh{
		datum,err = pool.getData(buf,&count,maxCount)
		if datum!=nil||err !=nil{
			break
		}
	}
	return
}

func (pool *myPool)getData(buf Buffer,count *uint32,maxCount uint32)(datum interface{},err error){
//
	if pool.Closed(){
		return nil,errors.New("")
	}
	 defer func() {
		 if *count>=maxCount && buf.Len() ==0 &&pool.BufferNumber()>1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber,^uint32(0))
			*count=0
			 return
		 }
		 pool.rwlock.RLock()
		 if pool.Closed(){
			atomic.AddUint32(&pool.bufferNumber,^uint32(0))
			err = errors.New("")
		 }else {
			pool.bufCh<-buf
		 }
		 pool.rwlock.RUnlock()
	 }()
	datum,err = buf.Get()
	if datum !=nil{
		atomic.AddUint64(&pool.total,^uint64(0))
		return
	}
	if err!=nil{
		return
	}
	(*count)++
	return
}


//type pool interface {
//	BufferCap() uint32
//	MaxBufferNumber() uint32
//	BufferNumber() uint32
//	Total() uint64
//	//本方法是阻塞的
//	Put(datum interface{})error
//	Get()(datum interface{},err error)
//	Close() bool
//	//返回缓冲池的关闭状态
//	Closed() bool

//}