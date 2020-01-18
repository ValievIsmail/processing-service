CREATE TABLE main.canceled (
)
INHERITS (main.transaction);

ALTER TABLE ONLY main.canceled ALTER COLUMN amount SET DEFAULT 0.0;

ALTER TABLE main.canceled OWNER TO postgres;
