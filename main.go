package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type pingResponse struct {
	Message string `json:"message"`
	Payload string `json:"payload,omitempty"`
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := initializer.RegisterRpc("ping", rpcPing); err != nil {
		return err
	}

	logger.Info("Registered RPC: ping")
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
