package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game_zinx/core"
	"mmo_game_zinx/pb/pb"
	"zinx/ziface"
	"zinx/znet"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	//解析客户端传入 proto协议数据
	protoMsg := &pb.Position{}
	err := proto.Unmarshal(request.GetMsgData(), protoMsg)
	if err != nil {
		fmt.Println("Unmarshal message err", err)
		return
	}
	//获得当前聊天的发起者
	pid, _ := request.GetConnection().GetProperty("pid")

	//根据 pid得到 Player
	player := core.WorldMgrObj.GetPlayerByPid(pid.(uint32))

	//更新当前玩家位置信息,并同步通知周边玩家
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
