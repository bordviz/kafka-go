package message

import (
	"go-kafka/internal/config"
	"go-kafka/internal/handlers"
	"go-kafka/internal/lib/logger/sl"
	"go-kafka/internal/services/producer"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	log      *slog.Logger
	producer *producer.ProducerService
	cfg      config.Kafka
}

func NewHandler(producer *producer.ProducerService, log *slog.Logger, cfg config.Kafka) *Handler {
	return &Handler{producer: producer, log: log, cfg: cfg}
}

func AddMessageHandler(log *slog.Logger, cfg config.Kafka, producer *producer.ProducerService) func(r chi.Router) {
	handler := NewHandler(producer, log, cfg)

	return func(r chi.Router) {
		r.Post("/", handler.Send())
	}
}

type Message struct {
	Message string `json:"message"`
}

func (h *Handler) Send() http.HandlerFunc {
	const op = "handlers.message.Send"

	return func(w http.ResponseWriter, r *http.Request) {
		h.log = slog.With(
			slog.String("op", op),
		)

		var message Message
		if err := render.Decode(r, &message); err != nil {
			h.log.Error("failed to decode request body", sl.Err(err))
			handlers.ErrorResponse(w, r, 400, "failed to decode body")
			return
		}
		if message.Message == "" {
			handlers.ErrorResponse(w, r, 422, "message is required")
			return
		}

		err := h.producer.ProduceMessage(message.Message, h.cfg.MessageTopic)
		if err != nil {
			h.log.Error("failed to send message to topic", sl.Err(err))
			handlers.ErrorResponse(w, r, 400, "failed to send message to topic")
			return
		}

		h.log.Info("message successfully received")
		handlers.SuccessResponse(w, r, 200, "message sent successfully")
	}
}
