-- create table courses (
--     course_id UUID primary key default uuid_generate_v4(),
--     course_teacher_id UUID not null,
--     department_id UUID not null,
--     course_name text not null,
--     course_fullname text not null,
--     course_credit int not null,
--     course_year int not null,
--     course_semester int not null,
--     course_start_shift int not null,
--     course_end_shift int not null,
--     course_day day not null,
--     max_enroller int not null,
--     current_enroller int not null,
--     course_room text not null,
--     created_at timestamp not null default now(),
--     updated_at timestamp not null default now()
-- );

-- name: CreateCourse :one
INSERT INTO courses (
    course_teacher_id, department_id, course_name, course_fullname, course_credit, course_year, course_semester, course_start_shift, course_end_shift, course_day, max_enroller, current_enroller, course_room, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;