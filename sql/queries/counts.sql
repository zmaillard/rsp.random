-- name: GetStateCounts :many
select aas.slug,count(*) as counter from sign.highwaysign
inner join sign.admin_area_state aas on aas.id = highwaysign.admin_area_state_id
group by aas.slug;

-- name: GetCountyCounts :many
select (s.slug || '_' || c.slug)::text as county_slug, county_counts.counter
from sign.admin_area_county c inner join sign.admin_area_state s
                                         on c.admin_area_stateid = s.id
                              inner join (
    select aacn.id,count(*) as counter from sign.highwaysign
                                                inner join sign.admin_area_county aacn on aacn.id = highwaysign.admin_area_county_id
    group by aacn.id) county_counts on c.id = county_counts.id;

-- name: GetPlaceCounts :many
select (s.slug || '_' || p.slug)::text as place_slug, place_counts.counter
from sign.admin_area_place p inner join sign.admin_area_state s
                                        on p.admin_area_stateid = s.id
                             inner join (
    select aap.id,count(*) as counter from sign.highwaysign
                                               inner join sign.admin_area_place aap on aap.id = highwaysign.admin_area_place_id
    group by aap.id) place_counts on p.id = place_counts.id;

-- name: GetTotalCounts :one
select count(*) from sign.highwaysign;

-- name: GetSigns :many
SELECT imageid::text, country_slug, state_slug, place_slug, county_slug FROM sign.vwhugohighwaysign;