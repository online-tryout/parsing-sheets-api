-- name: CreateModule :one
INSERT INTO "modules" (
        title,
        "tryoutId",
        "moduleOrder"
    )
VALUES ($1, $2, $3)
RETURNING *;