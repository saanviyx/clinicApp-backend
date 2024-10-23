-- Reset sequence for primary keys
ALTER SEQUENCE users_user_id_seq RESTART WITH 1000;
ALTER SEQUENCE appointment_appointment_id_seq RESTART WITH 100;
ALTER SEQUENCE schedules_schedule_id_seq RESTART WITH 100;
ALTER SEQUENCE slot_slot_id_seq RESTART WITH 10;

ALTER TABLE Schedules
ADD COLUMN availability VARCHAR(15) DEFAULT 'available';


-- Function to insert a schedule when a new doctor is added
CREATE OR REPLACE FUNCTION insert_schedule_for_doctor()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the new user's role is 'doctor'
    IF NEW.role = 'doctor' THEN
        -- Insert a new schedule for the doctor
        INSERT INTO Schedules (
            doctor_id, 
            date, 
            total_appointment_time, 
            total_appointments, 
            availability, 
            created_at)
        VALUES (
            NEW.user_id,  
            NULL,       
            '00:00:00',   
            0,          
            'available',   
            CURRENT_TIMESTAMP
        );
    END IF;
    -- Return the new row
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger that calls the above function after a new user is inserted
CREATE TRIGGER trigger_insert_schedule_for_doctor
AFTER INSERT ON Users
FOR EACH ROW
EXECUTE FUNCTION insert_schedule_for_doctor();



CREATE OR REPLACE FUNCTION update_schedule_record()
RETURNS TRIGGER AS $$
DECLARE
    appointment_duration INTERVAL;
BEGIN
    appointment_duration := NEW.end_time - NEW.start_time;

    IF EXISTS (
        SELECT 1
        FROM Schedules
        WHERE doctor_id = NEW.doctor_id
        AND (date = NEW.appointment_date OR date IS NULL)
    ) THEN
        -- Update the schedule with the new appointment details and correct the date if necessary
        UPDATE Schedules
        SET total_appointment_time = total_appointment_time + appointment_duration,
            total_appointments = total_appointments + 1,
            date = COALESCE(date, NEW.appointment_date),
            availability = CASE
                WHEN total_appointments + 1 >= 12 OR total_appointment_time + appointment_duration >= '08:00:00'
                THEN 'unavailable'
                ELSE 'available'
            END
        WHERE doctor_id = NEW.doctor_id
        AND (date = NEW.appointment_date OR date IS NULL);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_schedule_record
AFTER INSERT ON Appointment
FOR EACH ROW
EXECUTE FUNCTION update_schedule_record();



-- Function to update schedule metrics when an appointment is cancelled
CREATE OR REPLACE FUNCTION update_schedule_on_cancellation()
RETURNS TRIGGER AS $$
DECLARE
    appointment_duration INTERVAL;
BEGIN
    SELECT end_time - start_time INTO appointment_duration
    FROM Appointment
    WHERE appointment_id = OLD.appointment_id;

    UPDATE Schedules
    SET total_appointment_time = total_appointment_time - appointment_duration,
        total_appointments = total_appointments - 1,
        availability = CASE
                        WHEN total_appointments = 12 OR total_appointment_time + appointment_duration = '08:00:00'
                        THEN 'unavailable'
                        ELSE 'available'
                       END
    WHERE doctor_id = OLD.doctor_id
    AND date = OLD.appointment_date;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update schedule metrics before deleting an appointment
CREATE TRIGGER trigger_update_schedule_on_cancellation
BEFORE DELETE ON Appointment
FOR EACH ROW
EXECUTE FUNCTION update_schedule_on_cancellation();

