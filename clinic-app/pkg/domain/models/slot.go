package models

import "time"

type SlotDoc struct {
	SlotID        int       `json:"slot_id"`
	AppointmentID int       `json:"appointment_id"`
	PatientID     int       `json:"patient_id"`
	PatientName   string    `json:"patient_name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsBooked      bool      `json:"is_booked"`
	Duration      string    `json:"duration"`
}

type SlotPat struct {
	SlotID        int       `json:"slot_id"`
	AppointmentID int       `json:"appointment_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsBooked      bool      `json:"is_booked"`
	Duration      string    `json:"duration"`
}
