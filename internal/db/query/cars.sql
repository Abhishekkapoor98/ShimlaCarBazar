-- name: CreateCar :one
INSERT INTO cars (title, brand, model, condition_type, price, status, image_urls)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListCars :many
SELECT * FROM cars
ORDER BY created_at DESC;