CREATE TABLE IF NOT EXISTS pantries (
    id         UUID PRIMARY KEY,
    slug       TEXT        NOT NULL UNIQUE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Read-only catalog shared by every pantry; populated by the seed migration.
CREATE TABLE IF NOT EXISTS products (
    id               UUID PRIMARY KEY,
    code             TEXT    NOT NULL UNIQUE,
    name             TEXT    NOT NULL,
    image            TEXT    NOT NULL,
    default_category TEXT    NOT NULL CHECK (default_category IN ('FRESH', 'PANTRY', 'DRINKS', 'HOME_CARE')),
    default_view     TEXT    NOT NULL CHECK (default_view IN ('PRIMARY', 'SECONDARY', 'OTHER')),
    sort_order       INTEGER NOT NULL
);

-- Per-pantry state of a catalog product. "view" is spelled out as
-- pantry_view to avoid shadowing the SQL keyword.
CREATE TABLE IF NOT EXISTS pantry_items (
    pantry_id   UUID        NOT NULL REFERENCES pantries (id) ON DELETE CASCADE,
    product_id  UUID        NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    status      TEXT        NOT NULL CHECK (status IN ('OUT', 'LOW', 'OK')),
    pantry_view TEXT        NOT NULL CHECK (pantry_view IN ('PRIMARY', 'SECONDARY', 'OTHER')),
    category    TEXT        NOT NULL CHECK (category IN ('FRESH', 'PANTRY', 'DRINKS', 'HOME_CARE')),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (pantry_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_pantry_items_pantry_id ON pantry_items (pantry_id);
