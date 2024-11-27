/**-----------------------------------------------**/
/** 			 CNAD ASSIGNMENT  1				  **/
/**				CAR SHARING SYSTEM				  **/
/**												  **/
/**  	      Database Script for creating  	  **/
/**  			database tables and data.		  **/
/**						  **/
/**-----------------------------------------------**/

/**============Create the Tables==================**/

-- Drop database if already exists
DROP DATABASE IF EXISTS ECarShare;
CREATE DATABASE ECarShare;

use ECarShare;

-- Drop tables if already exist
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cars;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS pricing;
DROP TABLE IF EXISTS payment;

/**=======Table:  users=======**/ 
CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,						-- Unique identifier for the user
    email VARCHAR(255) UNIQUE,									-- Email for registration (ensure uniqueness)
    phone VARCHAR(8) UNIQUE,									-- Phone number for registration (optional, ensure uniqueness)
    password_hash VARCHAR(255) NOT NULL,						-- Encrypted password for authentication
    membership_tier ENUM('Basic', 'Premium', 'VIP') DEFAULT 'Basic',  -- Membership tier, default is 'Basic'
    name VARCHAR(100),											-- User's name
	username VARCHAR(50),										-- User's username (to log in)
    date_of_birth DATE,											-- Date of birth
    is_verified BOOLEAN DEFAULT FALSE							-- Whether the user has verified their email/phone
);

/**=======Table:  cars=======**/ 
CREATE TABLE cars (
    car_id INT AUTO_INCREMENT PRIMARY KEY,						-- Unique identifier for the car
    car_model VARCHAR(255),										-- Model of the car
    license_plate VARCHAR(50) UNIQUE,							-- License plate number (unique)
    status ENUM('Available', 'Reserved', 'In Maintenance', 'Unavailable') DEFAULT 'Available', -- car status
    current_location VARCHAR(255),								-- Current location of the car (could be GPS coordinates)
    charge_level INT DEFAULT 100,								-- Battery charge level (if electric), default 100%
    cleanliness_status ENUM('Clean', 'Needs Cleaning', 'Dirty') DEFAULT 'Clean', -- Cleanliness of the car
    last_serviced TIMESTAMP,									-- Timestamp when the car was last serviced
    added_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP				-- Timestamp when the car was added to the system
);

/**=======Table:  rentals=======**/ 
CREATE TABLE rentals (
    rental_id INT AUTO_INCREMENT PRIMARY KEY,              -- Unique rental identifier
    user_id INT,                                           -- Foreign key referencing the users table
    car_id INT,                                        -- Foreign key referencing the vehicles table
    start_datetime DATETIME,                               -- Start date and time of the rental
    end_datetime DATETIME,                                 -- End date and time of the rental
    total_amount DECIMAL(10, 2),                            -- Total amount after discount
    status ENUM('Pending', 'Completed', 'Cancelled') DEFAULT 'Pending', -- Rental status
    FOREIGN KEY (user_id) REFERENCES users(user_id),       -- Foreign key referencing the users table
    FOREIGN KEY (car_id) REFERENCES cars(car_id) -- Foreign key referencing the vehicles table
);


/**=======Table:  reservations=======**/ 
CREATE TABLE reservations (
    reservation_id INT AUTO_INCREMENT PRIMARY KEY,           -- Unique identifier for the reservation
    user_id INT,                                             -- Reference to the user making the reservation
    car_id INT,												 -- Reference to the car being reserved
    start_datetime DATETIME,                                 -- Start date and time of the reservation
    end_datetime DATETIME,                                   -- End date and time of the reservation
    status ENUM('Pending', 'Confirmed', 'Cancelled', 'Completed') DEFAULT 'Pending', -- Status of the reservation
    FOREIGN KEY (user_id) REFERENCES users(user_id),         -- Foreign key referencing the users table
    FOREIGN KEY (car_id) REFERENCES cars(car_id)			 -- Foreign key referencing the cars table
);

/**=======Table:  pricing=======**/ 
CREATE TABLE pricing (
    pricing_id INT AUTO_INCREMENT PRIMARY KEY,              -- Unique identifier for pricing record
    hourly_rate DECIMAL(10, 2),                             -- Hourly rate for the car
    weekend_rate DECIMAL(10, 2),                            -- Special weekend rate (if applicable)
    discount_percentage DECIMAL(5, 2) DEFAULT 0            	-- Discount percentage (e.g., 10.00 for 10% discount)
);

/**=======Table:  payment=======**/ 
CREATE TABLE payments (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,             	-- Unique payment identifier
    rental_id INT,                                         	-- Foreign key referencing the rentals table
    user_id INT,                                           	-- Foreign key referencing the users table
    payment_method ENUM('Credit Card', 'Debit Card', 'PayPal', 'Apple Pay') NOT NULL, -- Payment method
    payment_status ENUM('Pending', 'Completed', 'Failed') DEFAULT 'Pending', -- Payment status
    amount DECIMAL(10, 2),                                  -- Amount paid
	transaction_id CHAR(36) NOT NULL,						-- unique identifier for transaction, i.e. guids
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,       -- Date and time of the payment
    FOREIGN KEY (rental_id) REFERENCES rentals(rental_id), 	-- Foreign key referencing the rentals table
    FOREIGN KEY (user_id) REFERENCES users(user_id)        	-- Foreign key referencing the users table
);


/**  Creating Sample Data  **/

/** Creating Records for Table users **/

INSERT INTO users (email, phone, password_hash, membership_tier, name, username, date_of_birth, is_verified)
VALUES ('john.doe@example.com', '91234567', 'hashed_password_example', 'Premium', 'John Doe', 'johndoe', '1990-01-01', TRUE);


/** Creating Records for Table cars **/
INSERT INTO cars (car_model, license_plate, status, current_location, charge_level, cleanliness_status, last_serviced)
VALUES ('Tesla Model S', 'XYZ1234AB', 'Available', '40.7128° N, 74.0060° W', 85, 'Clean', '2024-11-01 14:30:00');

/** Creating Records for Table rentals **/
INSERT INTO rentals (user_id, car_id, start_datetime, end_datetime, total_amount, status)
VALUES (1, 1, '2024-12-01 09:00:00', '2024-12-01 17:00:00', 100.00, 'Completed');

/** Creating Records for Table reservations **/
INSERT INTO reservations (user_id, car_id, start_datetime, end_datetime, status)
VALUES (1, 1, '2024-12-01 10:00:00', '2024-12-01 14:00:00', 'Confirmed');

/** Creating Records for Table pricing **/
INSERT INTO pricing (hourly_rate, weekend_rate, discount_percentage)
VALUES (25.00, 30.00, 10.00);

/** Creating Records for Table payment **/
INSERT INTO payments (rental_id, user_id, payment_method, payment_status, amount, transaction_id)
VALUES (1, 1, 'Credit Card', 'Completed', 150.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479');


select*from users;
select*from rentals;
select*from reservations;
select*from payments;
select*from pricing;