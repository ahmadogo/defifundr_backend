-- name: CreateCampaignType :one

INSERT INTO
    campaigns (campaign_name)
VALUES ($1) RETURNING *;

-- name: GetAllCampaignType :many

SELECT * FROM campaigns;