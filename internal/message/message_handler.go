package message

import (
	"encoding/json"
	"fmt"
	"messagio_test_task/internal/models"
	"messagio_test_task/internal/producer"
	"messagio_test_task/internal/responses"
	"net/http"
)

type MessageHandler struct {
	MessageService
	*producer.Producer
}

func NewMessageHandler(service MessageService, producer *producer.Producer) *MessageHandler {
	return &MessageHandler{MessageService: service, Producer: producer}
}

// СreateMessage godoc
// @Summary Создает новое сообщение
// @Tags message
// @Accept json
// @Param message body messageReq true "Сообщение"
// @Success 201
// @Failure 400 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /messages/create [post]
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg messageReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&msg); err != nil {
		response(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &models.Message{Message: msg.Message}
	if err := h.MessageService.createMessage(r.Context(), message); err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Producer.SendMessage(message); err != nil {
		response(w, http.StatusInternalServerError, "Failed to send message to Kafka")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPendingMessagesCount godoc
// @Summary Выводит информацию о количестве необработанных сообщений
// @Tags message
// @Accept json
// @Success 200 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /messages/stats [get]
func (h *MessageHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	pendingCount, err := h.MessageService.getPendingMessagesCount(r.Context())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	processedCount, err := h.MessageService.getProcessedMessagesCount(r.Context())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusCreated, fmt.Sprintf("Pending/Processed messages: %d/%d", *pendingCount, *processedCount))
}

func response(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responses.Response{Message: message})
}
