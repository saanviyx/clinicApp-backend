package admin

const (
	// View availability of all doctors
	GetAllDoctorsAvailabilityQuery = `
		SELECT 
			Users.user_id, 
			Users.name, 
			Users.email,
			Schedules.date, 
			Schedules.total_appointments, 
			Schedules.total_appointment_time,
			Schedules.availability
		FROM Users
		INNER JOIN Schedules ON Users.user_id = Schedules.doctor_id;
	`

	// View doctors with the most appointments in a day
	GetDoctorsWithMostAppointmentsQuery = `
		SELECT 
            Users.user_id, 
            Users.name,
            Users.email, 
            Schedules.total_appointments
        FROM Users
        INNER JOIN Schedules ON Users.user_id = Schedules.doctor_id
        WHERE Schedules.date = $1
        GROUP BY Users.user_id, Schedules.total_appointments
        ORDER BY total_appointments DESC;
	`

	// View doctors with over 6 hours of appointments in a day
	GetDoctorsWithOverSixHoursQuery = `
		SELECT 
			Users.user_id, 
			Users.name, 
			Users.email,
			Schedules.total_appointment_time
		FROM Users
		INNER JOIN Schedules ON Users.user_id = Schedules.doctor_id
		WHERE Schedules.date = $1
		GROUP BY Users.user_id, Schedules.total_appointment_time
		HAVING Schedules.total_appointment_time > '06:00:00';
	`
)
