-- Migration: Add support for 2-level discounts (Global & Stan)
-- Date: 2025-01-31

-- Add new columns to diskons table
ALTER TABLE diskons 
ADD COLUMN tipe_diskon VARCHAR(20) NOT NULL DEFAULT 'global',
ADD COLUMN id_stan INTEGER REFERENCES stans(id) ON DELETE CASCADE;

-- Create index for better query performance
CREATE INDEX idx_diskons_id_stan ON diskons(id_stan);
CREATE INDEX idx_diskons_tipe ON diskons(tipe_diskon);

-- Update existing data to be global discount
UPDATE diskons SET tipe_diskon = 'global' WHERE tipe_diskon IS NULL OR tipe_diskon = '';

COMMENT ON COLUMN diskons.tipe_diskon IS 'Type of discount: global (superadmin) or stan (per-stan)';
COMMENT ON COLUMN diskons.id_stan IS 'NULL for global discount, Stan ID for stan-specific discount';
