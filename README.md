# Microservices Architecture for Vehicle Rental System

## Design Consideration of Your Microservices

### 1. **Separation of Concerns**

- The system is divided into several microservices, each responsible for specific domain logic:
  - **APIs**: Handles communication with the backend through RESTful APIs.
  - **Config**: Manages authentication, database connections, and environment configuration.
  - **Handlers**: Manages the logic and responses for the frontend application.
  - **Models**: Defines the database schema (structs) for different entities like users, cars, and payment.
  - **Services**: Contains functions for interacting with the database (insert, select, update) to manage user, car, and reservation and payment data.
  - **Static**: Stores static files like images, CSS, and JS used in the frontend.

### 2. **Scalability**

- Microservices are designed for independent scaling. The architecture ensures each service can scale based on load. Services like **APIs** and **Handlers** can be scaled horizontally to accommodate a high number of requests.

### 3. **Data Consistency**

- Communication between services is done via RESTful APIs. For consistency across microservices, the architecture might utilize event-driven communication or caching mechanisms where necessary.
- Each service uses a MySQL database for storing structured data such as users, vehicles, and reservations.

### 4. **Security**

- **Config** handles secure authentication, including JWT-based token management for user login and service-to-service communication.
- Passwords are stored securely using hashing algorithms, and sensitive data is encrypted.

### 5. **Fault Tolerance and Resilience**

- Each service implements basic error handling, retries, and fallback mechanisms.
- Load balancing is implemented to ensure requests are distributed evenly across the services.

### 6. **Inter-Service Communication**

- Services communicate using RESTful APIs to ensure easy integration and scalability. The API documentation is provided to help developers understand the endpoints available for interaction with the system.

---

## Architecture Diagram

```plaintext
+---------------------+        +---------------------+        +---------------------+
|       API.go        | <----> |   Handlers Service   | <----> |   Services          |
+---------------------+        +---------------------+        +---------------------+
         ^                           ^                             ^
         |                           |                             |
         |                           |                             |
+---------------------+        +---------------------+        +---------------------+
|   Config Service    |        |   Models            |        |   Static Files      |
+---------------------+        +---------------------+        +---------------------+
         |
         v
+---------------------+
|    Main.go          |
+---------------------+
```

# ECarShare Application Setup Guide

## Prerequisites

- Go (Golang) installed (version 1.20 or higher recommended)
- MySQL or compatible database server
- Git (optional, for cloning the repository)
- [godotenv](https://github.com/joho/godotenv) package for environment variable management

## Environment Configuration

1. Create a `.env` file in the project root directory
2. Add the following database configuration:

   ```
   DB_USER=user
   DB_PASSWORD=password
   DB_HOST=127.0.0.1:3306
   DB_NAME=ECarShare
   ```

3. Ensure `.env` is added to your `.gitignore` file to prevent sensitive information from being committed

## Database Setup

1. Connect to your MySQL database
2. Run the database initialization script:
   ```bash
   mysql -u [your_username] -p < ECarShare.sql
   ```
   Replace `[your_username]` with your database username
   You will be prompted to enter your database password

## Install Dependencies

```bash
go mod tidy
go get github.com/joho/godotenv
```

## Running the Application

Follow these steps in order:

### 1. Start the API Server

```bash
go run API.go
```

This will initialize the API endpoints and services

### 2. Run the Main Application

In a separate terminal window, run:

```bash
go run main.go
```

## Troubleshooting

- Verify database connection parameters in the `.env` file
- Ensure all required Go dependencies are installed
- Check that you have the necessary permissions to run database and Go scripts
- Confirm that the database is running and accessible

## Dependencies

- [godotenv](https://github.com/joho/godotenv) for environment variable management
- Database driver for your specific database system

# References

## Flatpickr

**APA-style citation:**  
Flatpickr. (n.d.). _A lightweight and powerful datetime picker_. Retrieved December 8, 2024, from [https://flatpickr.js.org/](https://flatpickr.js.org/)

---

## MDBootstrap (Material Design for Bootstrap)

MDBootstrap. (n.d.). _Material Design for Bootstrap_. Retrieved December 8, 2024, from [https://mdbootstrap.com/](https://mdbootstrap.com/)

---

## ChatGPT

OpenAI. (2024). _ChatGPT_ (Dec 2024 version). Retrieved from [https://chat.openai.com/](https://chat.openai.com/)

---

## Claude.ai

Anthropic. (n.d.). _Claude AI Assistant_. Retrieved December 8, 2024, from [https://claude.ai/](https://claude.ai/)

---

## Freepik

Freepik. (n.d.). _Free and Premium Graphic Resources_. Retrieved December 8, 2024, from [https://www.freepik.com/](https://www.freepik.com/)
