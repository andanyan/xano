package session

import (
	"sync"
	"xano/common"
)

// 用于管理member的session
type Member struct {
	sync.RWMutex
	// session集合
	Sessions map[uint64]*Session
	// 当前SID
	SID uint64
}

var member *Member

func GetMember() *Member {
	if member == nil {
		member = new(Member)
		member.Sessions = make(map[uint64]*Session)
	}
	return member
}

// session init
func (m *Member) SesssionInit(s *Session) {
	m.Lock()
	defer m.Unlock()

	var sid uint64
	for {
		if m.SID >= common.MaxSessionNum {
			m.SID = 0
		}
		m.SID++
		mchID := common.GetCache().GetUInt64(common.MchIDKey)
		sid = mchID*common.MaxSessionNum + m.SID
		// 一直寻找 直到找到合适的 一般情况下一次就可以找到
		if m.Sessions[sid] == nil {
			break
		}
	}

	s.SID = sid
	m.Sessions[sid] = s
}

// session close
func (m *Member) SessionClose(s *Session) {
	m.Lock()
	defer m.Unlock()
	sid := s.GetSid()
	delete(m.Sessions, sid)
}

// session find
func (m *Member) SessionFindByID(sid uint64) *Session {
	m.RLock()
	defer m.RUnlock()
	return m.Sessions[sid]
}

// 获取当前连接总数
func (m *Member) SessionCount() int {
	return len(m.Sessions)
}
