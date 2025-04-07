drop database IF EXISTS ord;
create database ord;
use ord;
drop table IF EXISTS stock;
drop table IF EXISTS order1;
drop table IF EXISTS user_coupon;
drop table IF EXISTS pay;

create table stock(
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT,
  product_id int(11) not null,
  stock int(11) not null,
  create_time datetime DEFAULT current_timestamp,
  update_time datetime DEFAULT current_timestamp,
  unique key(product_id)
);

insert into stock(product_id, stock) VALUES (1, 100), (2, 200);

create table order1(
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT,
  order_id varchar(45),
  user_id int(11) not null,
  product_id int(11) not null,
  amount int(11) not null,
  status varchar(45) not null,
  create_time datetime DEFAULT current_timestamp,
  update_time datetime DEFAULT current_timestamp,
  UNIQUE key(order_id),
  key(product_id),
  key(user_id)
);

create table user_coupon(
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT,
  user_id int(11) not null,
  coupon_id int(11) not null,
  used int(11) not null,
  create_time datetime DEFAULT current_timestamp,
  update_time datetime DEFAULT current_timestamp,
  UNIQUE key(user_id, coupon_id),
  key(coupon_id)
);

create table pay(
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT,
  user_id int(11) not null,
  order_id varchar(45) not null,
  amount int(11) not null,
  status varchar(45) not null,
  create_time datetime DEFAULT current_timestamp,
  update_time datetime DEFAULT current_timestamp,
  UNIQUE key(order_id),
  key(user_id)
);