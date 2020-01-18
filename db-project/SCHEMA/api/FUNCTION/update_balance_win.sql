CREATE OR REPLACE FUNCTION api.update_balance_win(_amount double precision, _src_type character varying) RETURNS void
    LANGUAGE sql
    AS $$
 UPDATE main."user"
 SET balance = balance + _amount
 where src_type = _src_type;
$$;

ALTER FUNCTION api.update_balance_win(_amount double precision, _src_type character varying) OWNER TO postgres;
