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
	Objs     []*PoolObj
	IdMin    int
	IdMax    int
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
	item.Init()
	p.Items[addr] = item
	return item
}

// 连接池对象初始化
func (item *PoolItem) Init() {
	item.Lock()
	defer item.Unlock()
	for i := len(item.Objs); i < item.IdMin; i++ {
		obj, err := item.NewObj()
		if err != nil {
			break
		}
		item.Objs = append(item.Objs, obj)
	}
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

// 获取连接
func (item *PoolItem) Get() (*PoolObj, error) {
	item.Lock()
	defer item.Unlock()

	// 不足max时，
	if len(item.Objs) < item.IdMax {
		return item.NewObj()
	}

	// 找到连接
	index := -1
	timeNow := time.Now().Unix()
	for i := 0; i < len(item.Objs); i++ {
		if item.Objs[i].EndTime > timeNow && item.Objs[i].Client.Status() {
			index = i
			break
		}
	}
	// 没有找到 执行初始化, 返回一个新创建的
	if index == -1 || len(item.Objs)-index < item.IdMin {
		item.Init()
		return item.NewObj()
	}

	obj := item.Objs[index]

	// 清理一次
	item.Objs = item.Objs[index+1:]

	return obj, nil
}

// 回收连接
func (item *PoolItem) Recycle(obj *PoolObj) {
	item.Lock()
	defer item.Unlock()

	item.Objs = append(item.Objs, obj)
}
