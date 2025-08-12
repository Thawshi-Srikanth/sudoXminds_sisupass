CREATE TABLE access_tokens (
    token_id uuid NOT NULL DEFAULT gen_random_uuid(),
    wallet_id uuid,
    public_token text NOT NULL UNIQUE,
    passkey_hash text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    revoked_at timestamp without time zone,
    CONSTRAINT access_tokens_pkey PRIMARY KEY (token_id),
    CONSTRAINT access_tokens_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id)
);
CREATE TABLE bookings (
    booking_id uuid NOT NULL DEFAULT gen_random_uuid(),
    wallet_id uuid,
    slot_id uuid,
    booking_date timestamp with time zone DEFAULT now(),
    status text CHECK (
        status = ANY (ARRAY ['pending', 'confirmed', 'cancelled'])
    ),
    details jsonb,
    payment_transaction_id uuid,
    CONSTRAINT bookings_pkey PRIMARY KEY (booking_id),
    CONSTRAINT bookings_slot_id_fkey FOREIGN KEY (slot_id) REFERENCES slots(slot_id),
    CONSTRAINT bookings_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id),
    CONSTRAINT bookings_payment_transaction_id_fkey FOREIGN KEY (payment_transaction_id) REFERENCES transactions(transaction_id)
);
CREATE TABLE documents (
    document_id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid,
    name text NOT NULL,
    document_type text NOT NULL,
    file_url text NOT NULL,
    file_name text,
    uploaded_at timestamp with time zone DEFAULT now(),
    metadata jsonb,
    CONSTRAINT documents_pkey PRIMARY KEY (document_id),
    CONSTRAINT documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE scan_logs (
    log_id uuid NOT NULL DEFAULT gen_random_uuid(),
    token_id uuid,
    scanned_at timestamp with time zone DEFAULT now(),
    scanned_by_wallet_id uuid,
    scan_method text CHECK (scan_method = ANY (ARRAY ['qr', 'nfc'])),
    ip_address inet,
    location jsonb,
    CONSTRAINT scan_logs_pkey PRIMARY KEY (log_id),
    CONSTRAINT scan_logs_token_id_fkey FOREIGN KEY (token_id) REFERENCES access_tokens(token_id),
    CONSTRAINT scan_logs_scanned_by_wallet_id_fkey FOREIGN KEY (scanned_by_wallet_id) REFERENCES wallets(wallet_id)
);
CREATE TABLE sessions (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    refresh_token text NOT NULL UNIQUE,
    revoked boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    expires_at timestamp with time zone NOT NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE slots (
    slot_id uuid NOT NULL DEFAULT gen_random_uuid(),
    slot_type text NOT NULL,
    title text NOT NULL,
    description jsonb,
    action jsonb,
    fields jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT slots_pkey PRIMARY KEY (slot_id)
);
CREATE TABLE tokens (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    token_type text NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    expires_at timestamp with time zone NOT NULL,
    CONSTRAINT tokens_pkey PRIMARY KEY (id),
    CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE transactions (
    transaction_id uuid NOT NULL DEFAULT gen_random_uuid(),
    from_wallet_id uuid,
    to_wallet_id uuid,
    amount numeric NOT NULL CHECK (amount > 0),
    transaction_date timestamp with time zone DEFAULT now(),
    transaction_type text NOT NULL CHECK (
        transaction_type = ANY (
            ARRAY ['topup', 'spend', 'transfer', 'booking_payment']
        )
    ),
    status text DEFAULT 'completed' CHECK (
        status = ANY (ARRAY ['pending', 'completed', 'failed'])
    ),
    description text,
    CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id),
    CONSTRAINT transactions_to_wallet_id_fkey FOREIGN KEY (to_wallet_id) REFERENCES wallets(wallet_id),
    CONSTRAINT transactions_from_wallet_id_fkey FOREIGN KEY (from_wallet_id) REFERENCES wallets(wallet_id)
);
CREATE TABLE user_identities (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    provider text NOT NULL,
    provider_user_id text NOT NULL,
    access_token text,
    refresh_token text,
    token_expiry timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT user_identities_pkey PRIMARY KEY (id),
    CONSTRAINT user_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE users (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    email character varying NOT NULL UNIQUE,
    encrypted_password character varying,
    role character varying NOT NULL DEFAULT 'user',
    is_active boolean DEFAULT true,
    email_confirmed_at timestamp with time zone,
    phone character varying UNIQUE,
    phone_confirmed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
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