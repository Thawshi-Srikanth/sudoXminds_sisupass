-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    wallet_id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    balance numeric DEFAULT 0 CHECK (balance >= 0),
    wallet_type text NOT NULL,
    parent_wallet_id uuid,
    access_qr_code text UNIQUE,
    access_nfc_code text UNIQUE,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT wallets_pkey PRIMARY KEY (wallet_id),
    CONSTRAINT wallets_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT wallets_parent_wallet_id_fkey FOREIGN KEY (parent_wallet_id) REFERENCES wallets(wallet_id)
);

-- Create indexes for better performance
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_wallet_type ON wallets(wallet_type);
CREATE INDEX idx_wallets_parent_wallet_id ON wallets(parent_wallet_id);
CREATE INDEX idx_wallets_created_at ON wallets(created_at);

-- Add comments for documentation
COMMENT ON TABLE wallets IS 'Digital wallets for users to store balance and access passes';
COMMENT ON COLUMN wallets.wallet_id IS 'Unique identifier for the wallet';
COMMENT ON COLUMN wallets.user_id IS 'Reference to the wallet owner';
COMMENT ON COLUMN wallets.balance IS 'Current balance in the wallet (cannot be negative)';
COMMENT ON COLUMN wallets.wallet_type IS 'Type of wallet (e.g., Bus Pass, Train Pass, Food Pass, etc.)';
COMMENT ON COLUMN wallets.parent_wallet_id IS 'Reference to parent wallet for hierarchical wallets';
COMMENT ON COLUMN wallets.access_qr_code IS 'QR code for wallet access containing structured data';
COMMENT ON COLUMN wallets.access_nfc_code IS 'NFC code for wallet access containing structured data';
COMMENT ON COLUMN wallets.created_at IS 'Timestamp when the wallet was created';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets CASCADE;
-- +goose StatementEnd
