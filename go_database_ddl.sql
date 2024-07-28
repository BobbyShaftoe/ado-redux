create table if not exists locations
(
    id   integer not null
        constraint locations_pk
            primary key,
    name varchar not null
);

alter table locations
    owner to postgres;