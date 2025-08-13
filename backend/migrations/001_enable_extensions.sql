-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- +goose StatementBegin
SELECT 'pgcrypto extension enabled';
-- +goose StatementEnd

-- +goose Down
DROP EXTENSION IF EXISTS "pgcrypto";

-- +goose StatementBegin
SELECT 'pgcrypto extension disabled';
-- +goose StatementEnd