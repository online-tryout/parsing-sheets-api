-- name: CreateTryout :one
INSERT INTO "tryouts" (
        title,
        price,
        status,
        "startedAt",
        "endedAt"
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;