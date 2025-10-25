-- Enable UUID extension for PostgreSQL
-- This script ensures UUID functions are available

-- Try to create extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Alternative: Enable crypto extension (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";