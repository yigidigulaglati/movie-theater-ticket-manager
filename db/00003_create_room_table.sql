-- +goose Up
-- +goose StatementBegin
create table if not exists room(
    id integer not null,
    room_name text not null,
    seat_price integer not null,

    check (seat_price >= 0),
    check (length(room_name) >= 2),

    primary key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists room;
-- +goose StatementEnd
