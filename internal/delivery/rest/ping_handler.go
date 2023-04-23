package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/darahayes/go-boom"
	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go_ping_kube/internal/domain/errs"
	"go_ping_kube/internal/domain/models"
	"go_ping_kube/internal/domain/services"
	"net/http"
)

type PingHandler struct {
	service services.IPingService
}

func NewPingHandler(pingService *services.PingService) *PingHandler {
	return &PingHandler{
		service: pingService,
	}
}

func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	logrus.Info("PING !!!")

	data, err := h.service.Add(r.Context(), &models.CreatePingData{
		IP:         readUserIP(r),
		RequestURI: r.RequestURI,
		UserAgent:  r.UserAgent(),
		Headers:    r.Header,
	})
	if err != nil {
		boom.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	return
}

func (h *PingHandler) Get(w http.ResponseWriter, r *http.Request) {
	logrus.Info("GET !!!")

	if r.Method != http.MethodGet {
		boom.MethodNotAllowed(w)
		return
	}

	get := r.URL.Query().Get("uuid")
	logrus.Info("uuid =>", get)

	uuid, err := uuid2.Parse(get)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		//JSONError(w, err, http.StatusBadRequest)
		boom.BadRequest(w, fmt.Sprintf("param uuid is invalid. %s", err.Error()))
		return
	}

	data, err := h.service.Get(r.Context(), uuid)
	if err != nil {
		// TODO: Move to func like "handleDomainErr"

		rErr := &errs.DomainNotFoundError{}
		if errors.As(err, &rErr) {
			boom.NotFound(w, err.Error())
			return
		}

		boom.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (h *PingHandler) All(w http.ResponseWriter, r *http.Request) {
	logrus.Info("ALL !!!")
	list, err := h.service.All(r.Context())
	if err != nil {
		boom.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	return
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
