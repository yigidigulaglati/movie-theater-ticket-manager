-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

create table if not exists seating(
    id integer not null,
    room_id integer not null,
    row_index integer not null,
    col_index integer not null,
    seat integer not null,

    foreign key(room_id) references room(id) on delete cascade,
    primary key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists seating;
-- +goose StatementEnd
