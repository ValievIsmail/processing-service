CREATE OR REPLACE FUNCTION api.update_balance_lost(_amount double precision, _src_type character varying) RETURNS void
    LANGUAGE sql
    AS $$
 UPDATE main."user"
 set balance = case
 	when (balance - _amount) < 0 then 0
 	else balance - _amount
 end
 where src_type = _src_type;
$$;

ALTER FUNCTION api.update_balance_lost(_amount double precision, _src_type character varying) OWNER TO postgres;
