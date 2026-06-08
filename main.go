package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/heroiclabs/nakama-common/runtime"
)

type registerCheckRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type registerCheckResponse struct {
	EmailExists    bool `json:"email_exists"`
	UsernameExists bool `json:"username_exists"`
	Available      bool `json:"available"`
}

type pingResponse struct {
	Message string `json:"message"`
	Payload string `json:"payload,omitempty"`
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := initializer.RegisterRpc("ping", rpcPing); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("register_check", rpcRegisterCheck); err != nil {
		return err
	}

	logger.Info("Registered RPC: ping")
	logger.Info("Registered RPC: register_check")
	return nil
}

func rpcPing(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	response, err := json.Marshal(pingResponse{
		Message: "pong",
		Payload: payload,
	})
	if err != nil {
		logger.Error("Failed to marshal ping response: %v", err)
		return "", err
	}

	return string(response), nil
}

func rpcRegisterCheck(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request registerCheckRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.Error("Failed to unmarshal register_check request: %v", err)
		return "", runtime.NewError("Invalid JSON payload.", 3)
	}

	email := strings.TrimSpace(request.Email)
	username := strings.TrimSpace(request.Username)
	if email == "" && username == "" {
		return "", runtime.NewError("Email or username is required.", 3)
	}

	var response registerCheckResponse
	err := db.QueryRowContext(ctx, `
SELECT
	CASE WHEN $1::text = '' THEN false ELSE EXISTS (
		SELECT 1 FROM users WHERE lower(email) = lower($1::text)
	) END AS email_exists,
	CASE WHEN $2::text = '' THEN false ELSE EXISTS (
		SELECT 1 FROM users WHERE username = $2::text
	) END AS username_exists
`, email, username).Scan(&response.EmailExists, &response.UsernameExists)
	if err != nil {
		logger.Error("Failed to check register duplicates: %v", err)
		return "", runtime.NewError("Failed to check register duplicates.", 13)
	}

	response.Available = !response.EmailExists && !response.UsernameExists

	result, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal register_check response: %v", err)
		return "", err
	}

	return string(result), nil
}
