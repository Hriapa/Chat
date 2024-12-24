CREATE TABLE IF NOT EXISTS users
(
    Id serial PRIMARY KEY,
    login varchar(30) NOT NULL,
    encrypted_password varchar(100) NOT NULL,
    name varchar(30),
    familyname varchar(30),
    surname varchar(30),
    birthdate date
);

CREATE TABLE IF NOT EXISTS rooms
(
    id serial PRIMARY KEY,
    user_id1 integer NOT NULL,
    user_id2 integer NOT NULL,
    quntity integer NOT NULL
);

CREATE TABLE IF NOT EXISTS nr_messages
(
    user_id integer NOT NULL,
    room_id integer NOT NULL,
    number integer NOT NULL
);

CREATE TABLE IF NOT EXISTS messages
(
    index_number integer NOT NULL,
    user_id integer NOT NULL,
    format smallint NOT NULL,
    message bytea,
    date timestamp,
    reference integer REFERENCES rooms (id)
)