-- Run only after 001 and after backfilling existing rows.
-- Check that these queries return zero before executing the ALTER statements.
SELECT COUNT(*) AS missing_sale_owner FROM dealer.dt_sale WHERE sell_member_no IS NULL;
SELECT COUNT(*) AS missing_purchase_quantity FROM dealer.dt_purchase WHERE quantity IS NULL;
SELECT COUNT(*) AS missing_history_quantity FROM dealer.dt_purchase_history WHERE quantity IS NULL;

-- Apply with strict SQL mode; do not allow implicit conversion of NULL to zero.
-- For populated databases, execute the checks separately and stop if any is nonzero.
ALTER TABLE dealer.dt_sale MODIFY COLUMN sell_member_no BIGINT NOT NULL;
ALTER TABLE dealer.dt_purchase MODIFY COLUMN quantity INT NOT NULL;
ALTER TABLE dealer.dt_purchase_history MODIFY COLUMN quantity INT NOT NULL;
