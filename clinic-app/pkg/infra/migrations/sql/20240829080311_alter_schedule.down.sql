ALTER TABLE Schedules
DROP COLUMN IF EXISTS total_appointment_time,
DROP COLUMN IF EXISTS total_appointments;

-- Drop Triggers
DROP TRIGGER IF EXISTS trigger_update_schedule ON Users CASCADE;
DROP TRIGGER IF EXISTS trigger_update_schedule_record ON Appointment CASCADE;
DROP TRIGGER IF EXISTS trigger_update_schedule_on_cancellation ON Appointment CASCADE;

-- Drop Functions
DROP FUNCTION IF EXISTS update_schedule() CASCADE;
DROP FUNCTION IF EXISTS update_schedule_record() CASCADE;
DROP FUNCTION IF EXISTS update_schedule_on_cancellation() CASCADE;
