-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
    FROM feed_follows
    LEFT JOIN feeds
    ON feed_follows.feed_id = feeds.id
    LEFT JOIN posts
    ON feeds.id = posts.feed_id 
    WHERE feed_follows.user_id = $1
    ORDER BY posts.updated_at DESC
    LIMIT $2;
