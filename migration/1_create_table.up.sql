BEGIN;

DROP TABLE IF EXISTS owner CASCADE;
CREATE TABLE owner (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,

    CONSTRAINT unique_owner UNIQUE (name, surname, patronymic)
);

DROP TABLE IF EXISTS car CASCADE;
CREATE TABLE car (
    id SERIAL PRIMARY KEY,
    regNum TEXT NOT NULL,
    mark TEXT NOT NULL,
    model TEXT NOT NULL,
    year INTEGER,
    owner_id INTEGER NOT NULL REFERENCES owner(id) ON DELETE CASCADE,

    CONSTRAINT unique_car UNIQUE (regNum, mark, model, year)
);

COMMIT;
