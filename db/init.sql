create schema lober;

create table lober.c2s (
                     id      serial primary key,
                     name    varchar(32)
);
create table lober.scope (
                      id      serial primary key,
                      name    varchar(32)
);
create table lober.events (
                     id          serial primary key,
                     name        varchar(32),
                     command_run varchar(256),
                     c2_id       serial,
                     scope_id    serial,
                     time        timestamp default current_timestamp
);