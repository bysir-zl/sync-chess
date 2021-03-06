package chess

type ActionType uint16
type ActionTypes []ActionType

// 定义了基本打牌动作, 如有其它的可再定义
const (
	AT_Get         ActionType = iota + 1 // 摸牌, 服务器不会下发这个命令 而是自动(比如杠牌后)下发通知告知玩家摸了那张牌
	AT_Play                              // 出牌
	AT_Peng                              // 碰
	AT_GangDian                          // 直杠
	AT_GangAn                            // 暗杠
	AT_GangBu                            // 补杠
	AT_HuDian                            // 点炮
	AT_HuZiMo                            // 自摸
	AT_HuQiangGang                       // 抢杠胡
	AT_LiangDao                          // 亮倒
	AT_Pass                              // 过, 可以过 杠,碰,胡
)

func (p ActionType) String() (s string) {
	switch p {
	case AT_Get:
		s = "Get"
	case AT_Play:
		s = "Play"
	case AT_Peng:
		s = "Pong"
	case AT_GangDian:
		s = "GangDian"
	case AT_GangBu:
		s = "GangBu"
	case AT_GangAn:
		s = "GangAn"
	case AT_HuDian:
		s = "HuDian"
	case AT_HuZiMo:
		s = "HuZiMo"
	case AT_HuQiangGang:
		s = "QiangGang"
	case AT_LiangDao:
		s = "LiangDao"
	case AT_Pass:
		s = "Pass"
	}
	return
}

func (p *ActionTypes) Contain(a ActionType) bool {
	for _, at := range *p {
		if at == a {
			return true
		}
	}
	return false
}

// 动作来至哪
type ActionFrom int32

const (
	AF_Auto    ActionFrom = iota + 1 // 自动打牌
	AF_Player                        // 来至玩家
	AF_Storage                       // 来至存档
)

func (p ActionFrom) String() string {
	s := ""
	switch p {
	case AF_Auto:
		s = "Auto"
	case AF_Player:
		s = "Player"
	case AF_Storage:
		s = "storage"
	}
	return s
}

// 玩家动作请求
type PlayerActionRequest struct {
	Types      ActionType `json:"Types"`
	Cards      Cards `json:"Cards"` // 动作哪几张牌, 比如亮倒隐藏刻子时有用
	Card       Card  `json:"Card"`  // 动作哪张牌
	ActionFrom ActionFrom
}
