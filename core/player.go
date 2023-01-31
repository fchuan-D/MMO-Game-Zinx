package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"mmo_game_zinx/pb/pb"
	"sync"
	"zinx/ziface"
)

// Player 玩家管理模块
type Player struct {
	Pid  uint32             //玩家ID
	Conn ziface.IConnection //当前玩家的连接(与客户端的连接)
	X    float32            //平面 X坐标
	Y    float32            //高度
	Z    float32            //平面 Y坐标
	V    float32            //旋转角度 (0-360)
}

/*
Player ID生成器
*/
var PidGen uint32 = 1 //玩家ID计数器
var IdLock sync.Mutex //计数器加锁

// NewPlayer 创建玩家
func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	//生产一个玩家ID
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建玩家对象
	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(150 + rand.Intn(12)),
		Y:    0,
		Z:    float32(150 + rand.Intn(12)),
		V:    0,
	}
}

/*
服务器向发送客户端消息
将 pb的 protobuf数据序列化后再调用 Zinx的 SendMsg()
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将 proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("Marshal data err:", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("Connection in player is nil")
		return
	}

	//将二进制文件 通过Zinx框架SendMsg()发送给客户端
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg err")
		return
	}
}

// SyncPid 告知客户端玩家 Pid，同步已生成的 Pid给客户端
func (p *Player) SyncPid() {
	//创建 msgId = 0 的 proto数据
	protoMsg := &pb.SyncPid{
		Pid: int32(p.Pid),
	}

	p.SendMsg(1, protoMsg)
}

// BroadCastStartPos 告知客户端玩家坐标位置
func (p *Player) BroadCastStartPos() {
	//创建 msgId = 200 的 proto数据
	protoMsg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  2, //广播坐标位置
		Data: &pb.BroadCast_P{
			//初始位置
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, protoMsg)
}

// Talk 广播聊天数据
func (p *Player) Talk(content string) {
	//创建 msgId = 200 的 proto数据
	protoMsg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  1, //广播世界聊天
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//获取当前所有在线玩家
	players := WorldMgrObj.GetAllPlayer()

	//遍历在线玩家并广播
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// SyncSurround 同步周边玩家,广播当前玩家位置信息
func (p *Player) SyncSurround(tp int32) {
	//根据当前玩家位置获取周围玩家ID
	pids := WorldMgrObj.AoiMgr.GetSurroundPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		//得到周围玩家集合
		players = append(players, WorldMgrObj.GetPlayerByPid(uint32(pid)))
	}

	//将当前玩家位置信息 通过 msgId = 200 发给周边玩家(让其他玩家看到自己)
	//创建 MsgId = 200 的 proto数据
	protoMsg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  tp, //广播出生位置信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}

	//将周边玩家位置信息 通过 msgID = 202 发给当前玩家客户端(让自己看到其他玩家)
	//创建 pb.Player切片
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	//遍历周边玩家
	for _, player := range players {
		//创建一个 pb.player
		p := &pb.Player{
			Pid: int32(player.Pid),
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, p)
	}

	syncPlayersProtoMsg := &pb.SyncPlayers{
		Ps: playersProtoMsg[:],
	}

	p.SendMsg(202, syncPlayersProtoMsg)
}

// UpdatePos 更新当前玩家位置信息,并同步通知周边玩家
func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	//同步通知周边玩家
	p.SyncSurround(4)
}

// OffLine 玩家下线业务
func (p *Player) OffLine() {
	//得到当前所有玩家
	players := WorldMgrObj.GetAllPlayer()

	//创建 proto数据,当前玩家下线
	protoMsg := &pb.SyncPid{
		Pid: int32(p.Pid),
	}

	//广播给所有玩家
	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}

	//将下线玩家移出 AoiMgr
	WorldMgrObj.AoiMgr.RemovePidByPos(int(p.Pid), p.X, p.Z)
	//将下线玩家移出 WorldMgr
	WorldMgrObj.RemovePlayer(p)
}
