create schema lober;

create table lober.users (
    id          serial primary key,
    username    varchar(32),
    email       varchar(32),
    team_id     serial
);
create table lober.teams (
    id          serial primary key,
    name        varchar(32),
    color       varchar(6),
    lead_id     serial
);
create table lober.c2s (
    id          serial primary key,
    name        varchar(32)
);
create table lober.scope (
    id          serial primary key,
    name        varchar(32)
);
create table lober.events (
    id          serial primary key,
    command     varchar(256),
    user_id     serial,
    c2_id       serial,
    scope_id    serial,
    time        timestamp default current_timestamp
);