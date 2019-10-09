package v1

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/ns3777k/mailcage/smtp"

	"github.com/ns3777k/mailcage/ws"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ns3777k/mailcage/storage"
)

type API struct {
	mailer   *smtp.Mailer
	storage  storage.Storage
	upgrader websocket.Upgrader
	wsHub    *ws.Hub
	logger   zerolog.Logger
}

type MessagesResponse struct {
	Total int
	Count int
	Start int
	Items []*storage.Message
}

func NewAPI(s storage.Storage, mailer *smtp.Mailer, logger zerolog.Logger) *API {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsHub := ws.NewHub()

	go func() {
		for e := range s.GetEvents() {
			wsHub.Broadcast(e)
		}
	}()

	return &API{
		storage:  s,
		upgrader: upgrader,
		wsHub:    wsHub,
		mailer:   mailer,
		logger:   logger.With().Str("api_version", "v1").Logger(),
	}
}

func (a *API) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/messages", a.GetMessages).Methods("GET")
	router.HandleFunc("/message", a.GetMessage).Methods("GET")
	router.HandleFunc("/ws", a.WebsocketUpgrade).Methods("GET")

	router.HandleFunc("/message", a.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/messages", a.DeleteMessages).Methods("DELETE")

	router.HandleFunc("/release", a.ReleaseMessage).Methods("POST")
}

func (a *API) GetMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	message, err := a.storage.GetOne(id)
	if err != nil {
		if err == storage.ErrMessageNotFound {
			respondError(w, http.StatusNotFound, "message not found")
			return
		}

		a.logger.Err(err).Str("handler", "GetMessage").Str("id", id).
			Msg("error getting a message from storage")

		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	respondOk(w, message)
}

func (a *API) GetMessages(w http.ResponseWriter, r *http.Request) {
	start, limit := getPager(r)
	messages, err := a.storage.Get(start, limit)
	if err != nil {
		a.logger.Err(err).Str("handler", "GetMessages").Msg("error getting messages from storage")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	cnt, err := a.storage.Count()
	if err != nil {
		a.logger.Err(err).Str("handler", "GetMessages").Msg("error counting messages")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	respondOk(w, &MessagesResponse{
		Total: cnt,
		Count: len(messages),
		Items: messages,
		Start: start,
	})
}

func (a *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := a.storage.DeleteOne(id); err != nil {
		if err == storage.ErrMessageNotFound {
			respondError(w, http.StatusNotFound, "message not found")
			return
		}

		a.logger.Err(err).Str("handler", "DeleteMessage").Str("id", id).
			Msg("error deleting a message")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteMessages(w http.ResponseWriter, r *http.Request) {
	if err := a.storage.DeleteAll(); err != nil {
		a.logger.Err(err).Str("handler", "DeleteMessages").Msg("error deleting messages")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) ReleaseMessage(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")
	id := r.URL.Query().Get("id")
	logger := a.logger.With().Str("id", id).
		Str("server", server).
		Str("handler", "ReleaseMessage").
		Logger()

	message, err := a.storage.GetOne(id)
	if err != nil {
		logger.Err(err).Msg("error getting a message from storage")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	if err := a.mailer.Send(server, message); err != nil {
		logger.Err(err).Msg("error releasing a message")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}
}

func (a *API) WebsocketUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Err(err).Str("handler", "WebsocketUpgrade").Msg("error upgrading connection")
		respondError(w, http.StatusInternalServerError, "something bad happened")
		return
	}

	client := ws.NewClient(a.wsHub, conn)

	go client.ReadPump()
	go client.WritePump()
}
