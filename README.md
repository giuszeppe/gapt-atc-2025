# ATC Simulation Platform

This project is an Air Traffic Control (ATC) simulation platform designed for training and testing ATC scenarios. It features a backend API (Go) and a frontend (Vue).

## Features

- **Simulations**: Multiple ATC scenarios (takeoff, enroute, landing) with realistic transcripts.
- **User Roles**: Tower and aircraft roles for interactive simulation.
- **Authentication**: Secure login for users.
- **WebSocket/HTTP API**: Real-time and RESTful communication.

## Project Structure

- `backend/`: Go API server, database seeds, and configuration.
- `frontend/`: Vue.js web client.

## Prerequisites

- Go (1.20+)
- Node.js (21.1.0+) & npm(10.4.0+)
- SQLite3

## Setup

### 1. Backend

```sh
cd backend
cp .env.example .env   # Edit as needed
go mod tidy
```

### 2. Frontend

```sh
cd frontend
npm install
```

### 3. Starting the project
```shell
cd <project_root>
make run-refresh #first time run
make run #subsequent runs
```

## Usage
### Logging in
- Two users: `admin:password` and `test:password` for testing.

### Using the web client
- Access the web client at `http://localhost:5173` (or as configured).
- If needed, access the user guide from the homepage.
- Start or join a simulation as tower or aircraft.
