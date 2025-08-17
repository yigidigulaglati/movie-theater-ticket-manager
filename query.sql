
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


-- name: SelectJoinRoomMovieAndSeatingWithRoomID :many
select room_movie.id, seating.row_index, seating.col_index, seating.seat from 
room_movie inner join seating 
on room_movie.room_id = seating.room_id
where room_movie.room_id = ?;
