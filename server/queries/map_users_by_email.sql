SELECT
  id AS user_id,
  LOWER(email) AS email
FROM users
WHERE LOWER(email) = ANY($1::text[]);
