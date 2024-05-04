-- name: CreateQuestion :one
INSERT INTO "questions" (
        content,
        "moduleId",
        "questionOrder"
    )
VALUES ($1, $2, $3)
RETURNING *;