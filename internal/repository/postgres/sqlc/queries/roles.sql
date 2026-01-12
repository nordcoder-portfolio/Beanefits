-- name: GetRolesByUserID :many
SELECT role_code
FROM user_roles
WHERE user_id = $1
ORDER BY role_code;

-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_code)
VALUES ($1, $2);

-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_code = $2;
