package chess

import "github.com/bysir-zl/sync-chess/core"

type PlayerCards struct {
	Cards core.Cards         // 手上的牌
	Pong  core.Cards         // 碰的牌
	Gang  map[core.Card]Gang // 杠的牌
}

type Gang struct {
	Score    int32         // 分数, 杠需要记录扣分的人. 杠上杠的情况分数不一样
	Receiver core.Player   // 接收者
	Giver    []core.Player // 给予者
	Types    core.GangType
}