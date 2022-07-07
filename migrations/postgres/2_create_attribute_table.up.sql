CREATE TYPE attribute_types AS ENUM ('datetime', 'text', 'number');


CREATE TABLE IF NOT EXISTS attribute (
    id uuid primary key,
    name varchar(255) not null,
    type attribute_types not null
);