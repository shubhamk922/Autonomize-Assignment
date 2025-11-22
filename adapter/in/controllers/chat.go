package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/team-monitoring/service"
)

type ChatController struct {
	Service *service.ChatBot
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func (controller *ChatController) Handle(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := controller.Service.Handle(context.Background(), req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("bot error: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ChatResponse{Response: resp})
}
