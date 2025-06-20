create table orders (
id integer primary key not null,
pet_id integer not null,
quantity integer not null,
ship_date text not null,
status text not null,
complete boolean not null
)