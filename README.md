# 🚗 Shimla Car Bazar - Dealership Management System (CDMS)

A robust, full-stack application designed to streamline the car buying and selling process. Built with a high-performance **Go (Golang)** backend and a responsive **Vanilla JS/Tailwind CSS** frontend, it features secure Role-Based Access Control (RBAC), multi-image cloud storage, and type-safe database operations.

---

## 🌟 Comprehensive Features

### 1. 🔐 Security & Authentication
* **Role-Based Access Control (RBAC):** Strict separation of workflows for `Admin` and `Customer`. API routes are protected using custom Go middlewares.
* **JWT Authentication:** Secure login sessions using JSON Web Tokens.
* **Password Hashing:** User passwords are encrypted using `Bcrypt` before being stored in the database.
* **Hardcoded Admin Guard:** Only a specific, authorized email (`admin@shimlacarbazar.com`) can register and access the system as an Admin.

### 2. 🖼️ Cloud Media Handling
* **Multi-Image Uploads:** Customers can upload multiple photos of their car (Front, Back, Interior) at once.
* **Cloudinary Integration:** Images are directly uploaded to **Cloudinary** via the Go backend, ensuring permanent, fast, and scalable media hosting.
* **Smart Storage:** Multiple image URLs are efficiently managed and stored as comma-separated strings in PostgreSQL.

### 3. 👥 User Workflows
**For Customers:**
* Browse the public marketplace of available cars without logging in.
* Submit a "Sell Car" proposal including car details, expected price, phone number, and multiple images.
* Access a private "My Requests" dashboard to track the live status (Pending/Approved) of their submissions.

**For Administrators:**
* Access a secured Admin Dashboard hidden from public view.
* Review incoming car sale proposals along with customer contact details and image galleries.
* **One-Click Approval:** Approving a proposal automatically transfers the car details and images into the main public inventory (`cars` table).

---

## 💻 Technology Stack

### Backend
* **Language:** Go (Golang)
* **Router:** `go-chi/chi/v5` (Fast, composable router)
* **Security:** `golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`
* **Media Storage:** `cloudinary-go/v2`

### Database
* **Database:** PostgreSQL (Run via Docker)
* **Query Builder:** `sqlc` (Generates fully type-safe Go code from SQL)
* **Migrations:** `golang-migrate`

### Frontend
* **Core:** HTML5, Vanilla JavaScript
* **Styling:** Tailwind CSS (via CDN)
* **Architecture:** Single Page Application (SPA) logic utilizing Fetch API with JWT Bearer tokens.

---

## 🗄️ Database Architecture

The system relies on a highly normalized PostgreSQL database:
1. **`users`**: Stores authentication details and roles (`admin` or `customer`).
2. **`cars`**: The main public inventory holding approved cars and their image galleries.
3. **`sell_proposals`**: Tracks customer requests to sell cars. Linked to `users` via foreign key. Includes contact numbers and multiple image URLs.
4. **`rc_transfers`**: Tracks the ownership transfer status between the dealership and buyers.

---

## 🌐 API Route Structure

| Method | Endpoint | Access Level | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/cars` | Public | View all available cars in inventory |
| `POST` | `/signup` | Public | Register a new user |
| `POST` | `/login` | Public | Authenticate and receive JWT |
| `POST` | `/upload` | Protected (Logged-in) | Upload multiple images to Cloudinary |
| `POST` | `/proposals` | Protected (Logged-in) | Submit a new car sale proposal |
| `GET` | `/proposals` | Protected (Smart) | Customers see their own; Admin sees all |
| `PUT` | `/proposals/{id}/approve`| Admin Only | Approve proposal & push to inventory |
| `POST` | `/cars` | Admin Only | Directly add a car to inventory |

---

## 🚀 Local Setup & Installation

Follow these instructions to run the Shimla Car Bazar project on your local machine.

### Prerequisites
* [Go](https://golang.org/doc/install) installed (v1.20+)
* [Docker](https://www.docker.com/) installed
* A free [Cloudinary](https://cloudinary.com/) account for image hosting.

### 1. Clone the Repository

git clone [https://github.com/Abhishekkapoor98/ShimlaCarBazar.git](https://github.com/Abhishekkapoor98/ShimlaCarBazar.git)
cd ShimlaCarBazar


Start a PostgreSQL container:

docker run --name car_bazar_db -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secretpassword -d postgres:alpine

Create the database and apply the schema:

docker exec -it car_bazar_db createdb --username=root --owner=root car_bazar
# Note: Apply your schema definitions from internal/db/migration/000001_init_schema.up.sql here

You must provide your Cloudinary API URL for the image uploads to work. Run this in your terminal:
export CLOUDINARY_URL=cloudinary://<API_KEY>:<API_SECRET>@<CLOUD_NAME>

Run the Server

go mod tidy
go run cmd/api/main.go

The application will be live at: http://localhost:8080

🧪 Testing the Application
To properly test the RBAC (Role-Based Access Control) system:

Admin Account:

Click 'Login/Signup'

Email: admin@shimlacarbazar.com (This is strictly validated in the backend)

Role: Admin

Customer Account:

Sign up with any other email (e.g., test@gmail.com)

Role: Customer

Log in with a customer account to submit a proposal with images, then log out and log in with the Admin account to approve it!

Built with ❤️ by Abhishek Kapoor.
