package v1

import (
    "github.com/gorilla/mux"
    "github.com/ns3777k/mailcage/storage"
    "net/http"
)

type API struct {
    storage storage.Storage
}

type MessagesResponse struct {
    Total int
    Count int
    Start int
    Items []*storage.Message
}

func NewAPI(storage storage.Storage) *API {
    return &API{storage: storage}
}

func (a *API) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/messages", a.GetMessages).Methods("GET")
    router.HandleFunc("/message", a.GetMessage).Methods("GET")

    router.HandleFunc("/message", a.DeleteMessage).Methods("DELETE")
    router.HandleFunc("/messages", a.DeleteMessages).Methods("DELETE")
}

func (a *API) GetMessage(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    message, err := a.storage.GetOne(id)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "something bad happened")
        return
    }

    respondOk(w, message)
}

func (a *API) GetMessages(w http.ResponseWriter, r *http.Request) {
    start, limit := getPager(r)
    messages, err := a.storage.Get(start, limit)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "something bad happened")
        return
    }

    respondOk(w, &MessagesResponse{
        Total: a.storage.Count(),
        Count: len(messages),
        Items: messages,
        Start: start,
    })
}

func (a *API) DeleteMessage(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    if err := a.storage.DeleteOne(id); err != nil {
        respondError(w, http.StatusInternalServerError, "something bad happened")
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteMessages(w http.ResponseWriter, r *http.Request) {
    if err := a.storage.DeleteAll(); err != nil {
        respondError(w, http.StatusInternalServerError, "something bad happened")
        return
    }

    w.WriteHeader(http.StatusOK)
}
