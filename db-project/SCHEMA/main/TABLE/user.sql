CREATE TABLE main."user" (
	src_type character varying NOT NULL,
	balance double precision DEFAULT 0 NOT NULL
);

ALTER TABLE main."user" OWNER TO postgres;
