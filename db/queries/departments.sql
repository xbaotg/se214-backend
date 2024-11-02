-- create table departments (
--     department_id UUID primary key default uuid_generate_v4(),
--     department_name text not null,
--     department_code text not null,
--     created_at timestamp not null default now(),
--     updated_at timestamp not null default now()
-- );


-- name: CreateDepartment :one
INSERT INTO departments (
    department_name, department_code, created_at, updated_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListDepartments :many
SELECT * FROM departments;
 