-- name: CreateSellProposal :one
INSERT INTO sell_proposals (customer_id, car_brand, car_model, expected_price, phone_number, image_urls, status)
VALUES ($1, $2, $3, $4, $5, $6, 'pending')
RETURNING *;

-- name: ListSellProposals :many
SELECT * FROM sell_proposals
ORDER BY created_at DESC;

-- name: GetSellProposal :one
SELECT * FROM sell_proposals
WHERE id = $1 LIMIT 1;

-- name: UpdateSellProposalStatus :one
UPDATE sell_proposals 
SET status = $2 
WHERE id = $1 
RETURNING *;