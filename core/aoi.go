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
	grids map[int]*Grid
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
		grids: make(map[int]*Grid),
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
			aoiMgr.grids[gid] = NewGrid(
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
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid.String())
	}

	return s
}
