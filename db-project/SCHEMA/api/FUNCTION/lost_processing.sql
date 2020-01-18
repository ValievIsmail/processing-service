CREATE OR REPLACE FUNCTION api.lost_processing(_transaction_id integer, _amount double precision, _src_type character varying) RETURNS void
    LANGUAGE sql
    AS $$
 with insert_transaction as (
 	INSERT INTO main."transaction"
	(id, amount, state, src_type, dt)
	VALUES(_transaction_id, _amount, 'lost', _src_type, now())
)
 select api.update_balance_lost(_amount, _src_type);
$$;

ALTER FUNCTION api.lost_processing(_transaction_id integer, _amount double precision, _src_type character varying) OWNER TO postgres;
