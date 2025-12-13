-- Migration: Add Path column to Images table
-- This allows storing the file system path where processed images are saved

BEGIN TRANSACTION;

-- Add Path column to Images table
ALTER TABLE Images ADD COLUMN Path TEXT NOT NULL DEFAULT '';

-- Update existing records to have empty path (if any exist)
-- In production, you'd want to backfill these with actual paths
UPDATE Images SET Path = '' WHERE Path IS NULL OR Path = '';

COMMIT;
