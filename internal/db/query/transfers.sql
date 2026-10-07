-- name: InitiateRCTransfer :one
INSERT INTO rc_transfers (car_id, buyer_id, current_status)
VALUES ($1, $2, 'documents_submitted')
RETURNING *;

-- name: ListRCTransfers :many
SELECT * FROM rc_transfers
ORDER BY created_at DESC;

-- name: UpdateTransferStatus :one
UPDATE rc_transfers 
SET current_status = $2 
WHERE id = $1 
RETURNING *;