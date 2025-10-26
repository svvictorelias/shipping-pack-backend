-- Schema initialization for pack calculation service using SERIAL integers

-- Table for available pack sizes
CREATE TABLE IF NOT EXISTS packs (
    id SERIAL PRIMARY KEY,
    size INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table for storing calculation history
CREATE TABLE IF NOT EXISTS calculations (
    id SERIAL PRIMARY KEY,
    items INTEGER NOT NULL,
    total_items INTEGER NOT NULL,
    pack_count INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table for pack usage details linked to a specific calculation
CREATE TABLE IF NOT EXISTS calculation_items (
    id SERIAL PRIMARY KEY,
    calculation_id INTEGER NOT NULL REFERENCES calculations(id) ON DELETE CASCADE,
    pack_size INTEGER NOT NULL,
    quantity INTEGER NOT NULL
);

-- Helpful indexes for query performance
CREATE INDEX IF NOT EXISTS idx_packs_size ON packs(size);
CREATE INDEX IF NOT EXISTS idx_calculation_items_calc_id ON calculation_items(calculation_id);
