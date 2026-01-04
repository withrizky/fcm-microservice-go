
# High-Performance FCM Microservice (Go)

A high-throughput Push Notification service built with **Golang** and **Firebase Cloud Messaging (FCM)**.

This service utilizes a **Worker Pool pattern** to handle thousands of concurrent notification requests efficiently. It separates the API layer from the heavy lifting of communicating with Google's servers, ensuring low latency for clients.

## Architecture

The system uses an **In-Memory Event Driven** architecture. HTTP requests are validated and pushed into a buffered channel. Background workers pick up jobs and execute the FCM delivery using a **Reusable Firebase Connection** (Singleton pattern).
## Key Features

* Reusable Connection: Initializes the Firebase App only once at startup, preventing connection overhead and rate limits.

* Dual Mode Support: Supports sending messages to a specific device (Token) or broadcasting to a group (Topic) via a single endpoint.

* In-Memory Worker Pool: Processes notifications in parallel using Go routines (Default: 50 concurrent workers).

* Graceful Shutdown: Ensures pending notifications are sent before the server stops.

* Pure Go: No external message brokers (RabbitMQ/Redis) required.

## 📂 Folder Structure

```
fcm_microservice/
├── cmd/
│   └── server/
│       └── main.go           # Application Entry Point
├── internal/
│   ├── model/            # Data Structures (Payloads)
│   ├── fcm/              # Firebase Client Wrapper
│   └── worker/           # Worker Pool Logic
├── signal-app-xxxx.json  # Firebase Service Account Key
├── .env                  # Configuration File
├── go.mod                # Go Modules
└── README.md             # Documentation
```

## 🛠 Prerequisites

* **Go** (version 1.18+)
* **Firebase Service Account Key** : A JSON file downloaded from the Firebase Console.

* *    Project Settings > Service Accounts > Generate New Private Key.



## 🚀 Installation & Setup

1. **Clone the repository**
```bash
git clone https://github.com/withrizky/fcm-microservice-go.git
cd fcm-microservice-go

```


2. **Install Dependencies**
```bash
go mod tidy

```


3. **Environment Configuration**
Create a `.env` file in the root directory:
```env
PORT=8083
# Make sure this filename matches your actual JSON file
GOOGLE_APPLICATION_CREDENTIALS=signal-app-xxxxx-firebase.json
```


4. **Run the Server**
```bash
go run cmd/server/main.go

```



## 📡 API Documentation

### Send Message

Sends a message to the processing queue.

* **URL**: `/send-fcm`
* **Method**: `POST`
* **Content-Type**: `application/json`

**Request Body Spesific Device:**

```json
{
    "title": "Order Update",
    "body": "Your package has been shipped!",
    "to": "device_token_xyz123...",
    "data": {
        "order_id": "ORD-555"
    }
}

```

**Request Body Send Topic**

```json
{
    "title": "Breaking News",
    "body": "New version available.",
    "topic": "news_updates"
}

```

**Response (Success):**

```json
{
    "status": "queued",
    "message": "Pesan masuk antrean"
}

```

* Status Code: `202 Accepted`

**Response (Error):**

```json
{
    "error": "Payload invalid"
}

```

* Status Code: `400 Bad Request`

## 📈 Performance Strategy

This service is optimized for high concurrency:

1. **Buffered Channels**: Can hold up to **5,000** (configurable) pending messages in RAM.
2. **Concurrency**: Spawns **50** (configurable) concurrent workers. This means 50 messages are processed in parallel every millisecond.
3. **Singleton Client**: Reuses the HTTP/2 connection to Firebase for maximum throughput.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
