CREATE SCHEMA main;

CREATE SCHEMA api;

CREATE TABLE main.transaction (
	id int4 NOT NULL,
	amount float8 NOT NULL DEFAULT 0.0,
	state varchar(10) NOT NULL,
	src_type varchar not null,
	dt timestamp NULL,
	CONSTRAINT transaction_un UNIQUE (id)
);

CREATE TABLE main.canceled (
)
INHERITS (main."transaction");

CREATE TABLE main."user" (
	src_type varchar NOT NULL,
	balance float8 NOT NULL DEFAULT 0
);

CREATE OR REPLACE FUNCTION api.get_last_odd_transactions()
 RETURNS TABLE(id integer, amount double precision, state character varying, src_type varchar)
 LANGUAGE sql
AS $$
 select id, amount, state, src_type
 from only main."transaction"
 where mod(id, 2) = 0
 order by id desc
 limit 10;
$$;

CREATE OR REPLACE FUNCTION api.update_balance_lost(_amount double precision, _src_type varchar)
 RETURNS void
 LANGUAGE sql
AS $$
 UPDATE main."user"
 set balance = case
 	when (balance - _amount) < 0 then 0
 	else balance - _amount
 end
 where src_type = _src_type;
$$;

CREATE OR REPLACE FUNCTION api.update_balance_win(_amount double precision, _src_type varchar)
 RETURNS void
 LANGUAGE sql
AS $$
 UPDATE main."user"
 SET balance = balance + _amount
 where src_type = _src_type;
$$;

CREATE OR REPLACE FUNCTION api.post_processing(_transaction_id integer, _amount double precision, state text, _src_type varchar)
 RETURNS void
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

CREATE OR REPLACE FUNCTION api.win_processing(_transaction_id integer, _amount float8, _src_type varchar)
 RETURNS void
 LANGUAGE sql
AS $$
 with insert_transaction as (
 	INSERT INTO main."transaction"
	(id, amount, state, src_type, dt)
	VALUES(_transaction_id, _amount, 'win', _src_type, now())
)
 select api.update_balance_win(_amount, _src_type);
$$;

CREATE OR REPLACE FUNCTION api.lost_processing(_transaction_id integer, _amount double precision, _src_type varchar)
 RETURNS void
 LANGUAGE sql
AS $$
 with insert_transaction as (
 	INSERT INTO main."transaction"
	(id, amount, state, src_type, dt)
	VALUES(_transaction_id, _amount, 'lost', _src_type, now())
)
 select api.update_balance_lost(_amount, _src_type);
$$;

INSERT INTO main."user" ("src_type", "balance")
select 'client', 0;

INSERT INTO main."user" ("src_type", "balance")
select 'game', 0;

INSERT INTO main."user" ("src_type", "balance")
select 'server', 0;

INSERT INTO main."user" ("src_type", "balance")
select 'payment', 0;

WITH generate_client_data as (
	select 
 	generate_series(1, 10) id, 
 	round((random() * 100 + 1)::numeric, 2)::float8 amount, 
 	'win' state, 
 	'client' src_type,
 	generate_series(timestamp '2020-01-10 10:00:00', 
 					timestamp '2020-01-10 19:00:00', interval '1 hour') dt
 ), sum_amount as (
 	select sum(gd.amount) as sum from generate_client_data as gd
 ), update_balance as (
  	update main."user"
 	set balance = balance + sa.sum
 	from sum_amount as sa
 	where src_type = 'client'
 )
insert into main."transaction" (id, amount, state, src_type, dt)
select id, amount, state, src_type, dt from generate_client_data;

WITH generate_game_data as (
	select 
 	generate_series(11, 20) id, 
 	round((random() * 100 + 1)::numeric, 2)::float8 amount, 
 	'win' state, 
 	'game' src_type,
 	generate_series(timestamp '2020-02-10 10:00:00', 
 					timestamp '2020-02-10 19:00:00', interval '1 hour') dt
 ), sum_amount as (
 	select sum(gd.amount) as sum from generate_game_data as gd
 ), update_balance as (
  	update main."user"
 	set balance = balance + sa.sum
 	from sum_amount as sa
 	where src_type = 'game'
 )
insert into main."transaction" (id, amount, state, src_type, dt)
select id, amount, state, src_type, dt from generate_game_data;

WITH generate_server_data as (
	select 
 	generate_series(21, 30) id, 
 	round((random() * 100 + 1)::numeric, 2)::float8 amount, 
 	'win' state, 
 	'server' src_type,
 	generate_series(timestamp '2020-03-10 10:00:00', 
 					timestamp '2020-03-10 19:00:00', interval '1 hour') dt
 ), sum_amount as (
 	select sum(gd.amount) as sum from generate_server_data as gd
 ), update_balance as (
  	update main."user"
 	set balance = balance + sa.sum
 	from sum_amount as sa
 	where src_type = 'server'
 )
insert into main."transaction" (id, amount, state, src_type, dt)
select id, amount, state, src_type, dt from generate_server_data;

WITH generate_payment_data as (
	select 
 	generate_series(31, 40) id, 
 	round((random() * 100 + 1)::numeric, 2)::float8 amount, 
 	'win' state, 
 	'payment' src_type,
 	generate_series(timestamp '2020-04-10 10:00:00', 
 					timestamp '2020-04-10 19:00:00', interval '1 hour') dt
 ), sum_amount as (
 	select sum(gd.amount) as sum from generate_payment_data as gd
 ), update_balance as (
  	update main."user"
 	set balance = balance + sa.sum
 	from sum_amount as sa
 	where src_type = 'payment'
 )
insert into main."transaction" (id, amount, state, src_type, dt)
select id, amount, state, src_type, dt from generate_payment_data;