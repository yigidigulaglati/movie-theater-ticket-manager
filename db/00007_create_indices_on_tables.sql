-- +goose Up
-- +goose StatementBegin
create index username_index on user(username);

create index movie_title_index on movie(title);
create index movie_genre_index on movie(genre);
create index movie_rating_index on movie(rating);

create index seating_room_id_index on seating(room_id);

create index room_movie_seating_room_movie_id_row_col_index on room_movie_seating(room_movie_id, row_index, col_index);

create index room_movie_movie_id_index on room_movie(movie_id);
create index room_movie_start_date_unix_index on room_movie(start_date_unix);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists movie_rating_index;
drop index if exists movie_genre_index;
drop index if exists movie_title_index;
drop index if exists seating_room_id_index;
drop index if exists room_movie_seating_room_movie_id_row_col_index;
drop index if exists room_movie_movie_id_index;
drop index if exists room_movie_start_date_unix_index,
-- +goose StatementEnd
