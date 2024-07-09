create table if not exists cities (
    id serial,
    name varchar(255),
    country varchar(255),
    latitude double precision,
    longitude double precision,
    primary key (id),
    unique (name, country)
);

create table if not exists forecasts (
    id serial,
    city_id int,
    temp real,
    date date,
    forecast_json jsonb,
    primary key (id),
    unique (city_id, date),
    foreign key (city_id) references cities(id) on delete cascade
);