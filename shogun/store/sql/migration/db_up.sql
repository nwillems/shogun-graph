
CREATE TABLE nodes (
  id          INTEGER PRIMARY KEY,
  node_name   VARCHAR NOT NULL UNIQUE,
  properties  JSONB
);


CREATE TABLE edges (
  id      INTEGER PRIMARY KEY,
  source  INTEGER NOT NULL,
  target  INTEGER NOT NULL,
  type    TEXT NOT NULL DEFAULT '',

  properties JSONB,

  FOREIGN KEY(source) REFERENCES nodes(id),
  FOREIGN KEY(target) REFERENCES nodes(id)
);
