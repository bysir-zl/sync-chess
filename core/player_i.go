package core

import (
	"context"
	"encoding/json"
)

type Player interface {
	// 设置手牌
	SetCards(cards Cards)

	// 用户的唯一Id
	GetId() string

	// 获能进行的动作,应该根据手上的牌判断返回
	// isRounder是否该你出牌
	// card 别人打的牌
	CanActions(isRounder bool, card Card) ActionTypes

	// 通知玩家需要进行操作
	// actions 为需要玩家的动作
	NotifyNeedAction(types ActionTypes)

	// 阻塞获取玩家操作
	WaitAction(ctx context.Context) (playerAction *PlayerActionRequest, err error)

	// 当玩家超时, 或者亮倒时, 执行自动打牌
	// 应当从actions中选取一个动作
	// card为要动作的牌, 部分情况为空
	RequestActionAuto(actions ActionTypes, card Card) (playerAction *PlayerActionRequest)

	// 响应玩家操作(通常在这里发送消息给客户端)
	ResponseAction(response *PlayerActionResponse) ()

	// 来自其他玩家的动作
	NotifyFromOtherPlayerAction(notice *PlayerActionNotice) ()

	// 所有打牌动作,当用户请求时应该handle
	// playerDe 被操作者(点炮,点杠,碰等)
	DoAction(action *PlayerActionRequest, playerDe Player) (response *PlayerActionResponse,err error)

	// 序列化, 保存玩家状态(手上的牌,碰,杠,胡)
	Marshal() (bs []byte, err error)
	Unmarshal(bs []byte) (err error)
}

// 玩家动作请求
type PlayerActionRequest struct {
	Types ActionType
	Cards Cards // 动作哪几张牌, 比如亮倒隐藏刻子时有用
	Card  Card  // 动作哪张牌
}

// 玩家动作相应
type PlayerActionResponse struct {
	Types ActionType
	Card  Card  // 动作哪张牌
}

// 通知玩家来自其他玩家动作
type PlayerActionNotice struct {
	PlayerFrom Player
	Types      ActionType
	Card       Card // 动作哪张牌,有部分情况会为空 如摸牌
}

func NewActionResponse() *PlayerActionResponse {
	json.Marshal()
	json.Unmarshal()
	return &PlayerActionResponse{}
}
