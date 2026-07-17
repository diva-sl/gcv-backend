# GCV | Enterprise Mailing & Communications Backend

A high-performance, secure Go backend API designed to handle communication pipelines, consultation requests, and newsletter subscriptions.

---

## 👨‍💻 Lead Developer
* **Developer Name:** Divakaran S
* **Role:** Lead Full-Stack Engineer / Architect
* **Contribution:** Built the complete Go backend infrastructure from scratch. Designed the secure concurrent SMTP mail dispatchers, route groups, and env configurations.

---

## 🛠️ Tech Stack & Architecture

* **Language:** Go / Golang (Compiled, concurrent, high-concurrency architecture)
* **Web Framework:** Gin-Gonic (High-performance HTTP routing)
* **Encryption:** `crypto/tls` (Ensures secure SSL/TLS handshakes over port 465)
* **Variables:** GoDotEnv (Seamless configuration management)

---

## 🚀 Key Features Implemented

1. **Sequential Mail Dispatcher:** Designed an automated helper function (`sendEmail`) that connects via TLS to Hostinger's SMTP relays to handle communication sequences.
2. **Double-Inquiry Notification:**
   * **Admin Warning:** Formulates a detailed HTML summary and sends it directly to the inbox of the administrators.
   * **Client Auto-Receipt:** Sends a personalized, styled HTML confirmation to the customer's email with a custom `Reply-To` header routing replies back to the sender.
3. **Newsletter Subscription Route (`/api/subscribe`):**
   * Processes email entries from the website's footer.
   * Sends real-time alert notifications directly to **`contact@gcvdanta.com`**.
4. **CORS Security:** Built-in middleware permitting secure cross-origin queries while avoiding header blocks during preflight pre-checks.

---

## 📦 Local Installation & Setup

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/diva-sl/gcv-backend.git
   cd gcv-backend
   ```

2. **Configure Environment Variables:**
   Create a `.env` file in the root folder:
   ```env
   PORT=5000
   SMTP_HOST=smtp.hostinger.com
   SMTP_PORT=465
   SMTP_USER=rajeshbandila@gcvdanta.com
   SMTP_PASS=Gcvdanta@123
   SMTP_FROM=connect@gcvdanta.com
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run Server Locally:**
   ```bash
   go run main.go
   ```

5. **Compile Binary:**
   ```bash
   go build -o main main.go
   ```

---

## 📄 License
Project developed for GCV Digital Product & Engineering. All rights reserved by Divakaran S.
