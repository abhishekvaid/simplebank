-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES( $1, $2, $3)
RETURNING *; 

-- name: GetAccount :one
SELECT * FROM accounts 
WHERE id = $1 AND owner = $2
LIMIT 1 ; 

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
WHERE id = $1 AND owner = $2 
LIMIT 1
FOR NO KEY UPDATE ;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $3
WHERE id = $1 and owner = $2
RETURNING * ; 

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg('amount')
WHERE id = sqlc.arg('id') and owner = sqlc.arg('owner')
RETURNING * ; 

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1 and owner = $2;