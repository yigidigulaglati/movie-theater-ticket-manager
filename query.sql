
-- name: SelectUsernamePassword :one
select id, username, password from user
where id = ? limit 1;

-- name: InsertUser :execrows
insert into user(username, password) 
values (?, ?);

-- name: SelectUserWithUsername :one
select id, username, password from user
where username = ?;



-- name: InsertNewMovie :execrows
insert into movie(title, director, release_year, genre, rating, description, duration)
values (?,?,?,?,?,?,?);

-- name: DeleteMovieWithID :execrows
delete from movie where id = ?;

-- name: UpdateMovieWithID :execrows
update movie
set title = ?, director = ?, release_year = ?, genre = ?, rating = ?, description = ?, duration = ?
where id = ?;

-- name: SelectMovieWithID :one
select id from movie
where id = ? limit 1;

-- name: SelectMovieAllColumnsWithID :one
select id, title, director, release_year, genre, rating, description, duration 
from movie 
where id = ?
limit 1;


-- name: InsertNewRoom :execrows
insert into room(room_name, seat_price)
values (?, ?);

-- name: DeleteRoomWithID :execrows
delete from room where id = ?;

-- name: UpdateRoomWithID :execrows
update room
set room_name = ?, seat_price = ?
where id = ?;


-- name: SelectAllSeatingOfRoomWithID :many
select id, room_id, row_index, col_index, seat from seating
where id = ?;

-- name: DeleteStaleMovieSchedules :execrows
delete from room_movie where start_date_unix < ?;

-- name: SelectRoomsOfMovieJoinRoomJoinSeating :many
select room_movie.id, room_movie.room_id, room_movie.movie_id, room_movie.start_date_unix, room.room_name, room.seat_price, room_movie_seating.seat from room_movie inner join room on room_movie.room_id = room.id inner join room_movie_seating on room_movie.id = room_movie_seating.room_movie_id where room_movie.movie_id = ? and room_movie.start_date_unix > ? order by room_movie.room_id desc, room_movie.start_date_unix desc;

-- name: SelectRoomOfMovieJoinRoomJoinSeating :many
select room_movie.id, room_movie_seating.id as room_movie_seating_id, room_movie.room_id, room_movie.movie_id, room_movie.start_date_unix, room.room_name, room.seat_price, room_movie_seating.seat, room_movie_seating.row_index, room_movie_seating.col_index from room_movie inner join room on room_movie.room_id = room.id inner join room_movie_seating on room_movie.id = room_movie_seating.room_movie_id where room_movie.id = ?;

-- name: UpdateSeatInRoomMovieSeating :execrows
update room_movie_seating set seat = ? where room_movie_id = ? and row_index = ? and col_index = ?;

-- name: SelectSeatInRoomMovieSeating :one
select room_movie.id as room_movie_id, room_movie.room_id as room_id, room_movie.movie_id as movie_id, room_movie.start_date_unix as start_date_unix, room_movie_seating.id as room_movie_seating_id, room_movie_seating.row_index as row_index, room_movie_seating.col_index as col_index, room_movie_seating.seat as seat from room_movie_seating join room_movie on room_movie_seating.room_movie_id = room_movie.id where room_movie_seating.room_movie_id = ? and room_movie_seating.row_index = ? and room_movie_seating.col_index = ?;