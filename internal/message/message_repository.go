package message

import (
	"context"
	"database/sql"
	"log/slog"
)

type repository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &repository{db: db}
}

func (r *repository) createMessage(ctx context.Context, message *message) (*idRes, error) {
	var id idRes
	query := "INSERT INTO messages(message) VALUES($1) RETURNING id"
	if err := r.db.QueryRowContext(ctx, query, message.Message).Scan(&id.Id); err != nil {
		slog.Error("Error inserting message", "error", err)
		return nil, err
	}
	return &id, nil
}
