package core

import (
	"sync"
)

// 初始常量值
const (
	AOI_MIN_X  = 0
	AOI_MAX_X  = 400
	AOI_CNTS_X = 10
	AOI_MIN_Y  = 0
	AOI_MAX_Y  = 400
	AOI_CNTS_Y = 10
)

// WorldManager 当前游戏的世界管理模块
type WorldManager struct {
	//AOIManager 当前世界地图AOI的管理模块
	AoiMgr *AOIManager
	//当前全部在线的 Players集合
	Players map[uint32]*Player
	//保护锁
	PLock sync.RWMutex
}

// WorldMgrObj WorldMgr 提供一个对外 WorldManager句柄
var WorldMgrObj *WorldManager

// 初始化
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[uint32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(p *Player) {
	//加入在线玩家集合
	wm.PLock.Lock()
	wm.Players[p.Pid] = p
	wm.PLock.Unlock()

	//player 添加进AOIManager
	wm.AoiMgr.AddPidToGridByPos(int(p.Pid), p.X, p.Z)
}

// RemovePlayer 删除一个玩家
func (wm *WorldManager) RemovePlayer(p *Player) {
	wm.PLock.Lock()
	delete(wm.Players, p.Pid)
	wm.PLock.Unlock()

	wm.AoiMgr.RemovePidByPos(int(p.Pid), p.X, p.Z)
}

// GetPlayerByPid 通过 Pid查询一个玩家
func (wm *WorldManager) GetPlayerByPid(pid uint32) *Player {
	wm.PLock.RLock()
	defer wm.PLock.RUnlock()
	return wm.Players[pid]
}

// GetAllPlayer 获取全部玩家
func (wm *WorldManager) GetAllPlayer() []*Player {
	wm.PLock.RLock()
	defer wm.PLock.RUnlock()

	players := make([]*Player, 0)

	for _, p := range wm.Players {
		players = append(players, p)
	}

	return players
}
