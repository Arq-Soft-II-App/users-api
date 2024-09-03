# Users API

This project is a **Users API** built with **Golang**, providing a well-structured, scalable architecture for managing user data. The API leverages several modern tools and design patterns to ensure maintainability, efficiency, and ease of extension.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete operations for user entities.
- **MongoDB Integration**: Persistent data storage using MongoDB with a Singleton pattern to manage database connections.
- **RESTful API**: Built with the Gin framework, providing fast and easy routing for handling HTTP requests.
- **Data Transfer Objects (DTOs)**: Used for validation and structuring of input/output data between layers.
- **Middleware**: Custom middleware for API Key authentication to secure access to the API.
- **Error Handling**: Robust error handling across all layers, ensuring meaningful responses and logging.
- **Environment Configuration**: Environment variables managed using `godotenv` for easy configuration and deployment.

## Architecture

The API is designed with a clean and modular architecture, ensuring separation of concerns and ease of maintenance:

- **Router**: Manages API routes and applies middleware.
- **Controller**: Handles HTTP requests, interacts with the service layer, and returns responses.
- **Service**: Contains business logic, processes data from the controller, and interacts with the repository.
- **Repository**: Responsible for interacting directly with MongoDB, performing CRUD operations.
- **DTOs (Data Transfer Objects)**: Ensures data validation and integrity between the API's layers.

## Tools and Technologies

- **Golang**: The core programming language used to build the API.
- **Gin**: A high-performance HTTP web framework for routing and middleware.
- **MongoDB**: NoSQL database used for persistent storage.
- **Mongo Driver**: Official MongoDB driver for Golang, used to manage database operations.
- **godotenv**: Manages environment variables, simplifying the configuration process.
- **bcrypt**: Used for hashing passwords securely.
- **Sync/Once**: Singleton pattern for managing a single instance of the database connection.
