package buffer

type Pool interface {
	BufferCap() uint32
	MaxBufferNumber() uint32
	BufferNumber() uint32
	Total() uint64
	//本方法是阻塞的
	Put(datum interface{})error
	Get()(datum interface{},err error)
	Close() bool
	//返回缓冲池的关闭状态
	Closed() bool
}