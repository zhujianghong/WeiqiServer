// RoomManager
package main

import (
	"sync"
	"sync/atomic"
)

type RoomManager struct {
	autoRoomId uint64
	RoomMap    map[uint64]*Room
	rwLock     sync.RWMutex
}

func newRoomManager() *RoomManager {
	m := &RoomManager{
		autoRoomId: 0,
		RoomMap:    make(map[uint64]*Room),
	}

	return m
}

func (m *RoomManager) createRoom() *Room {
	roomId := atomic.AddUint64(&m.autoRoomId, 1)
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	room := newRoom(roomId)
	m.RoomMap[roomId] = room
	return room
}

func (m *RoomManager) delRoom(id uint64) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	if _, ok := m.RoomMap[id]; ok {
		delete(m.RoomMap, id)
	}
}

func (m *RoomManager) findRoom(id uint64) *Room {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	if v, ok := m.RoomMap[id]; ok {
		return v
	}
	return nil
}
