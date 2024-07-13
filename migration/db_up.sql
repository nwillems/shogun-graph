
CREATE TABLE nodes (
  id INTEGER PRIMARY KEY,
  node_name varchar NOT NULL,
  properties JSONB
)


CREATE TABLE edges (
  id INTEGER PRIMARY KEY,
  fst INTEGER NOT NULL,
  snd INTEGER NOT NULL,
  type TEXT DEFAULT '',

  properties JSONB,

  FOREIGN KEY(left) REFERENCES nodes(id),
  FOREIGN KEY(right) REFERENCES nodes(id)
)
