-- +goose Up
-- +goose StatementBegin
create table if not exists user(
    id integer not null,
    username text not null unique,
    password text not null,

    check(length(username) >= 5),

    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user;
-- +goose StatementEnd
