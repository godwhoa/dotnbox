package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

type Server struct {
	mutex      sync.RWMutex
	rooms      map[string]*Room
	maxRooms   int
	log        *zap.Logger
	gcInterval time.Duration
}

func (s *Server) garbageCollector(ctx context.Context) {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.mutex.Lock()
			for id, room := range s.rooms {
				if room.CanBeDeleted() {
					delete(s.rooms, id)
				}
			}
			s.mutex.Unlock()
		}
	}
}

func New(log *zap.Logger) *Server {
	return &Server{
		mutex:      sync.RWMutex{},
		rooms:      make(map[string]*Room),
		maxRooms:   100,
		log:        log,
		gcInterval: time.Minute * 10,
	}
}

type RoomCreationRequest struct {
	M int `json:"m"`
	N int `json:"n"`
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	s.mutex.RLock()
	_, exists := s.rooms[roomID]
	s.mutex.RUnlock()
	if exists {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}

	if len(s.rooms) >= s.maxRooms {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var req RoomCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	s.rooms[roomID] = NewRoom(roomID, req.M, req.N, s.log)
	s.mutex.Unlock()
	s.log.Info("Created room", zap.String("roomID", roomID), zap.String("addr", r.RemoteAddr))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) Run(ctx context.Context, addr string) error {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/room/{roomID}", s.handleRoom)
	r.Post("/room/{roomID}", s.handleCreateRoom)

	go s.garbageCollector(ctx)

	s.log.Info("Starting server", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		s.log.Error("Failed to start server", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) handleRoom(w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		s.log.Error("Failed to accept WS connection", zap.Error(err))
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Sorry bob")

	roomID := chi.URLParam(r, "roomID")
	if _, ok := s.rooms[roomID]; !ok {
		s.log.Info("Closing connection, room not found", zap.String("roomID", roomID), zap.String("addr", r.RemoteAddr))
		conn.Close(websocket.StatusPolicyViolation, "Room does not exist")
		return
	}

	s.log.Info("Accepted WS connection", zap.String("addr", r.RemoteAddr), zap.String("roomID", roomID))

	room := s.rooms[roomID]
	err = room.HandleConn(r.Context(), conn)
	s.log.Info("room.HandleCoon returned", zap.String("roomID", roomID), zap.String("addr", r.RemoteAddr), zap.Any("err", err))
	switch err {
	case ErrGameOver:
		conn.Close(websocket.StatusGoingAway, "Game over")
		return
	case ErrRoomFull:
		conn.Close(websocket.StatusPolicyViolation, "Room full")
		return
	case nil:
		conn.Close(websocket.StatusNormalClosure, "")
		return
	default:
		s.log.Error("Failed to handle WS connection", zap.Error(err))
	}
}
