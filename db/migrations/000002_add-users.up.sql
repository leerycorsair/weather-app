create table if not exists users (
    id serial,
    login varchar(255) unique,
    password varchar(255),
    email varchar(255),
    primary key (id)
);

create table if not exists users_cities (
    id serial,
    user_id int,
    city_id int,
    primary key (id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (city_id) references cities(id) on delete cascade,
    unique (user_id, city_id)
);