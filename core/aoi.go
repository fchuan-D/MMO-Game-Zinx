package core

import "fmt"

// AOIManager AOI区域管理模块
type AOIManager struct {
	//区域左边界坐标
	MinX int
	//区域右边界坐标
	MaxX int
	//X方向格子数量
	CntsX int
	//区域上边界坐标
	MinY int
	//区域下边界坐标
	MaxY int
	//Y方向格子数量
	CntsY int
	//当前区域中格子map.key = 格子ID
	Grids map[int]*Grid
}

// NewAOIManager 初始化AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}

	/*
		初始化区域中所有的格子
	*/
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//根据x,y 计算格子ID
			//格子ID = idy * cntX + idy
			gid := y*cntsX + x

			//初始化gid格子
			aoiMgr.Grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}
	return aoiMgr
}

//得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//打印格子信息
func (m *AOIManager) String() string {
	//打印AOI
	s := fmt.Sprintf("AOIManager:\nMinX:%d ,MaxX:%d ,MinY:%d ,MaxY:%d\nGrids in AOI:\n", m.MinX, m.MaxX, m.MinY, m.MaxY)

	//打印所有格子
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid.String())
	}
	return s
}

// GetSurroundGridsByGid 根据格子GID得到周围九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	//判断gid是否在AOIManager中
	if _, ok := m.Grids[gid]; !ok {
		return
	}

	//将当前gid本身加入九宫格切片中
	grids = append(grids, m.Grids[gid])

	//判断gid的左边是否有格子?右边是否有格子
	//需要通过gid得到当前格子x轴的编号 (idx = gid % CntsX)
	idx := gid % m.CntsX

	//判断idx编号左边是否有集合
	if idx > 0 {
		grids = append(grids, m.Grids[gid-1])
	}
	//判断idx编号右边是否有集合
	if idx < m.CntsX-1 {
		grids = append(grids, m.Grids[gid+1])
	}

	//取出当前x轴格子,进行遍历,再得到每个格子上下是否还有格子
	//得到当前x轴格子的ID集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	//遍历gidsX集合中每个格子的gid
	for _, v := range gidsX {
		//需要通过gid得到当前格子y轴的编号 (idY = gid % CntsX)
		idy := v / m.CntsX
		//gid上边有无格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CntsX])
		}
		//gid下边有无格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.Grids[v+m.CntsX])
		}
	}
	return grids
}

// GetGidByPos 通过x,y横纵轴坐标得到当前格子的 GID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := int(x) - m.MinX/m.gridWidth()
	idy := int(y) - m.MinY/m.gridLength()

	return idy*m.CntsX + idx
}

// GetPidByPos 通过横纵轴坐标得到周边九宫格内全部 playerIDs
func (m *AOIManager) GetPidByPos(x, y float32) (playerIDs []int) {
	//获取当前坐标对应格子ID
	gid := m.GetGidByPos(x, y)

	//获得周围九宫格格子集合
	grids := m.GetSurroundGridsByGid(gid)

	//遍历九宫格,得到所有 playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}

	return playerIDs
}

// AddToGridByPid 添加一个 playerID到一个格子中
func (m *AOIManager) AddToGridByPid(pid, gid int) {
	m.Grids[gid].Add(pid)
}

// RemovePidByGid 移除格子中一个 playerID
func (m *AOIManager) RemovePidByGid(pid, gid int) {
	m.Grids[gid].Delete(pid)
}

// GetPidsByGid 通过 GID获取全部 playerID
func (m *AOIManager) GetPidsByGid(gid int) (playerIDs []int) {
	playerIDs = m.Grids[gid].GetPlayerIDs()
	return playerIDs
}

// AddToGridByPos 通过坐标将 playerID添加到对应格子中
func (m *AOIManager) AddToGridByPos(pid int, x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.Grids[gid].Add(pid)
}

// RemovePidByPos 通过坐标将 playerID从一个格子中删除
func (m *AOIManager) RemovePidByPos(pid int, x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.Grids[gid].Delete(pid)
}
