CREATE OR REPLACE FUNCTION update_slot_on_appointment()
RETURNS TRIGGER AS $$
DECLARE
    appointment_duration INTERVAL;
BEGIN
    appointment_duration := NEW.end_time - NEW.start_time;

    IF appointment_duration < '00:15:00' THEN
        RAISE EXCEPTION 'Appointment duration is too short. Minimum duration is 15 minutes.'
        USING ERRCODE = 'P0002';
    ELSIF appointment_duration > '02:00:00' THEN
        RAISE EXCEPTION 'Appointment duration exceeds the maximum limit. Maximum duration is 2 hours.'
        USING ERRCODE = 'P0003';
    END IF;

    -- Check if a slot already exists
    IF NOT EXISTS (
        SELECT 1 
        FROM Slot 
        WHERE doctor_id = NEW.doctor_id 
          AND start_time = NEW.start_time 
          AND end_time = NEW.end_time
    ) THEN
        INSERT INTO Slot (appointment_id, doctor_id, start_time, end_time, duration, is_booked, created_at)
        VALUES (
            NEW.appointment_id,
            NEW.doctor_id,
            NEW.start_time,
            NEW.end_time,
            appointment_duration, 
            TRUE, 
            CURRENT_TIMESTAMP 
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_slot_on_appointment
AFTER INSERT ON Appointment
FOR EACH ROW
EXECUTE FUNCTION update_slot_on_appointment();
