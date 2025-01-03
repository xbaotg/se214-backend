CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create type role as enum ('admin', 'user', 'lecturer');
create type day as enum ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday');
create type tu_status as enum ('paid', 'unpaid');
create type co_status as enum ('done', 'failed', 'progressing');
create table users (
	    id UUID primary key default uuid_generate_v4(),
	    username varchar(50) not null,
	    password text not null,
	    user_fullname text not null,
	    user_role role not null default 'user',
	    year int not null,
	    created_at timestamp not null default now(),
	    updated_at timestamp not null default now()
);

create table courses (
	    id UUID primary key default uuid_generate_v4(),
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
		confirmed boolean not null default false,
	    max_enroller int not null,
	    current_enroller int not null,
	    course_room text not null,
	    created_at timestamp not null default now(),
	    updated_at timestamp not null default now()
);

create table registered_courses (
	    course_id UUID not null,
	    user_id UUID not null,
	    course_year int not null,
	    course_semester int not null,
	    status co_status not null default 'progressing',
	    created_at timestamp not null default now(),
	    updated_at timestamp not null default now()
);

create table departments (
	    id UUID primary key default uuid_generate_v4(),
	    department_name text not null,
	    department_code text not null,
	    created_at timestamp not null default now(),
	    updated_at timestamp not null default now()
);

create table prerequisite_courses (
	    course_id text not null,
	    prerequisite_id text not null,
	    created_at timestamp not null default now(),
	    updated_at timestamp not null default now()
);

create table tuitions (
	    id UUID primary key default uuid_generate_v4(),
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

create table all_courses (
	course_name text primary key not null,
	course_fullname text not null,
	status boolean not null default true,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now()
);

alter table tuitions add constraint fk_user_id foreign key (user_id) references users (id);
alter table courses add constraint fk_course_user_id foreign key (course_teacher_id) references users (id);
alter table courses add constraint fk_department_id foreign key (department_id) references departments (id);
alter table registered_courses add constraint fk_course_id foreign key (course_id) references courses (id);
alter table registered_courses add constraint fk_resgiter_users_id foreign key (user_id) references users (id);
-- alter table prerequisite_courses add constraint fk_course_id foreign key (course_id) references courses (id);
-- alter table prerequisite_courses add constraint fk_prerequisite_id foreign key (prerequisite_id) references courses (id);
alter table courses add constraint fk_course_name foreign key (course_name) references all_courses (course_name);
alter table prerequisite_courses add constraint fk_prerequisite_courses_all_courses foreign key (course_id) references all_courses (course_name);
alter table prerequisite_courses add constraint fk_prerequisite_courses_prerequisite foreign key (prerequisite_id) references all_courses (course_name);

ALTER TABLE prerequisite_courses
ADD CONSTRAINT unique_course_prerequisite
UNIQUE (course_id, prerequisite_id);

alter table users add constraint unique_username unique (username);
alter table departments add constraint unique_department_name unique (department_code);

insert into users (username, password, user_fullname, user_role, year) values ('admin', '$2a$10$VFKP3WQvhZRb7CGVag1li.6DjtTKqp3tIoTpDLGPIY4pGQvwC1QXm', 'admin', 'admin', 0);