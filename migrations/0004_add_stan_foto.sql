-- Migration: Add foto column to stans table
-- Date: 2026-01-31

-- Add foto column to stans table
ALTER TABLE stans
ADD COLUMN foto VARCHAR(255);

-- Add comment for documentation
COMMENT ON COLUMN stans.foto IS 'Path to stan photo/image file, supports base64 upload and URL';