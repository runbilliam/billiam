CREATE TABLE users (
   id         CHAR(26) PRIMARY KEY,
   version    INTEGER NOT NULL DEFAULT 1,
   email      TEXT UNIQUE NOT NULL,
   password   TEXT NOT NULL,
   timezone   TEXT NOT NULL,
   active     BOOLEAN NOT NULL DEFAULT false,
   created_at TIMESTAMPTZ NOT NULL,
   updated_at TIMESTAMPTZ,
   login_at   TIMESTAMPTZ
);

---- create above / drop below ----

DROP TABLE IF EXISTS users CASCADE;
