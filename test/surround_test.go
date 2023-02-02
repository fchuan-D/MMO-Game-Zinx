package test

import (
	"fmt"
	"mmo_game_zinx/core"
	"testing"
)

func TestAOIManagerGetSurround(t *testing.T) {
	aoiMgr := core.NewAOIManager(0, 200, 2, 0, 400, 4)

	for gid := range aoiMgr.Grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid:", gid, "grids len =", len(grids))
		gids := make([]int, 0, 9)
		for _, grid := range grids {
			gids = append(gids, grid.GID)
		}
		fmt.Println("surround grid GID:", gids)
	}
}
