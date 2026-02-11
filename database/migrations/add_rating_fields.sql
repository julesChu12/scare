-- Add rating fields to service_requests table
ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS rating INT DEFAULT 0;
ALTER TABLE service_requests ADD COLUMN IF NOT EXISTS feedback TEXT;
