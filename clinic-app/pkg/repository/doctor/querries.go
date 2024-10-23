package doctor

const (
	// View all doctors
	GetAllDoctorsQuery = `
		SELECT 
			Users.user_id AS doctor_id, 
			Users.name, 
			Users.email,
			Schedules.availability
		FROM Users
		INNER JOIN Schedules ON Users.user_id = Schedules.doctor_id
		WHERE Users.role = 'doctor'
		AND Schedules.availability = 'available';
	`

	// View specific doctor information
	GetDoctorByIdQuery = `
		SELECT 
			Users.user_id AS doctor_id, 
			Users.name, 
			Users.email,
			Schedules.availability
		FROM Users
		INNER JOIN Schedules ON Users.user_id = Schedules.doctor_id
		WHERE user_id = $1 
		AND role = 'doctor'
		AND Schedules.availability = 'available';
	`
	// View available time slots by doctor
	GetSlotsByDoctorQuery = `
		SELECT 
			s.slot_id,
			s.appointment_id, 
			p.user_id,
			p.name,
			s.start_time, 
			s.end_time, 
			s.is_booked,
			s.duration
		FROM Slot s
		LEFT JOIN Appointment a ON s.doctor_id = a.doctor_id AND s.start_time = a.start_time
		LEFT JOIN Users p ON a.patient_id = p.user_id
		WHERE s.doctor_id = $1
	`

	// View available time slots by patient
	GetSlotsByPatientQuery = `
		SELECT 
			s.slot_id,
			s.appointment_id, 
			s.start_time, 
			s.end_time, 
			s.is_booked,
			s.duration
		FROM Slot s
		LEFT JOIN Appointment a ON s.doctor_id = a.doctor_id AND s.start_time = a.start_time
		LEFT JOIN Users p ON a.patient_id = p.user_id
		WHERE s.doctor_id = $1;
	`
)
