-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

create table if not exists room_movie(
    id integer not null,
    room_id integer not null,
    movie_id integer not null,
    start_date_unix integer not null,

    foreign key(room_id) references room(id) on delete cascade,
    foreign key(movie_id) references movie(id) on delete cascade,
    primary key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists room_movie;
-- +goose StatementEnd
