CREATE TABLE main.transaction (
	id integer NOT NULL,
	amount double precision DEFAULT 0.0 NOT NULL,
	state character varying(10) NOT NULL,
	src_type character varying NOT NULL,
	dt timestamp without time zone
);

ALTER TABLE main.transaction OWNER TO postgres;

--------------------------------------------------------------------------------

ALTER TABLE main.transaction
	ADD CONSTRAINT transaction_un UNIQUE (id);
