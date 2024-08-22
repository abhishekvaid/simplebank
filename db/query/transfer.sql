-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES($1, $2, $3) RETURNING *; 

-- name: GetTransferById :one
SELECT * FROM transfers 
WHERE id = $1 LIMIT 1; 

-- name: GetTransfersByFromAccountId :many
SELECT * FROM transfers 
WHERE from_account_id = $1 LIMIT $2 OFFSET $3; 

-- name: GetTransfersByToAccountId :many
SELECT * FROM transfers 
WHERE to_account_id = $1 LIMIT $2 OFFSET $3;

-- name: GetTransfersByFromAndToAccountId :many
SELECT * FROM transfers 
WHERE from_account_id = $1 AND to_account_id = $2 
LIMIT $3 OFFSET $4;