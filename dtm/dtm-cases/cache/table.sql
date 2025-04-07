drop database IF EXISTS cache1;
create database cache1;
use cache1;
drop table IF EXISTS ver;

create table ver (
  k VARCHAR(100) PRIMARY KEY,
  v VARCHAR(100) default '',
  time_cost varchar(45) default 0
);