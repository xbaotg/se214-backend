CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create type role as enum ('admin', 'user');
create type day as enum ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday');
create type tu_status as enum ('paid', 'unpaid');
create type co_status as enum ('done', 'failed', 'progressing');

create table users (
    user_id UUID primary key default uuid_generate_v4(),
    username varchar(50) not null,
    password text not null,
    user_email text not null,
    user_fullname text not null,
    user_role role not null default 'user',
    year int not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table courses (
    course_id UUID primary key default uuid_generate_v4(),
    course_teacher_id UUID not null,
    department_id UUID not null,
    course_name text not null,
    course_fullname text not null,
    course_credit int not null,
    course_year int not null,
    course_semester int not null,
    course_start_shift int not null,
    course_end_shift int not null,
    course_day day not null,
    max_enroller int not null,
    current_enroller int not null,
    course_room text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table registed_courses (
    course_id UUID not null,
    user_id UUID not null,
    course_year int not null,
    course_semester int not null,
    status co_status not null default 'progressing',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table departments (
    department_id UUID primary key default uuid_generate_v4(),
    department_name text not null,
    department_code text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table prerequisite_courses (
    course_id UUID not null,
    prerequisite_id UUID not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table tuition (
    tuition_id UUID primary key default uuid_generate_v4(),
    user_id UUID not null,
    tuition int not null,
    paid int not null default 0,
    total_credit int not null,
    year int not null,
    semester int not null,
    tuition_status tu_status not null default 'unpaid',
    tuition_deadline timestamp not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table sessions (
    session_id UUID primary key default uuid_generate_v4(),
    user_id UUID not null,
    refresh_token text not null,
    expires_in timestamp not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    is_active boolean not null default true
);

alter table tuition add constraint fk_user_id foreign key (user_id) references users (user_id);
alter table courses add constraint fk_user_id foreign key (course_teacher_id) references users (user_id);
alter table courses add constraint fk_department_id foreign key (department_id) references departments (department_id);
alter table registed_courses add constraint fk_course_id foreign key (course_id) references courses (course_id);
alter table registed_courses add constraint fk_user_id foreign key (user_id) references users (user_id);
alter table prerequisite_courses add constraint fk_course_id foreign key (course_id) references courses (course_id);
alter table prerequisite_courses add constraint fk_prerequisite_id foreign key (prerequisite_id) references courses (course_id);

