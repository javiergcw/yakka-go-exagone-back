-- Create custom types for Yakka Backend
-- This script creates the required custom types for the application

-- Create user_status enum type
DO $$ BEGIN
    CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');
    RAISE NOTICE 'Created user_status enum type';
EXCEPTION
    WHEN duplicate_object THEN 
        RAISE NOTICE 'user_status enum type already exists';
END $$;

-- Create user_role enum type
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');
    RAISE NOTICE 'Created user_role enum type';
EXCEPTION
    WHEN duplicate_object THEN 
        RAISE NOTICE 'user_role enum type already exists';
END $$;

-- Verify types were created
SELECT typname, typtype FROM pg_type WHERE typname IN ('user_status', 'user_role');
