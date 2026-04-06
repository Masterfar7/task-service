-- Add recurrence fields to tasks table
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS is_template BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS recurrence_type TEXT NOT NULL DEFAULT 'none';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS recurrence_config JSONB;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS next_occurrence DATE;

-- Create index for template tasks
CREATE INDEX IF NOT EXISTS idx_tasks_is_template ON tasks (is_template) WHERE is_template = TRUE;

-- Create index for next occurrence
CREATE INDEX IF NOT EXISTS idx_tasks_next_occurrence ON tasks (next_occurrence) WHERE next_occurrence IS NOT NULL;

-- Add check constraint for recurrence_type
ALTER TABLE tasks ADD CONSTRAINT check_recurrence_type
    CHECK (recurrence_type IN ('none', 'daily', 'monthly', 'specific_dates', 'even_odd'));
