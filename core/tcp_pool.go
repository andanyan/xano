package core

import (
	"fmt"
	"sync"
	"time"
	"xlq-server/common"
)

type Pool struct {
	sync.RWMutex
	// 连接池 addr -- pool
	Items map[string]*PoolItem
}

// 连接池对象
type PoolItem struct {
	sync.RWMutex
	Objs chan *PoolObj
	// 最小数量
	IdMin int
	// 最大数量
	IdMax int
	// 当前数量 总数
	Num      int
	LifeTime int64
	Addr     string
}

// 连接对象
type PoolObj struct {
	EndTime int64
	Client  *TcpClient
}

var pool *Pool

func GetPool(addr string) *PoolItem {
	if pool == nil {
		pool = new(Pool)
		pool.Items = make(map[string]*PoolItem)
	}

	pool.Lock()
	defer pool.Unlock()

	item, ok := pool.Items[addr]
	if ok {
		return item
	}

	item = pool.NewItem(addr)

	// 新建pool
	return item
}

// 创建连接池
func (p *Pool) NewItem(addr string) *PoolItem {
	item := new(PoolItem)
	item.IdMin = common.TcpPoolIdMin
	item.IdMax = common.TcpPoolIdMax
	item.LifeTime = common.TcpPoolLifeTime
	item.Addr = addr
	item.Objs = make(chan *PoolObj, common.TcpPoolIdMax)
	item.Init()
	p.Items[addr] = item
	return item
}

// 连接池对象初始化
func (item *PoolItem) Init() {
	for i := 0; i < item.IdMin; i++ {
		obj, err := item.NewObj()
		if err != nil {
			break
		}
		item.Objs <- obj
	}
	item.Num = item.IdMin
}

// 新建连接
func (item *PoolItem) NewObj() (*PoolObj, error) {
	cli, err := NewTcpClient(item.Addr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	timeNow := time.Now().Unix()
	obj := &PoolObj{
		EndTime: timeNow + item.LifeTime,
		Client:  cli,
	}
	return obj, nil
}

// 获取连接：优先提供空闲，无空闲并且没有达到最大连接数则创建，否则等待最长时间后再获取
func (item *PoolItem) Get() (*PoolObj, error) {
	// 无空闲、未达最大连接
	if item.Num < item.IdMax {
		obj, err := item.NewObj()
		if err != nil {
			return nil, err
		}
		item.Num++
		return obj, nil
	}
	// 有空闲立即获取
	// 无空闲、已达最大连接
	t := time.NewTimer(time.Duration(common.TcpPoolMaxWaitTime) * time.Millisecond)
	select {
	case obj := <-item.Objs:
		return obj, nil
	case <-t.C:
		return nil, fmt.Errorf("timeout")
	}
}

// 回收连接
func (item *PoolItem) Recycle(obj *PoolObj) {
	item.Lock()
	defer item.Unlock()

	item.Objs <- obj
}
