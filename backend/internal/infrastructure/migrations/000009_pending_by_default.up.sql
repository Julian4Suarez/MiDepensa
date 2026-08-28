-- Every active product starts undecided. Archived products remain archived.
UPDATE products SET default_status = 'PENDING' WHERE default_status <> 'ARCHIVED';
UPDATE pantry_items SET status = 'PENDING' WHERE status <> 'ARCHIVED';
ALTER TABLE products ALTER COLUMN default_status SET DEFAULT 'PENDING';
