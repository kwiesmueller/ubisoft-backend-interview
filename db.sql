CREATE TABLE IF NOT EXISTS entries (
    id            serial,
    session_id    VARCHAR(50) NOT null,
    user_id       VARCHAR(50) NOT null,
    rating        INT8 NOT null,
    comment       VARCHAR(50),
    PRIMARY key (session_id, user_id)
);
CREATE UNIQUE INDEX ON entries (id);
