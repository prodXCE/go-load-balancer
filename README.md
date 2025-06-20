# Go HTTP Load Balancer

A simple, robust, and highly concurrent HTTP load balancer built from scratch in Go. This project was developed as a step-by-step guide to learning core Go concepts (structs, methods, interfaces, concurrency) while simultaneously understanding fundamental principles of system design (load balancing, high availability, health checks).

The application listens for incoming HTTP requests and distributes them among a pool of backend servers using a Round Robin algorithm. It includes a concurrent health checker to automatically detect and handle server failures, ensuring traffic is only sent to healthy backends.

---

## Features

- **HTTP Reverse Proxy**: Sits in front of backend servers and forwards client requests.
- **Round Robin Load Balancing**: Distributes incoming requests sequentially across the available servers.
- **Concurrent Health Checks**: Periodically checks the health of all backend servers in the background without blocking request handling.
- **High Availability**: Automatically removes unhealthy servers from the rotation and adds them back once they recover.
- **Concurrency-Safe**: Uses Go's mutexes to safely handle shared state across multiple goroutines, preventing race conditions.
- **Simple & Extensible**: Built with the standard library and minimal code, making it easy to understand and extend with new features.

---

## Getting Started

Follow these instructions to get the load balancer and its test backend servers running on your local machine.

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.18 or higher) installed on your system.
- Git for version control.
- Two separate terminal windows.

### Installation & Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/prodXCE/go-load-balancer.git
    cd go-load-balancer
    ```

2.  **Run the Backend Servers:**
    In your **first terminal window**, start the dummy backend servers. These will listen on ports `8081`, `8082`, and `8083`.
    ```bash
    go run backend.go
    ```
    You will see logs indicating the servers have started.

3.  **Run the Load Balancer:**
    In your **second terminal window**, start the load balancer. This will listen on port `8080`.
    ```bash
    go run main.go
    ```
    You will see logs indicating the load balancer is running and has configured the backend servers.

4.  **Test the Load Balancer:**
    Open your web browser and navigate to `http://localhost:8080`.

    Each time you refresh the page, you should see a message from a different backend server (e.g., "Hello from backend server on port: 8081"). This confirms the Round Robin distribution is working.

5.  **Test the Health Check:**
    - Go to the **first terminal** (running `backend.go`) and stop the process with `Ctrl + C`.
    - Wait about 10 seconds. In the **second terminal**, you will see logs like "Server ... is down."
    - Refresh your browser at `http://localhost:8080`. You will now receive a "Service not available" error.
    - Restart the backend servers in the first terminal (`go run backend.go`).
    - Wait another 10 seconds. The load balancer will log that the servers are "back up."
    - Refresh your browser, and traffic will flow again.

---

## How It Works

The system is composed of two simple applications:

-   **`main.go` (The Load Balancer)**: This is the core application. It has several key components:
    -   **Server & ServerPool Structs**: Custom data types that model our backend servers and the pool that manages them.
    -   **Round Robin Logic**: The `GetNextPeer()` method safely selects the next healthy server from the pool.
    -   **Reverse Proxy**: Uses `net/http/httputil.NewSingleHostReverseProxy` to forward requests to the chosen backend.
    -   **Health Checker**: A background goroutine runs every 10 seconds, concurrently checking the status of all servers and updating their `Alive` status in a thread-safe manner.

-   **`backend.go` (The Dummy Servers)**: A simple application that launches three independent web servers on different ports. Each server simply responds with a message indicating which port it's on. This allows us to easily simulate a real-world server farm for testing purposes.

---

## Core Concepts Illustrated

This project provides a practical demonstration of several key concepts:

### Go Language Concepts
- **Packages & Modules**: Project setup with `go mod`.
- **Structs & Methods**: Creating our own data types (`Server`, `ServerPool`).
- **Slices**: Managing the list of backend servers.
- **Concurrency**:
  - **Goroutines**: For running the health checker in the background.
  - **Mutexes (`sync.Mutex`, `sync.RWMutex`)**: For protecting shared data (`current` index, `Alive` status) from race conditions.
- **Standard Library**: Extensive use of `net/http`, `net/url`, `log`, `sync`, and `time`.
- **Error Handling**: The standard `if err != nil` pattern.
- **`defer` Keyword**: For safely unlocking mutexes.

### System Design Concepts
- **Load Balancing**: The core principle of distributing traffic.
- **Reverse Proxies**: The mechanism used to forward traffic.
- **High Availability**: How health checks and server pools keep the service online despite failures.
- **Scalability**: The design allows for easily adding more servers to the pool.
- **Decoupling**: The load balancer and backend servers are separate applications that can be run and managed independently.


