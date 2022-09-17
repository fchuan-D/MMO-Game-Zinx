package test

import (
	"fmt"
	"lib/mmo_game_zinx/core"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := core.NewAOIManager(100, 300, 4, 200, 450, 5)

	fmt.Println(aoiMgr.String())
}
