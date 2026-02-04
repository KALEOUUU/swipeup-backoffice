-- Migration: Add superadmin role, inventory (stock), and payment channel features
-- Date: 2026-02-03

-- Add new fields to stan table for payment channels
ALTER TABLE stan ADD COLUMN IF NOT EXISTS qris_image VARCHAR(255);
ALTER TABLE stan ADD COLUMN IF NOT EXISTS accept_cash BOOLEAN DEFAULT TRUE;
ALTER TABLE stan ADD COLUMN IF NOT EXISTS accept_qris BOOLEAN DEFAULT FALSE;

-- Add stock and availability fields to menu table for inventory management
ALTER TABLE menu ADD COLUMN IF NOT EXISTS stock INTEGER DEFAULT 0;
ALTER TABLE menu ADD COLUMN IF NOT EXISTS is_available BOOLEAN DEFAULT TRUE;

-- Update existing menu items to be available if no stock tracking was done before
UPDATE menu SET is_available = TRUE WHERE is_available IS NULL;
UPDATE menu SET stock = 0 WHERE stock IS NULL;

-- Note: The 'superadmin' role is now valid in the user.role field
-- Existing roles: admin_stan, siswa
-- New role: superadmin

-- Insert default superadmin user (password: superadmin123 - should be changed immediately)
-- Password hash for 'superadmin123' using bcrypt
-- You can generate this hash using: 
-- go run -e 'import "golang.org/x/crypto/bcrypt"; h,_:=bcrypt.GenerateFromPassword([]byte("superadmin123"),bcrypt.DefaultCost); println(string(h))'

-- IMPORTANT: Run this INSERT manually or via seed script
-- INSERT INTO users (username, password, role, created_at, updated_at) 
-- VALUES ('superadmin', '$2a$10$YOUR_BCRYPT_HASH_HERE', 'superadmin', NOW(), NOW());

COMMENT ON COLUMN stan.qris_image IS 'Path to QRIS image for payment';
COMMENT ON COLUMN stan.accept_cash IS 'Whether this stan accepts cash payment';
COMMENT ON COLUMN stan.accept_qris IS 'Whether this stan accepts QRIS payment';
COMMENT ON COLUMN menu.stock IS 'Current stock quantity of menu item';
COMMENT ON COLUMN menu.is_available IS 'Whether menu item is available for order';
