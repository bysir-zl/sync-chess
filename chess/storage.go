package chess

import (
	"github.com/bysir-zl/bygo/log"
	"bytes"
	"encoding/json"
	"errors"
)

// 用于保存牌局, 实现down机恢复, 重放等功能
type Storage struct {
	manager       *Manager
	playerActionC map[string]chan *PlayerActionRequest // 玩家记录存档
}

type SnapShoot struct {
	Players            Players
	RoundStartPlayerId string
	SurplusCards       Cards
}

type Step struct {
	PlayerId      string
	ActionRequest *PlayerActionRequest
}

// 保存每轮开始(出牌玩家开始出牌之前)快照
func (p *Storage) SnapShoot() {
	players := p.manager.Players
	roundStartPlayer := p.manager.RoundStartPlayer
	surplusCards := p.manager.CardGenerator.GetCardsSurplus()

	s := SnapShoot{
		Players:            players,
		RoundStartPlayerId: roundStartPlayer.Id,
		SurplusCards:       surplusCards,
	}
	sBs, err := s.M()
	if err != nil {
		panic(err)
	}
	bs := make([]byte, len(sBs)+1)
	bs[0] = 1
	copy(bs[1:], sBs)

	Redis.RPUSH("STO", bs)

	log.Info("Storage SnapShoot ", s)
}

// 恢复快照,并且读取待运行的操作
func (p *Storage) Recovery() (has bool) {
	// 拉取所有记录, 找到最近一次SnapShoot并恢复
	// 将Step记录转换为manager.AutoAction, 当有AutoAction时manager不在询问player而是自动action
	ss, err := Redis.LRANGE("STO", 0, -1)
	if err != nil {
		log.Error("Recovery ERR: ", err)
		return
	}

	steps := []*Step{}
	for i := len(ss) - 1; i >= 0; i-- {
		bs := ss[i].([]byte)
		switch bs[0] {
		case 1:
			snap := SnapShoot{}
			err := snap.UnM(bs[1:], p.manager.PlayerLeader.PlayerCreator)
			if err != nil {
				log.Error("snap.UnM Err:", err)
				return
			}
			p.manager.Players = snap.Players
			p.manager.RoundStartPlayer, _ = p.manager.Players.Find(snap.RoundStartPlayerId)
			p.manager.CardGenerator.SetCardsSurplus(snap.SurplusCards)

			log.Info("Storage Recovery", snap)
			has = true
		case 2:
			step := &Step{}
			err := step.UnM(bs[1:])
			if err != nil {
				log.Error("step.UnM Err:", err)
				return
			}
			steps = append(steps, step)
		}
		if has {
			break
		}
	}

	// 处理step
	for _, step := range steps {
		if _, ok := p.playerActionC[step.PlayerId]; !ok {
			p.playerActionC[step.PlayerId] = make(chan *PlayerActionRequest, 100)
		}
		p.playerActionC[step.PlayerId] <- step.ActionRequest
		log.Info("Storage Recovery Step", step)
	}

	return
}

// 保存玩家操作日志
func (p *Storage) Step(player *Player, request *PlayerActionRequest) {
	s := Step{
		PlayerId:      player.Id,
		ActionRequest: request,
	}
	sBs, err := json.Marshal(&s)
	if err != nil {
		panic(err)
	}
	bs := make([]byte, len(sBs)+1)
	bs[0] = 2
	copy(bs[1:], sBs)

	Redis.RPUSH("STO", bs)

	log.Info("Storage Step ", player, request)
}

// 清空这局存档
func (p *Storage) Clean() {
	Redis.DEL("STO")

	log.Info("Storage Cleaned")
	return
}

// 获取玩家动作存档
func (p *Storage) HasStep(playerId string) (has bool) {
	c, ok := p.playerActionC[playerId]
	if !ok {
		return
	}
	has = len(c) != 0
	return
}

// 获取玩家动作存档
func (p *Storage) PopStep(playerId string) (action *PlayerActionRequest, has bool) {
	c, ok := p.playerActionC[playerId]
	if !ok {
		return
	}
	select {
	case action = <-c:
		has = true
		return
	default:
		return
	}
	return
}

// ---------------

var sp = []byte("@#$%$#@")
var spPlayer = []byte("^&*(*&^")
var spPlayerId = []byte("^&ID&^")

func (s *SnapShoot) M() (bs []byte, err error) {
	var buff bytes.Buffer
	playerBs := [][]byte{}
	for _, player := range s.Players {
		pbs, e := player.PlayerI.Marshal()
		if e != nil {
			err = e
			return
		}
		pbs = append(pbs, spPlayerId...)
		// 添加id
		pbs = append(pbs, []byte(player.Id)...)
		playerBs = append(playerBs, pbs)
	}
	// 写入玩家
	buff.Write(bytes.Join(playerBs, spPlayer))
	// 写入局头人
	buff.Write(sp)
	buff.Write([]byte(s.RoundStartPlayerId))
	// 写入剩余卡牌
	buff.Write(sp)
	cardsBs, err := json.Marshal(s.SurplusCards)
	if err != nil {
		return
	}
	buff.Write(cardsBs)

	bs = buff.Bytes()

	return
}

func (s *SnapShoot) UnM(bs []byte, PlayerCreator func() PlayerInterface) (err error) {
	bsp := bytes.Split(bs, sp)
	if len(bsp) != 3 {
		err = errors.New("bad format")
		return
	}
	// 读取玩家
	if len(bsp[0]) == 0 {
		err = errors.New("bad format: player")
		return
	}
	playerBs := bytes.Split(bsp[0], spPlayer)
	if len(playerBs) == 0 {
		err = errors.New("bad format: player len is zero")
		return
	}
	lenPlayer := len(playerBs)
	players := make(Players, lenPlayer)
	for i, pbs := range playerBs {
		idAI := bytes.Split(pbs, spPlayerId)
		if len(idAI) != 2 {
			err = errors.New("bad format: playerId")
			return
		}
		id := string(idAI[1])
		player := NewPlayer(id, PlayerCreator)
		err = player.PlayerI.Unmarshal(idAI[0])
		if err != nil {
			return
		}
		players[i] = player
	}
	s.Players = players
	// 读取id
	s.RoundStartPlayerId = string(bsp[1])
	// 剩余卡牌
	err = json.Unmarshal(bsp[2], &s.SurplusCards)
	if err != nil {
		return
	}

	return
}
func (s *Step) M() (bs []byte, err error) {
	bs, err = json.Marshal(s)
	return
}

func (s *Step) UnM(bs []byte) (err error) {
	err = json.Unmarshal(bs, s)
	return
}

func NewStorage(manager *Manager) *Storage {
	return &Storage{
		manager:       manager,
		playerActionC: map[string]chan *PlayerActionRequest{},
	}
}