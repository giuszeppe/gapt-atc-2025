# ATC Comm Simulator backend

## Setup
1. Clone the repository:
   ```bash
   git clone
   cd gapt-atc-2025
    ```
2. Populate the env-file
    ```bash
    cp .env.example .env # replace with your own values
    ```
3. Create a database:
    ```bash
   go build -o backend
   ./backend database refresh
    ```
4. Run the server:
    ```bask
   ./backend server
   ```