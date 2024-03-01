CREATE DATABASE IF NOT EXISTS orders;

CREATE TABLE IF NOT EXISTS orders
(
    id          varchar(36) not null
        primary key,
    price       double      not null,
    tax         double      null,
    final_price double      null
);