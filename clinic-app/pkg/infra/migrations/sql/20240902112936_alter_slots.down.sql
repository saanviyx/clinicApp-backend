DROP TRIGGER IF EXISTS trigger_update_slot_on_appointment ON Appointment CASCADE;

DROP FUNCTION IF EXISTS update_slot_on_appointment() CASCADE;