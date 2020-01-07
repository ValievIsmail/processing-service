CREATE OR REPLACE FUNCTION api.post_processing(_transaction_id integer, _amount double precision, state text, _src_type character varying) RETURNS void
    LANGUAGE sql
    AS $$
with cancel_transaction as (
 INSERT INTO main."canceled"
 (id, amount, state, src_type, dt)
 VALUES(_transaction_id, _amount, state, _src_type, now())
), delete_old as (
 DELETE FROM ONLY main."transaction"
 WHERE id = _transaction_id
)
 select case when state = 'win' 
 then (select api.update_balance_lost(_amount, _src_type))
 else (select api.update_balance_win(_amount, _src_type))
 end;
$$;

ALTER FUNCTION api.post_processing(_transaction_id integer, _amount double precision, state text, _src_type character varying) OWNER TO postgres;
