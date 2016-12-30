// PlayerManager
package main

import (
	"sync"

	"github.com/gansidui/gotcp"
)

type PlayerManager struct {
	playerIdMap   map[int32]*Player
	PlayerNameMap map[string]*Player
	rwLock        sync.RWMutex
}

func newPlayerManager() *PlayerManager {
	return &PlayerManager{
		playerIdMap:   make(map[int32]*Player),
		PlayerNameMap: make(map[string]*Player),
	}
}

func (m *PlayerManager) login(p *Player) bool {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	if v, ok := m.playerIdMap[p.userid]; ok {
		if v.conn != p.conn {
			delete(m.playerIdMap, p.userid)
			delete(m.PlayerNameMap, p.account)
		} else {
			return false
		}
		//TODO 通知对方账号在别的地方登陆
		p.conn.Close()
	}
	m.playerIdMap[p.userid] = p
	m.PlayerNameMap[p.account] = p
	return true
}

func (m *PlayerManager) logout(userid int32, conn *gotcp.Conn) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	if v, ok := m.playerIdMap[userid]; ok {
		if v.conn == conn {
			account := v.account
			delete(m.playerIdMap, userid)
			delete(m.PlayerNameMap, account)
		}
	}
}

func (m *PlayerManager) findPlayerById(userid int32) *Player {
	m.rwLock.RLocker()
	defer m.rwLock.RUnlock()
	if v, ok := m.playerIdMap[userid]; ok {
		return v
	}
	return nil
}

func (m *PlayerManager) findPlayerByAccount(account string) *Player {
	m.rwLock.RLocker()
	defer m.rwLock.RUnlock()
	if v, ok := m.PlayerNameMap[account]; ok {
		return v
	}
	return nil
}
