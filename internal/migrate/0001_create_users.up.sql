create table users (
id serial primary key,
username text not null,
firstname text not null,
lastname text not null,
email text not null,
password text not null,
phone text not null,
create_at timestamp not null,
update_at timestamp default null,
delete_at timestamp default null
);

create table black_lists (
id serial primary key,
jti text not null
)