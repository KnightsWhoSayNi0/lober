create table c2s (
                     id      bigint generated always as identity,
                     name    varchar(32)
);
create table scope(
                      id      bigint generated always as identity,
                      name    varchar(32)
);
create table logs(
                     id          bigint generated always as identity,
                     name        varchar(32),
                     command_run varchar(256),
                     c2_id       bigint,
                     scope_id    bigint
);
