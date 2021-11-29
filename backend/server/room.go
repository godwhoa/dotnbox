package server

import (
	"context"
	"dotnbox/dotnbox"
	"dotnbox/proto"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Room struct {
	ID        string
	createdat time.Time
	mutex     sync.Mutex
	conns     map[dotnbox.Owner]*websocket.Conn
	game      *dotnbox.Game
	log       *zap.Logger
}

func NewRoom(ID string, M, N int, log *zap.Logger) *Room {
	return &Room{
		ID:        ID,
		createdat: time.Now().UTC(),
		conns:     make(map[dotnbox.Owner]*websocket.Conn),
		game:      dotnbox.NewGame(M, N),
		log:       log,
	}
}

func stateToTurn(state dotnbox.State) dotnbox.Owner {
	switch state {
	case dotnbox.PlayerOneTurn:
		return dotnbox.PlayerOne
	case dotnbox.PlayerTwoTurn:
		return dotnbox.PlayerTwo
	default:
		return dotnbox.PlayerNone
	}

}

var ErrRoomFull = errors.New("Only two players per room")
var ErrGameOver = errors.New("Game over in this room")

func (r *Room) assignPlayer() dotnbox.Owner {
	player := dotnbox.PlayerOne
	if _, exists := r.conns[dotnbox.PlayerOne]; exists {
		player = dotnbox.PlayerTwo
	}
	return player
}

func (r *Room) IsFull() bool {
	return len(r.conns) >= 2
}

func (r *Room) CanBeDeleted() bool {
	return len(r.conns) == 0 && time.Since(r.createdat) > time.Hour*1
}

func (r *Room) HandleConn(ctx context.Context, conn *websocket.Conn) error {
	if len(r.conns) >= 2 {
		return ErrRoomFull
	}

	player := r.assignPlayer()
	r.conns[player] = conn
	defer delete(r.conns, player)

	// Send config
	wsjson.Write(ctx, conn, proto.ToPayload(proto.GameConfig{M: r.game.M, N: r.game.N, Player: player}))

	// Handle game start
	if len(r.conns) == 2 && r.game.State() == dotnbox.Waiting {
		r.log.Info("Starting game", zap.String("room", r.ID))
		r.game.Evaluate()
	}

	// Handle game resume
	if len(r.conns) == 2 && r.game.State() == dotnbox.Paused {
		r.log.Info("Resuming game", zap.String("room", r.ID))
		r.game.Resume()
	}
	r.BroadcastState(ctx)

	for {
		var payload proto.Payload
		err := wsjson.Read(ctx, conn, &payload)
		if err != nil {
			r.log.Error("Error reading payload", zap.Error(err))
			r.log.Info("Player disconnected", zap.String("room", r.ID), zap.Any("player", player))
			r.game.Pause() // Pause the game on connection drop
			delete(r.conns, player)
			r.BroadcastState(ctx)
			return err
		}

		if err := r.ProcessPayload(ctx, payload, player); err != nil {
			r.SendPlayer(ctx, player, proto.ToPayload(proto.Error{Error: err.Error()}))
			continue
		}

		r.game.Evaluate()
		r.BroadcastState(ctx)
	}
}

func (r *Room) ProcessPayload(ctx context.Context, payload proto.Payload, player dotnbox.Owner) error {
	switch payload.Type {
	case "PLACE":
		var line dotnbox.Line
		if err := json.Unmarshal(payload.Payload, &line); err != nil {
			r.log.Error("Error unmarshalling PLACE payload", zap.Error(err), zap.String("data", string(payload.Payload)))
			return err
		}
		if err := r.game.Place(line, player); err != nil {
			r.SendPlayer(ctx, player, proto.ToPayload(proto.Error{Error: err.Error()}))
			return err
		}
	case "REMATCH":
		if err := r.game.Rematch(); err != nil {
			r.SendPlayer(ctx, player, proto.ToPayload(proto.Error{Error: err.Error()}))
			return err
		}
		r.Broadcast(ctx, proto.ToPayload(proto.GameConfig{M: r.game.M, N: r.game.N, Player: player}))
	}
	return nil
}

func (r *Room) BroadcastState(ctx context.Context) {
	r.Broadcast(ctx, proto.ToPayload(proto.FromGame(r.game)))
}

func (r *Room) SendPlayer(ctx context.Context, player dotnbox.Owner, payload proto.Payload) {
	wsjson.Write(ctx, r.conns[player], payload)
}

func (r *Room) Broadcast(ctx context.Context, payload proto.Payload) {
	for _, conn := range r.conns {
		wsjson.Write(ctx, conn, payload)
	}
}
