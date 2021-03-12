create table users (
    id serial primary key,
    uname varchar(255) not null,
    email varchar(255) not null,
    pword varchar(255) not null,
    created_at timestamp not null
);

create table groups (
    id serial primary key,
    name varchar(255)
);

create table memberships (
    joined_at timestamp not null,
    user_id integer references users(id) on delete cascade,
    group_id integer references groups(id) on delete cascade,
    PRIMARY KEY (user_id, group_id)
);

create table roles (
    id serial primary key,
    name varchar(255) not null,
    permissions integer,
    gets_chores boolean,
    group_id integer references groups(id) ON DELETE CASCADE
);

create table role_assignments (
    role_id integer references roles(id) ON DELETE CASCADE,
    user_id integer references users(id) ON DELETE CASCADE 
);

create table chores (
    id serial primary key,
    description varchar(255),
    name varchar (255) not null,
    duration integer,
    group_id integer references groups(id) ON DELETE CASCADE
);

create table chore_assignments (
    complete boolean,
    date_assigned timestamp not null,
    date_complete timestamp,
    date_due timestamp,
    chore_id integer references chores(id) ON DELETE CASCADE,
    user_id integer references users(id) ON DELETE CASCADE,
    PRIMARY KEY (chore_id, user_id)
);

create table sessions (
    uuid  not null primary key,
    values varchar,
    created timestamp not null,
    user_id integer references users(id) ON DELETE CASCADE
);