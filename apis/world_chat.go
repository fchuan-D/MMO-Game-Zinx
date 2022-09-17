package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"lib/mmo_game_zinx/core"
	"lib/mmo_game_zinx/pb/pb"
	"lib/zinx/ziface"
	"lib/zinx/znet"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//解析客户端传入 proto协议数据
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetMsgData(), protoMsg)
	if err != nil {
		fmt.Println("Unmarshal message err", err)
	}

	//获得当前聊天的发起者
	pid, _ := request.GetConnection().GetProperty("pid")

	//根据 pid得到 Player
	player := core.WorldMgrObj.GetPlayerByPid(pid.(uint32))

	//广播聊天数据
	player.Talk(protoMsg.Content)
}
