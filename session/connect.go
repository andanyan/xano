package session

import (
	"fmt"
	"sync"
)

// 连接数管理
type Connect struct {
	sync.Mutex
	// 当前连接总数
	count uint64
	// 最后sid
	sid uint64
	// sid可取值范围
	scope []uint64
}

var connect *Connect

// 获取connect
func GetConnect() *Connect {
	if connect == nil {
		connect = new(Connect)
	}
	return connect
}

// 判断sid是否充足
func (c *Connect) IsEnough() bool {
	c.Lock()
	defer c.Unlock()

	slen := len(c.scope)
	if slen == 0 {
		return false
	}

	// 找到有效值的范围索引
	var index int
	for i := 0; i < slen; i += 2 {
		if c.scope[i] > c.sid || c.scope[i+1] > c.sid {
			index = i
			break
		}
	}
	if index > 0 {
		c.scope = c.scope[index:]
	}

	// 完全不足
	if index >= slen {
		return false
	}

	// 最后的取值范围
	if index+2 >= slen {
		total := c.scope[index+1] - c.scope[index] + 1
		sum := c.scope[index+1] - c.sid
		if sum*10 < total {
			return false
		}
	}

	// 剩余的比较充足
	return true
}

// 获取一个sid
func (c *Connect) GetSid() (uint64, error) {
	c.Lock()
	defer c.Unlock()

	slen := len(c.scope)
	if slen == 0 {
		return 0, fmt.Errorf("not enough sid")
	}

	sid := c.sid + 1
	if sid < c.scope[0] || sid > c.scope[1] {
		return 0, fmt.Errorf("not enough sid")
	}

	c.count++
	c.sid = sid

	return sid, nil
}

// 消除一个sid
func (c *Connect) DelSid(sid uint64) {
	c.Lock()
	defer c.Unlock()
	c.count--
}

// 增加新的范围
func (c *Connect) AddScope(min, max uint64) {
	c.Lock()
	defer c.Unlock()
	c.scope = append(c.scope, min, max)
}

// 获取当前连接数
func (c *Connect) GetCount() uint64 {
	return c.count
}
