package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI地图的格子类型
*/
type Grid struct {
	//格子ID
	GID int

	//格子左边界
	MinX int

	//格子右边界
	MaxX int

	//格子上边界
	MinY int

	//格子下边界
	MaxY int

	//当前格子内玩家或物体ID集合
	playerIDs map[int]bool

	//当前集合的保护锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化当前格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 给当前格子增加玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true

}

// Delete 给当前格子删除玩家
func (g *Grid) Delete(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 得到当前格子所有玩家ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

//调式使用——打印出格子基本信息
func (g *Grid) String() string {
	s := fmt.Sprintf("GID:%d ,MinX:%d ,MaxX:%d ,MinY:%d ,MaxY:%d ,playerIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
	return s
}
