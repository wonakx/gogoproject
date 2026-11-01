-- Reference schema for a NEW installation, including the approved additions.
-- Do not run against existing tables; use migrations/ instead.
CREATE DATABASE IF NOT EXISTS dealer;

CREATE TABLE dealer.dt_member (
    member_no BIGINT AUTO_INCREMENT PRIMARY KEY,
    id VARCHAR(100) NOT NULL,
    email VARCHAR(200) NULL,
    address VARCHAR(300) NULL,
    phone_number VARCHAR(100) NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NULL,
    CONSTRAINT dt_member_id_unique_key UNIQUE (id)
);

CREATE TABLE dealer.dt_product (
    product_no BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NULL,
    description LONGTEXT NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NULL,
    CONSTRAINT dt_product_product_code_key UNIQUE (product_code)
);

CREATE TABLE dealer.dt_sale (
    sale_no BIGINT AUTO_INCREMENT PRIMARY KEY,
    sell_member_no BIGINT NOT NULL,
    name VARCHAR(200) NULL,
    description LONGTEXT NULL,
    product_no BIGINT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(30) NOT NULL,
    start_at DATETIME NOT NULL,
    end_at DATETIME NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NULL
);

CREATE TABLE dealer.dt_purchase (
    purchase_no BIGINT AUTO_INCREMENT PRIMARY KEY,
    sale_no BIGINT NOT NULL,
    sell_member_no BIGINT NOT NULL,
    buy_member_no BIGINT NOT NULL,
    quantity INT NOT NULL,
    purchase_status VARCHAR(30) NOT NULL,
    purchase_at DATETIME NOT NULL,
    create_at DATETIME NOT NULL
);

CREATE TABLE dealer.dt_purchase_history (
    purchase_history_no BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_no BIGINT NULL,
    sale_no BIGINT NOT NULL,
    sell_member_no BIGINT NOT NULL,
    buy_member_no BIGINT NOT NULL,
    quantity INT NOT NULL,
    purchase_status VARCHAR(30) NOT NULL,
    purchase_at DATETIME NOT NULL,
    create_at DATETIME NOT NULL
);
