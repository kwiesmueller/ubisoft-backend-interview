CREATE TABLE IF NOT EXISTS entries (
    id            serial,
    session_id    VARCHAR(50) NOT null,
    user_id       VARCHAR(50) NOT null,
    rating        INT8 NOT null,
    comment       VARCHAR(50),
    PRIMARY key (session_id, user_id)
);
CREATE INDEX IF NOT EXISTS entries_id ON entries (id);
