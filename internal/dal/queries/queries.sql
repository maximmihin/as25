-- name: GetDateFilteredPVZWithReceptionsAndProducts :many
WITH filtered_receptions AS (
    SELECT
        r.id,
        r.created_at,
        r.pvz_id,
        r.status
    FROM
        receptions r
    WHERE
        r.created_at BETWEEN @start_date AND @end_date
)
SELECT
    p.id AS pvz_id,
    p.created_at AS pvz_created_at,
    p.located_at AS pvz_city,
    COALESCE(
        json_agg(
            json_build_object(
                'reception_id', fr.id,
                'reception_created_at', fr.created_at,
                'reception_status', fr.status
            ) ORDER BY fr.created_at                    -- todo not pvz_created_at?
        ) FILTER (WHERE fr.id IS NOT NULL),
        '[]'::json
    ) AS receptions
FROM pvz p
    LEFT JOIN filtered_receptions fr ON p.id = fr.pvz_id
GROUP BY
    p.id
ORDER BY
    p.created_at
LIMIT @ft_limit
OFFSET @ft_offset
;