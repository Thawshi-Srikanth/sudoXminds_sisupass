-- +goose Up
-- +goose StatementBegin
CREATE TABLE slots (
    slot_id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    slot_type text NOT NULL,
    title text NOT NULL,
    description text,
    action text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    fields jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT slots_pkey PRIMARY KEY (slot_id),
    CONSTRAINT slots_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT slots_status_check CHECK (status IN ('active', 'inactive', 'expired', 'locked')),
    CONSTRAINT slots_slot_type_check CHECK (slot_type IN ('nfc', 'qr', 'auth', 'data', 'trigger', 'payment', 'access')),
    CONSTRAINT slots_action_check CHECK (action IN ('trigger', 'authenticate', 'authorize', 'data_read', 'data_write', 'payment', 'unlock', 'lock'))
);

-- Create indexes for better performance
CREATE INDEX idx_slots_user_id ON slots(user_id);
CREATE INDEX idx_slots_slot_type ON slots(slot_type);
CREATE INDEX idx_slots_action ON slots(action);
CREATE INDEX idx_slots_status ON slots(status);
CREATE INDEX idx_slots_created_at ON slots(created_at);
CREATE INDEX idx_slots_updated_at ON slots(updated_at);
CREATE INDEX idx_slots_title ON slots(title);

-- Create partial indexes for active slots
CREATE INDEX idx_slots_active ON slots(user_id, status) WHERE status = 'active';

-- Create GIN index for JSONB fields for better search performance
CREATE INDEX idx_slots_fields_gin ON slots USING GIN (fields);

-- Create text search index for title and description
CREATE INDEX idx_slots_text_search ON slots USING GIN (to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Add comments for documentation
COMMENT ON TABLE slots IS 'Interactive slots for NFC/QR actions and data management';
COMMENT ON COLUMN slots.slot_id IS 'Unique identifier for the slot';
COMMENT ON COLUMN slots.user_id IS 'Reference to the slot owner';
COMMENT ON COLUMN slots.slot_type IS 'Type of slot (nfc, qr, auth, data, trigger, payment, access)';
COMMENT ON COLUMN slots.title IS 'Display title for the slot';
COMMENT ON COLUMN slots.description IS 'Optional description explaining the slot purpose';
COMMENT ON COLUMN slots.action IS 'Action to execute when slot is triggered';
COMMENT ON COLUMN slots.status IS 'Current status of the slot (active, inactive, expired, locked)';
COMMENT ON COLUMN slots.fields IS 'Additional configuration data stored as JSON';
COMMENT ON COLUMN slots.created_at IS 'Timestamp when the slot was created';
COMMENT ON COLUMN slots.updated_at IS 'Timestamp when the slot was last updated';

-- Create trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_slots_updated_at 
    BEFORE UPDATE ON slots 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_slots_updated_at ON slots;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS slots;
-- +goose StatementEnd
