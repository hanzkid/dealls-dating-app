# Dealls Dating App Documentation

## Repository

- **URL**: [Dealls Dating App GitHub Repository](https://github.com/hanzkid/dealls-dating-app)

---

## Functional Requirements

### 1. **User Authentication**

Endpoints for managing user access and credentials:

- **Login**  
  - **Endpoint**: `/login`  
  - **Method**: POST  
  - **Description**: Allows users to log in to their accounts.  
  - **Security**: JWT-based authentication is issued upon successful login.  

- **Register**  
  - **Endpoint**: `/register`  
  - **Method**: POST  
  - **Description**: Enables new users to create accounts.  
  - **Password Storage**: Passwords are securely hashed using bcrypt.

---

### 2. **User Profile Management**

Endpoints for viewing and managing user profiles:

- **View Profile**  
  - **Endpoint**: `/me`  
  - **Method**: GET  
  - **Description**: Returns the authenticated user's profile information.  

- **Update Profile**  
  - **Endpoint**: `/me`  
  - **Method**: PUT  
  - **Description**: Allows users to update their profile information.  

- **Subscribe to Premium Services**  
  - **Endpoint**: `/subscribe`  
  - **Method**: POST  
  - **Description**: Enables users to subscribe to premium features.

---

### 3. **Dating Features**

Core functionalities for matchmaking and interaction:

- **View Profiles**  
  - **Endpoint**: `/profile`  
  - **Method**: GET  
  - **Description**: Displays a random user profile available for interaction. Free users can view up to 10 profiles, while premium users have unlimited access

- **Swipe Profiles**  
  - **Endpoint**: `/swipe`  
  - **Method**: POST  
  - **Description**: Allows users to swipe (like or dislike) other profiles.  

- **View Matches**  
  - **Endpoint**: `/match`  
  - **Method**: GET  
  - **Description**: Returns a list of profiles that mutually liked the authenticated user.

---

## Non-Functional Requirements

### 1. **Security**

- **Authentication**:  
  - JWT-based authentication is implemented for all protected routes to ensure secure access.  
- **Password Storage**:  
  - Passwords are securely hashed using **bcrypt** before being stored in the database.

---

### 2. **Testing**

- Comprehensive unit and integration tests were implemented to ensure system reliability and correctness.  
- Tests cover:
  - Authentication flow
  - Profile management
  - Core dating features (e.g., swipe, match retrieval)

---

## TechStack

1. **Golang**: Chosen for its performance, simplicity, and ability to handle high-concurrency applications efficiently.  
2. **JWT-based authentication**: Provides a stateless and secure method for authenticating users across endpoints.  
3. **Bcrypt hash**: Ensures secure storage of passwords by employing a strong hashing algorithm with built-in salting.  
4. **PostgreSQL**: Offers reliability, scalability, and advanced querying capabilities for managing relational data.  
5. **Docker containerization**: Simplifies deployment by creating consistent environments across different systems.


## Design System

### Entity Relationship Diagram ( ERD )
![image](https://github.com/user-attachments/assets/1cddeeb6-8986-4df5-926f-359ed5e3316e)

### Simplified Sequence Diagram 
![Untitled](https://github.com/user-attachments/assets/7084b6d7-d7a2-4b31-8381-9897b4f0bb0f)


