package db_connection

import (
	"database/sql"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS delivery(
	delivery_id     INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_receiver VARCHAR(70),
    phone_number  VARCHAR(15),
    zip_code      VARCHAR(10),
    name_city     VARCHAR(35),
    address     VARCHAR(100),
    region        VARCHAR(20),
    email         VARCHAR(255)
);
CREATE TABLE  IF NOT EXISTS payment(
    payment_id     INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    transaction_id VARCHAR(32) NOT NULL,
    request_id     VARCHAR(32),
    name_currency  VARCHAR(10),
    name_provider  VARCHAR(15) NOT NULL,
    total_amount   INT,
    payment_dt     BIGINT,
    name_bank      VARCHAR(30),
    delivery_cost  INT,
    goods_total    INT,
    custom_fee     INT
);
CREATE TABLE IF NOT EXISTS buy(
    buy_id             INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_uid          VARCHAR(32),
    track_number       VARCHAR(30),
    entry         	   VARCHAR(15),
    locale             VARCHAR(10),
    internal_signature VARCHAR(10),
    customer_id      VARCHAR(32) NOT NULL,
    delivery_service   VARCHAR(15),
    shardkey           VARCHAR(10),
    sm_id              INT NOT NULL,
    date_created       TIMESTAMPTZ NOT NULL,
    oof_shard          VARCHAR(10)
);
CREATE TABLE IF NOT EXISTS item(
    item_id      INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id		 INT,
    chrt_id     BIGINT NOT NULL,
    track_number VARCHAR(30),
    price        INT,
    rid          VARCHAR(32),
    name_item    VARCHAR(50),
    sale         INT,
    size_item    VARCHAR(15),
    total_price  INT,
    nm_id        BIGINT,
    brand        VARCHAR(50),
    status   	 INT NOT NULL
);
	`)
	if err != nil {
		return err
	}
	return nil
}
