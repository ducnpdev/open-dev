create schema dtm_busi;


CREATE SEQUENCE IF NOT EXISTS dtm_busi.user_account_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;
CREATE TABLE dtm_busi.user_account (
 id bigint NOT NULL DEFAULT nextval('dtm_busi.user_account_id_seq'::regclass) PRIMARY KEY,
  user_id bigint not NULL UNIQUE ,
  balance decimal(10,2) NOT NULL DEFAULT '0.00',
  trading_balance decimal(10,2) NOT NULL DEFAULT '0.00',
  create_time timestamp without time zone DEFAULT now(),
  update_time timestamp without time zone DEFAULT now()
);