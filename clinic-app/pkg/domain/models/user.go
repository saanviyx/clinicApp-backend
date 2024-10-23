package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Doctor struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Availability string `json:"availability"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DoctorAvailability represents a doctor's availability
type DoctorAvailability struct {
	DoctorID          int    `json:"doctor_id"`
	DoctorName        string `json:"doctor_name"`
	DoctorEmail       string `json:"doctor_email"`
	Date              any    `json:"appointment_date"`
	TotalAppointments int    `json:"total_appointments"`
	TotalTime         string `json:"total_time"`
	Availability      string `json:"availability"`
}

// DoctorAppointments represents a doctor and their total appointments on a given day
type DoctorMostAppointments struct {
	DoctorID          int    `json:"doctor_id"`
	DoctorName        string `json:"doctor_name"`
	DoctorEmail       string `json:"doctor_email"`
	TotalAppointments int    `json:"total_appointments"`
}

// DoctorAppointments represents a doctor and their total appointments on a given day
type DoctorOverTime struct {
	DoctorID    int    `json:"doctor_id"`
	DoctorName  string `json:"doctor_name"`
	DoctorEmail string `json:"doctor_email"`
	TotalTime   string `json:"total_time"`
}
