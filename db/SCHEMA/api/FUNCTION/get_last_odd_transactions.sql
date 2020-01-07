CREATE OR REPLACE FUNCTION api.get_last_odd_transactions() RETURNS TABLE(id integer, amount double precision, state character varying, src_type character varying)
    LANGUAGE sql
    AS $$
 select id, amount, state, src_type
 from only main."transaction"
 where mod(id, 2) = 0
 order by id desc
 limit 10;
$$;

ALTER FUNCTION api.get_last_odd_transactions() OWNER TO postgres;
