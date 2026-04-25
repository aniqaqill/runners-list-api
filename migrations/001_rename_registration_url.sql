-- Migration: Rename registeration_url → registration_url
-- =============================================================================
-- Context
-- -------
-- The original Go domain struct used the misspelled field "RegisterationURL".
-- GORM mapped this to the column "registeration_url" in Postgres.
-- The field has been renamed to "RegistrationURL" in Go code, and a
-- gorm:"column:registeration_url" tag keeps the old column name in production.
--
-- This migration renames the DB column to match the corrected spelling.
-- After running this, remove the explicit column tag from domain/events.go.
--
-- Safety pattern (dual-write migration):
--   Step 1: Add new column (this script).
--   Step 2: Deploy app with dual-write to both columns.
--   Step 3: Back-fill old rows.
--   Step 4: Switch reads to new column, remove old.
--
-- Since this is a single-instance hobby app with low traffic you can run
-- the simpler one-step rename below with near-zero risk.
-- =============================================================================

-- ONE-STEP (simple rename, zero downtime for this scale):
ALTER TABLE events
  RENAME COLUMN registeration_url TO registration_url;

-- After running: remove the gorm:"column:registeration_url" tag from
-- internal/core/domain/events.go so GORM derives the column name naturally.
