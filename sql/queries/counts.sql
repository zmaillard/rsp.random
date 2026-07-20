-- name: GetSigns :many
SELECT imageid::text, country_slug, state_slug, place_slug, county_slug, has_processed FROM sign.vwhugohighwaysign;