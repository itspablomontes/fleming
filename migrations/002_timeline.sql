-- Medical Timeline Events
-- Stores chronological historical health data for patients

CREATE TABLE IF NOT EXISTS timeline_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id   VARCHAR(42) NOT NULL,
    type         VARCHAR(50) NOT NULL,
    title        TEXT NOT NULL,
    description  TEXT,
    provider     TEXT,
    timestamp    TIMESTAMPTZ NOT NULL,
    blob_ref     TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    metadata     JSONB DEFAULT '{}',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

-- Index for efficient chronological retrieval per patient
CREATE INDEX IF NOT EXISTS idx_timeline_patient_timestamp ON timeline_events (patient_id, timestamp DESC);

COMMENT ON TABLE timeline_events IS 'Stores clinical events and documents in a patient timeline';
COMMENT ON COLUMN timeline_events.patient_id IS 'Ethereum wallet address of the data owner';
COMMENT ON COLUMN timeline_events.blob_ref IS 'Reference to encrypted blob in storage (e.g. MinIO/IPFS)';
