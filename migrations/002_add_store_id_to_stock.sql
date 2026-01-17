-- Migration: Add store_id to stock table
-- This migration adds store_id to the stock table to support stock per store
-- This migration is idempotent and can be executed multiple times safely

-- Step 1: Add store_id column (nullable initially to handle existing data)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'stock' AND column_name = 'store_id'
    ) THEN
        ALTER TABLE stock
        ADD COLUMN store_id UUID REFERENCES stores(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Step 2: Populate store_id from products for existing stock records (only where NULL)
UPDATE stock s
SET store_id = p.store_id
FROM products p
WHERE s.product_id = p.id AND s.store_id IS NULL;

-- Step 3: Make store_id NOT NULL now that all records have values
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'stock'
        AND column_name = 'store_id'
        AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE stock
        ALTER COLUMN store_id SET NOT NULL;
    END IF;
END $$;

-- Step 4: Drop the old unique constraint
ALTER TABLE stock
DROP CONSTRAINT IF EXISTS stock_tenant_id_product_id_key;

-- Step 5: Add new unique constraint including store_id (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'stock_tenant_id_product_id_store_id_key'
    ) THEN
        ALTER TABLE stock
        ADD CONSTRAINT stock_tenant_id_product_id_store_id_key UNIQUE(tenant_id, product_id, store_id);
    END IF;
END $$;

-- Step 6: Drop old index that is no longer optimal
DROP INDEX IF EXISTS idx_stock_tenant_product;

-- Step 7: Add new indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_stock_tenant_product_store ON stock(tenant_id, product_id, store_id);
CREATE INDEX IF NOT EXISTS idx_stock_store ON stock(store_id);
