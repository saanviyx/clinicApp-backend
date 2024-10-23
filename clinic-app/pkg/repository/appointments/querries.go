package appointments

const (
	// Book an appointment
	BookAppointmentQuery = `
	WITH check_appointment AS (
		SELECT 1
    	FROM Appointment
    	WHERE doctor_id = $1
    	AND appointment_date = $3
    	AND start_time = $4
	),
	check_schedule AS (
    	SELECT total_appointment_time, total_appointments
    	FROM Schedules
    	WHERE doctor_id = $1
    	AND (date = $3 OR date IS NULL)
	),
	valid_duration AS (
    	SELECT
        	CASE
				WHEN COALESCE(total_appointments, 0) + 1 > 12 THEN FALSE
            	WHEN COALESCE(total_appointment_time, '00:00:00') + ($5 - $4) > '08:00:00' THEN FALSE
            	ELSE TRUE
        	END AS is_valid
    	FROM check_schedule
	),
	insert_appointment AS (
    	INSERT INTO Appointment (
        	doctor_id, 
        	patient_id, 
        	appointment_date, 
        	start_time, 
        	end_time
    	)
    	SELECT $1, $2, $3, $4, $5
    	WHERE NOT EXISTS (SELECT 1 FROM check_appointment)
    	AND EXISTS (SELECT 1 FROM check_schedule)
    	AND EXISTS (SELECT 1 FROM valid_duration WHERE is_valid = TRUE)
    	RETURNING appointment_id
	),
	valid_status AS (
    	SELECT 
        	CASE
            	WHEN EXISTS (SELECT 1 FROM check_appointment) THEN 'Appointment Exists'
            	WHEN NOT EXISTS (SELECT 1 FROM check_schedule) THEN 'Schedule Not Found'
            	WHEN EXISTS (SELECT 1 FROM valid_duration WHERE is_valid = FALSE) THEN 'Doctor Overbooked'
            	ELSE 'Valid'
        	END AS status
	)
	SELECT 
	(SELECT appointment_id FROM insert_appointment) AS appointment_id,
    (SELECT status FROM valid_status) AS status;
	`

	// View appointment details
	GetAppointmentByIdQuery = `
		SELECT 
			Appointment.appointment_id, 
			Appointment.patient_id,
			Patient.name AS patient_name,
			Doctor.name AS doctor_name,
			Appointment.start_time, 
			Appointment.end_time,
			Appointment.status
		FROM Appointment
		INNER JOIN Users AS Patient ON Appointment.patient_id = Patient.user_id
		INNER JOIN Users AS Doctor ON Appointment.doctor_id = Doctor.user_id
		WHERE Appointment.appointment_id = $1 
		AND (Patient.user_id = $2 OR Doctor.user_id = $2);
	`

	// View patient appointment history
	GetPatientAppointmentHistoryQuery = `
		SELECT 
			Appointment.appointment_id, 
			Patient.user_id AS patient_id,
			Doctor.name AS doctor_name,
			Patient.name AS patient_name,
			Appointment.start_time, 
			Appointment.end_time,
			Appointment.status
		FROM Appointment
		INNER JOIN Users AS Patient ON Appointment.patient_id = Patient.user_id
		INNER JOIN Users AS Doctor ON Appointment.doctor_id = Doctor.user_id
		WHERE Appointment.patient_id = $1
		ORDER BY Appointment.appointment_id DESC;
	`

	// Cancel an appointment
	CancelAppointmentQuery = `
		DELETE FROM Appointment
		WHERE appointment_id = $1;
	`

	// Delete slot on delete appointment
	DeleteSlotQuery = `
	DELETE FROM Slot 
		WHERE appointment_id = $1; 
	`
)
