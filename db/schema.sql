CREATE TABLE IF NOT EXISTS properties (
  id serial,
  name citext NOT NULL UNIQUE,
  PRIMARY KEY(id)
);