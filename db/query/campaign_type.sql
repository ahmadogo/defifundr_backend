-- name: CreateCampaignType :one

INSERT INTO
    campaign_types (campaign_types)
VALUES ($1) RETURNING *;

-- name: GetAllCampaignType :many

SELECT * FROM campaign_types;