-- Add role_id foreign key column to users table with proper constraints
-- This migration establishes the relationship between users and roles

-- Step 1: Add the role_id column with a temporary default value
-- The default value of 1 allows existing users to have a valid role during migration
ALTER TABLE IF EXISTS users
ADD COLUMN role_id INT REFERENCES roles(id) DEFAULT 1;

-- Step 2: Update all existing users to have the 'user' role
-- This ensures all existing users get assigned to the default user role
UPDATE users 
SET role_Id = (
      SELECT id FROM roles WHERE name = 'user'
);

-- Step 3: Remove the default value constraint
-- After all users have been assigned roles, we no longer need the default
ALTER TABLE users
ALTER COLUMN role_id DROP DEFAULT; 

-- Step 4: Make role_id NOT NULL
-- Now that all users have valid role assignments, enforce the NOT NULL constraint
ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;