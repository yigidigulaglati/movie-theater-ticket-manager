-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

create table if not exists room_movie_seating(
    id integer not null,
    room_movie_id integer not null,

    row_index integer not null,
    col_index integer not null,
    seat integer not null,

    check(row_index >= 0),
    check(col_index >= 0),
    check(seat = 1 or seat = 2 or seat = -1 or seat -2),

    foreign key(room_movie_id) references room_movie(id) on delete cascade,
    primary key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists room_movie_seating;
-- +goose StatementEnd
