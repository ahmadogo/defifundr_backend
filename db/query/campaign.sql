-- name: CreateCampaignType :one

INSERT INTO
    campaigns (campaign_name, image)
VALUES ($1, $2) RETURNING *;

-- name: GetAllCampaignType :many

SELECT * FROM campaigns;
