-- Previous per-pantry choices cannot be reconstructed after a bulk reset.
ALTER TABLE products ALTER COLUMN default_status SET DEFAULT 'DISCARDED';
UPDATE products SET default_status = 'DISCARDED' WHERE default_status = 'PENDING';
