package models

import "time"

type BookAppointment struct {
	AppointmentID int       `json:"appointment_id"`
	DoctorID      int       `json:"doctor_id"`
	PatientID     int       `json:"patient_id"`
	Date          time.Time `json:"appointment_date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
}

type Appointment struct {
	AppointmentID int       `json:"appointment_id"`
	PatientID     int       `json:"patient_id"`
	PatientName   string    `json:"patient_name"`
	DoctorName    string    `json:"doctor_name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
}
