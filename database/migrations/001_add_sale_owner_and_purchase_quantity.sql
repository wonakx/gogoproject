-- Existing dealer tables: expand with nullable columns so old rows are preserved.
-- Run only once. This project does not execute migrations automatically.
ALTER TABLE dealer.dt_sale
    ADD COLUMN sell_member_no BIGINT NULL AFTER sale_no;

ALTER TABLE dealer.dt_purchase
    ADD COLUMN quantity INT NULL AFTER buy_member_no;

ALTER TABLE dealer.dt_purchase_history
    ADD COLUMN quantity INT NULL AFTER buy_member_no;

-- Before migration 002, backfill every existing row with actual values:
-- dt_sale.sell_member_no: the actual member_no of the seller.
-- dt_purchase.quantity: the actual purchased quantity.
-- dt_purchase_history.quantity: the quantity at the time of that history entry.
-- These values cannot be inferred safely from the original schema.
