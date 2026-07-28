package api

import (
	"encoding/json"
	"html/template"
	"message-feed/internal/models"
	pb "message-feed/proto"
	"net/http"

	coreactor "github.com/TSachin36/message-feed-core/actor"
)

func toModels(pbMessages []*pb.Message) []models.Message {

	messages := make(
		[]models.Message,
		0,
		len(pbMessages),
	)

	for _, msg := range pbMessages {

		messages = append(
			messages,
			models.Message{
				UserID: msg.UserID,
				Text:   msg.Text,
			},
		)
	}

	return messages
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		userID := r.URL.Query().Get("user")

		if userID == "" {
			http.Error(
				w,
				"Missing user",
				http.StatusBadRequest,
			)
			return
		}

		messages, err := getMessages(
			r.Context(),
			userID,
		)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

	case http.MethodPost:

		defer r.Body.Close()
		var msg models.Message

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(
				w,
				"Invalid JSON",
				http.StatusBadRequest,
			)
			return
		}

		err = saveMessage(
			r.Context(),
			msg,
		)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		coreactor.Broadcast(
			msg.UserID,
			msg.Text,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(msg)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

	default:

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user")

	if userID == "" {
		http.Error(
			w,
			"Missing user",
			http.StatusBadRequest,
		)
		return
	}

	messages, err := getMessages(
		r.Context(),
		userID,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/messages.html",
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = tmpl.Execute(w, messages)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}
