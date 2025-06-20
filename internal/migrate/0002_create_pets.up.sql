create table categories (
id serial primary key,
name text not null
);

create table tags (
id serial primary key,
name text not null
);

create table pets (
id integer primary key,
name text not null,
photo_url text,
category_id integer not null references categories(id),
status text not null
);

create table pet_tags (
pet_id integer not null references pets(id),
tag_id integer not null references tags(id),
primary key (pet_id, tag_id)
);

insert into categories ("id", "name") values
(1, 'dog'),
(2, 'cat');

insert into tags ("id", "name") values 
(1, 'friendly'),
(2, 'wild'),
(3, 'trained')