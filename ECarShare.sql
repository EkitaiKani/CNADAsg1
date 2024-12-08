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
DROP TABLE IF EXISTS payments;

/**=======Table:  users=======**/ 
CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,						-- Unique identifier for the user
    email VARCHAR(255) UNIQUE,									-- Email for registration (ensure uniqueness)
    -- phone VARCHAR(8) UNIQUE,									-- Phone number for registration (optional, ensure uniqueness)
	username VARCHAR(50) UNIQUE,										-- User's username (to log in)
    password_hash VARCHAR(255) NOT NULL,						-- Encrypted password for authentication
    membership_tier ENUM('Basic', 'Premium', 'VIP') DEFAULT 'Basic',  -- Membership tier, default is 'Basic'
    firstname VARCHAR(100),											-- User's first name
    lastname VARCHAR(100),											-- User's last name
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
    added_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,				-- Timestamp when the car was added to the system
    rate INT													-- hourly rate of the car
);

/**=======Table:  reservations=======**/ 
CREATE TABLE reservations (
    reservation_id INT AUTO_INCREMENT PRIMARY KEY,           -- Unique identifier for the reservation
    user_id INT,                                             -- Reference to the user making the reservation
    car_id INT,												 -- Reference to the car being reserved
    start_datetime DATETIME,                                 -- Start date and time of the reservation
    end_datetime DATETIME,                                   -- End date and time of the reservation
    status ENUM('Pending', 'Confirmed', 'Ongoing', 'Cancelled', 'Completed') DEFAULT 'Pending', -- Status of the reservation
    FOREIGN KEY (user_id) REFERENCES users(user_id),         -- Foreign key referencing the users table
    FOREIGN KEY (car_id) REFERENCES cars(car_id)			 -- Foreign key referencing the cars table
);

/**=======Table:  MembershipTiers =======**/ 
CREATE TABLE MembershipTiers (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for discount record
    tier_name VARCHAR(50) NOT NULL, -- name of the tier
    discount_percentage INT NOT NULL -- Discount percentage (e.g., 10 for 10% discount)
);


/**=======Table:  payment=======**/ 
CREATE TABLE payments (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,             	-- Unique payment identifier
    reservation_id INT,                                         	-- Foreign key referencing the rentals table
    user_id INT,                                           	-- Foreign key referencing the users table
    payment_method ENUM('Credit Card', 'Debit Card', 'PayPal', 'Apple Pay', 'Pending'), -- Payment method
    payment_status ENUM('Pending', 'Completed', 'Failed') DEFAULT 'Pending', -- Payment status
    amount DECIMAL(10, 2),                                  -- Amount paid
	transaction_id CHAR(36) NOT NULL,						-- unique identifier for transaction, i.e. guids
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,       -- Date and time of the payment
    FOREIGN KEY (reservation_id) REFERENCES reservations(reservation_id), 	-- Foreign key referencing the rentals table
    FOREIGN KEY (user_id) REFERENCES users(user_id)        	-- Foreign key referencing the users table
);


/**  Creating Sample Data  **/

/** Creating Records for Table users **/

INSERT INTO users (email, password_hash, membership_tier, firstname, lastname, username, date_of_birth, is_verified)
VALUES ('john.doe@example.com', 'hashed_password_example', 'Premium', 'John', 'Doe', 'johndoe', '1990-01-01', TRUE);


/** Creating Records for Table cars **/
INSERT INTO cars (car_model, license_plate, status, current_location, charge_level, cleanliness_status, last_serviced, rate)
VALUES
('Tesla Model S', 'XYZ1234AB', 'Available', '40.7128° N, 74.0060° W', 85, 'Clean', '2024-11-01 14:30:00', 15),
('Honda Civic', 'ABC5678CD', 'Reserved', '34.0522° N, 118.2437° W', 90, 'Clean', '2024-10-15 09:00:00', 10),
('Ford Mustang', 'LMN1234OP', 'In Maintenance', '51.5074° N, 0.1278° W', 75, 'Needs Cleaning', '2024-09-10 18:00:00', 20),
('Chevrolet Bolt', 'DEF2345GH', 'Available', '37.7749° N, 122.4194° W', 100, 'Clean', '2024-08-05 11:30:00', 18),
('BMW 3 Series', 'JKL3456MN', 'Unavailable', '48.8566° N, 2.3522° E', 50, 'Dirty', '2024-07-20 13:15:00', 25),
('Audi Q5', 'PQR9876ST', 'Reserved', '52.3676° N, 4.9041° E', 65, 'Clean', '2024-10-30 12:00:00', 22);

/** Creating Records for Table payment **/
INSERT INTO payments (reservation_id, user_id, payment_method, payment_status, amount, transaction_id)
VALUES (1, 2, 'Credit Card', 'Completed', 150.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479');

INSERT INTO MembershipTiers (tier_name, discount_percentage) 
VALUES 
    ('Basic', 0),
    ('Premium', 5),
    ('VIP', 10);
    
use ECarShare;

select*from users;
select*from cars;
select*from MembershipTiers;
select*from reservations;
select*from payments;

