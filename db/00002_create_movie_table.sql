-- +goose Up
-- +goose StatementBegin
create table if not exists movie(
    id integer not null,
    title text not null,
    director text,
    release_year integer,
    genre text,
    rating real,
    description text,
    duration text,

    check (length(title) >= 2),

    primary key(id) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists movie;
-- +goose StatementEnd
