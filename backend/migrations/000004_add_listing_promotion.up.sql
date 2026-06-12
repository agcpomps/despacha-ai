ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS featured_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS bumped_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_listings_featured
    ON listings (is_featured, featured_until)
    WHERE is_featured = TRUE;
