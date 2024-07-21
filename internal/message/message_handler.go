package message

import (
	"encoding/json"
	"io"
	"messagio_test_task/internal/responses"
	"net/http"
)

type MessageHandler struct {
	MessageService
}

func NewMessageHandler(service MessageService) *MessageHandler {
	return &MessageHandler{MessageService: service}
}

// СreateMessage godoc
// @Summary Создает новое сообщение
// @Tags message
// @Accept json
// @Param message body message true "Сообщение"
// @Success 201
// @Failure 400 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /messages/create [post]
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg message

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			response(w, http.StatusBadRequest, err.Error())
			return
		}
	}(r.Body)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&msg); err != nil {
		response(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.MessageService.createMessage(r.Context(), &msg)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func response(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responses.Response{Message: message})
}
