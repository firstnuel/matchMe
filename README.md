# Match-Me 💘

Match-Me is a modern, full-stack dating application designed to connect people through location-based matching and real-time interaction. It features a sleek React/TypeScript frontend and a powerful Go backend, providing a seamless and responsive user experience.

\*\*

-----

## 🚀 Getting Started

Follow these instructions to get a local copy up and running for development and testing.

### Prerequisites

Ensure you have the following installed on your system:

  * **Node.js**: `v18` or later
  * **Go**: `v1.24` or later
  * **PostgreSQL**: with the **PostGIS** extension enabled

### Installation & Setup

1.  **Clone the Repository**

    ```bash
    git clone https://gitea.kood.tech/ravikantpandit/match-me.git
    cd match-me
    ```

2.  **Set Up the Database**
    This project uses PostgreSQL with the PostGIS extension for location services.

    ```bash
    # On macOS (using Homebrew)
    brew install postgresql postgis

    # On Debian/Ubuntu
    sudo apt-get update
    sudo apt-get install postgresql postgresql-contrib postgis

    # Create the database and enable the extension
    createdb matchme
    psql -d matchme -c "CREATE EXTENSION postgis;"
    ```

3.  **Configure Environment Variables**
    Create a `.env` file inside the `server/` directory and populate it with your configuration.

    ```bash
    # server/.env
    APP_ENV=development
    PORT=8080
    DATABASE_URL=postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/matchme?
    DATABASE_NAME=postgres
    sslmode=disable
    JWT_SECRET=a-very-strong-and-secret-key
    CLOUDINARY_URL=your-cloudinary-api-environment-variable
    ```

      Create a `.env` file inside the `client/` directory and add this line.

      ```bash
      VITE_API_BASE_URL=http://localhost:8080
      ``` 

      > Remember that the port set in both the client and server must match, for example if you change the `PORT` in `server/.env` to 3000 then the client should be http://localhost:3000


4.  **Install Dependencies and Run**
    The `Makefile` contains all the necessary commands to install dependencies and run the application.

    ```bash
    # Install dependencies for both client and server
    make dev-setup
    ```

    Next, run the frontend and backend servers in **two separate terminal windows**:

    ```bash
    # In terminal 1: Start the client dev server
    
    make dev-client
    # Frontend will be available at http://localhost:5173
    ```

    ```bash
    # In terminal 2: Start the backend server
    make run-server
    # Backend will be running at http://localhost:8080
    ```
   
5.  if you do not have make installed
    
      ```bash
      # In terminal 1: Install the client and start the frontend dev server
      
      cd client 
      npm i && npm run dev
      # Frontend will be available at http://localhost:5173
      ```

      ```bash
      # In terminal 2: Start the backend server
      cd server && go run ./cmd/server
      # Backend will be running at http://localhost:8080
      ```
   

-----
### Database Management

The backend includes helpful commands for managing the database during development or testing.

```bash
# Navigate to the server directory
cd server

# Display all available database commands
go run ./cmd/server -h

# Seed the database with 50 test users
go run ./cmd/server -p 50

# Completely reset the database (drop all data)
go run ./cmd/server -r

# Reset the database and then seed it with 25 test users
go run ./cmd/server -rp 25
```

## ⚙️ Usage

### Makefile Commands

A `Makefile` at the root of the project simplifies common tasks.

| Command | Description |
| :--- | :--- |
| `make help` | Displays a list of all available commands. |
| `make dev-setup` | Installs all dependencies for both client and server. |
| `make dev-client` | Starts the frontend development server with hot-reloading. |
| `make run-server`| Starts the backend API server. |
| `make build-all` | Creates production-ready builds for both client and server. |
| `make clean` | Removes build artifacts and `node_modules`. |

-----

## ✨ Features

  * **✅ Secure User Authentication**: JWT-based authentication for secure sessions and profile management.
  * **📍 Geospatial Matching**: Utilizes PostGIS to discover and connect with potential matches nearby.
  * **💬 Real-Time Chat**: Instant messaging between connected users, powered by WebSockets for a fluid conversation experience.
  * **📸 Cloud-Based Image Handling**: Efficient and secure photo uploads and storage managed via Cloudinary.
  * **🤝 Connection Management**: A complete system to send, accept, and manage connection requests.
  * **📱 Fully Responsive Design**: A beautiful and intuitive interface that works flawlessly on both desktop and mobile devices.

-----

## 🛠️ Tech Stack

The project is built with a modern and robust technology stack, separating concerns between a client-side application and a server-side API.

| **Component** | **Technology** | **Purpose** |
| :--- | :--- | :--- |
| **Frontend** | React 19 (TypeScript) | UI development |
| | Vite | Build tooling & dev server |
| | Zustand & React Query | State management & server-state synchronization |
| | Bootstrap | Styling and responsive layout |
| | WebSocket API | Real-time communication |
| **Backend** | Go 1.24+ | Core application logic |
| | Gin | High-performance HTTP web framework |
| | Ent | ORM for type-safe database access |
| **Database** | PostgreSQL + PostGIS | Relational data and geospatial queries |
| **Infrastructure** | JWT | Authentication |
| | Cloudinary | Cloud-based image storage |


-----

## 📁 Project Structure

The repository is organized into two main parts: a `client` directory for the frontend and a `server` directory for the backend.

```
match-me/
├── client/           # React TypeScript frontend application
│   ├── src/
│   │   ├── features/ # Feature-based modules (e.g., auth, chat)
│   │   └── shared/   # Reusable components, hooks, and utilities
│   └── package.json
├── server/           # Go backend application
│   ├── cmd/server/   # Main application entrypoint and flags
│   ├── api/          # HTTP routes, handlers, and middleware
│   ├── ent/          # Auto-generated ORM code, models, and migrations
│   └── internal/     # Core business logic and services
└── Makefile          # Commands for building and running the project
```

