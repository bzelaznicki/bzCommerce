-- +goose Up
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    iso_code CHAR(2) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO countries (name, iso_code, sort_order) VALUES
('United States', 'US', 1),
('Canada', 'CA', 2),
('United Kingdom', 'GB', 3),
('Germany', 'DE', 4),
('France', 'FR', 5),
('Australia', 'AU', 6),
('New Zealand', 'NZ', 7),
('India', 'IN', 8),
('Japan', 'JP', 9),
('Brazil', 'BR', 10),
('Mexico', 'MX', 11),
('Netherlands', 'NL', 12),
('Spain', 'ES', 13),
('Italy', 'IT', 14),
('Sweden', 'SE', 15),
('Norway', 'NO', 16),
('Finland', 'FI', 17),
('Switzerland', 'CH', 18),
('Austria', 'AT', 19),
('Ireland', 'IE', 20),
('Denmark', 'DK', 21),
('Belgium', 'BE', 22),
('Poland', 'PL', 23),
('Portugal', 'PT', 24),
('South Korea', 'KR', 25),
('Singapore', 'SG', 26);

-- +goose Down
DROP TABLE IF EXISTS countries;
