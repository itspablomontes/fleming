-- Auth Challenges table for SIWE nonce storage
-- Challenges are short-lived and cleaned up periodically

CREATE TABLE IF NOT EXISTS auth_challenges (
    address     VARCHAR(42) PRIMARY KEY,
    message     TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Index for efficient expired challenge cleanup
CREATE INDEX IF NOT EXISTS idx_auth_challenges_expires ON auth_challenges (expires_at);

COMMENT ON TABLE auth_challenges IS 'Stores SIWE challenge messages pending signature verification';
COMMENT ON COLUMN auth_challenges.address IS 'Ethereum wallet address (checksummed)';
COMMENT ON COLUMN auth_challenges.message IS 'Full EIP-4361 SIWE message to be signed';
COMMENT ON COLUMN auth_challenges.expires_at IS 'Challenge expiration time (5 minutes from creation)';
