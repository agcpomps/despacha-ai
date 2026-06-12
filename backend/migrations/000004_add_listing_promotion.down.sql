DROP INDEX IF EXISTS idx_listings_featured;

ALTER TABLE listings
    DROP COLUMN IF EXISTS is_featured,
    DROP COLUMN IF EXISTS featured_until,
    DROP COLUMN IF EXISTS bumped_at;
