// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: module.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createModule = `-- name: CreateModule :one
INSERT INTO "modules" (
        title,
        "tryoutId",
        "moduleOrder"
    )
VALUES ($1, $2, $3)
RETURNING id, title, "tryoutId", "moduleOrder", "updatedAt", "createdAt"
`

type CreateModuleParams struct {
	Title       string        `json:"title"`
	TryoutId    uuid.UUID     `json:"tryoutId"`
	ModuleOrder sql.NullInt32 `json:"moduleOrder"`
}

func (q *Queries) CreateModule(ctx context.Context, arg CreateModuleParams) (Modules, error) {
	row := q.db.QueryRowContext(ctx, createModule, arg.Title, arg.TryoutId, arg.ModuleOrder)
	var i Modules
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.TryoutId,
		&i.ModuleOrder,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
