-- name: CreateOption :one
INSERT INTO "options" (
        content,
        "questionId",
        "isTrue",
        "optionOrder"
    )
VALUES ($1, $2, $3, $4)
RETURNING *;