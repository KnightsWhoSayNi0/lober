create schema lober;

create table lober.teams (
    id          serial primary key,
    name        varchar(32) unique not null,
    color       varchar(6) not null default '000000'
);
create table lober.users (
    id          serial primary key,
    username    varchar(32) unique not null,
    password    varchar(64) not null,
    team_id     serial not null references teams /* todo; users can be on multiple teams */
);
create table lober.c2s (
    id          serial primary key,
    name        varchar(32) unique not null
);
create table lober.scope (
    id          serial primary key,
    name        varchar(32) not null,
    ip_addr     inet unique not null,
    team        int not null
);
create table lober.events (
    id          serial primary key,
    command     text not null,
    user_id     serial not null references users,
    c2_id       serial not null references c2s,
    scope_id    serial not null references scope,
    time        timestamp not null default current_timestamp
);
create table lober.tokens (
    id          serial primary key,
    prefix      varchar(3) not null unique,
    token       varchar not null unique,
    user_id     serial not null references users,
    c2_id       serial not null references c2s,
    created     timestamp not null default current_timestamp,
    expires     timestamp not null
);