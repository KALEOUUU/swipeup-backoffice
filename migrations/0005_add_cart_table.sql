-- Migration: Add cart table for mobile basket functionality
-- Date: 2026-01-31

-- Create cart table for temporary order storage
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    id_siswa INTEGER NOT NULL REFERENCES siswas(id) ON DELETE CASCADE,
    id_menu INTEGER NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    qty INTEGER NOT NULL CHECK (qty > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(id_siswa, id_menu) -- One cart item per siswa-menu combination
);

-- Add indexes for performance
CREATE INDEX idx_carts_id_siswa ON carts(id_siswa);
CREATE INDEX idx_carts_id_menu ON carts(id_menu);

-- Add comments
COMMENT ON TABLE carts IS 'Temporary cart storage for mobile app basket functionality';
COMMENT ON COLUMN carts.id_siswa IS 'Reference to siswa who owns the cart item';
COMMENT ON COLUMN carts.id_menu IS 'Reference to menu item in cart';
COMMENT ON COLUMN carts.qty IS 'Quantity of menu item in cart';