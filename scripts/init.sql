-- Enable UUID extension for PostgreSQL
-- This script ensures UUID functions are available

-- Try to create extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Alternative: Enable crypto extension (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create custom function as fallback
DO $$
BEGIN
    -- Check if gen_random_uuid exists, if not create it
    IF NOT EXISTS (
        SELECT 1 FROM pg_proc
        WHERE proname = 'gen_random_uuid'
    ) THEN
        -- Create gen_random_uuid function using uuid-ossp
        CREATE OR REPLACE FUNCTION gen_random_uuid()
        RETURNS uuid AS $$
        BEGIN
            RETURN uuid_generate_v4();
        END;
        $$ LANGUAGE plpgsql;
    END IF;
END $$;