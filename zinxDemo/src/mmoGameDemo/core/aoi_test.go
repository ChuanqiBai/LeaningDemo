package core_test

import (
	"fmt"
	"testing"

	. "github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/core"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManagerSurroundGridsByGid(t *testing.T) {

	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	for gid, _ := range aoiMgr.Grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		t.Log("gid: ", gid, " grid len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		t.Log("surrounding grid IDs are ", gIDs)
	}
}
