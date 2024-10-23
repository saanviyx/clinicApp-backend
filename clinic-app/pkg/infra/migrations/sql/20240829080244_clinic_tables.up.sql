CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL UNIQUE PRIMARY KEY,
    username VARCHAR(25) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(25) NOT NULL,
    role VARCHAR(25) CHECK (role IN ('patient', 'doctor', 'admin')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Appointment (
    appointment_id SERIAL UNIQUE PRIMARY KEY,
    doctor_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    patient_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    appointment_date TIMESTAMP,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(50) CHECK (status IN ('scheduled', 'completed', 'canceled')) DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Slot (
    slot_id SERIAL UNIQUE PRIMARY KEY,
    appointment_id INT REFERENCES Appointment(appointment_id) ON DELETE CASCADE,
    doctor_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INTERVAL DEFAULT '0 hours',
    is_booked BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Schedules (
    schedule_id SERIAL PRIMARY KEY,
    doctor_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    date TIMESTAMP,
    total_appointment_time INTERVAL DEFAULT '0',
    total_appointments INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
