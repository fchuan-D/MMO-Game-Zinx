package main

import (
	"fmt"
	"mmo_game_zinx/apis"
	"mmo_game_zinx/core"
	"zinx/ziface"
	"zinx/znet"
)

// OnConnectionAdd 客户端建立连接之后的 Hook
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个 Player对象
	player := core.NewPlayer(conn)

	//给客户端发送 msgId=1 的消息,同步玩家ID
	player.SyncPid()

	//给客户端发送 msgId=200 的消息,同步玩家位置
	player.BroadCastStartPos()

	//将新上线玩家加入 WorldMgr
	core.WorldMgrObj.AddPlayer(player)

	//将连接与当前玩家ID绑定
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家,广播当前玩家位置信息
	player.SyncSurround(2)

	fmt.Println("playerID:", player.Pid, "Login...")
}

// OnConnectionStop 客户端断开连接之前的 Hook
func OnConnectionStop(conn ziface.IConnection) {
	//获取下线玩家ID
	pid, _ := conn.GetProperty("pid")

	//获得下线玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(uint32))

	//玩家下线业务
	fmt.Println("playerID:", pid, "Logout...")
	player.OffLine()
}

func main() {
	//创建zinx server句柄
	s := znet.NewServer()

	//连接创建和销毁的Hook函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionStop)

	/*
		注册业务路由函数
	*/
	//世界聊天业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Serve()
}
