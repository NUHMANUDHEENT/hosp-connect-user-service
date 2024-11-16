# User Service

The **User Service** is a core part of the HospConnect system, responsible for managing users, doctors, and admins while enabling seamless communication and scheduling functionalities.

---

## **Features**

### **User Management**
- Handles patient registration, profile updates, and authentication.
- Secure role-based access control (Admin, Doctor, Patient).

### **Doctor & Admin Management**
- Admin can manage doctors' profiles, availability, and roles.
- Doctors can sync availability with **Google Calendar** for streamlined scheduling.

### **Data Storage**
- **MongoDB**: For user profiles and role management data.
- **Redis**: For caching frequently accessed data to improve performance.

### **Communication Features**
- **Video Call Integration**: Powered by **Jitsi**, enabling virtual consultations.
- **Real-Time Chat**: Patients and admins can communicate using WebSocket.

### **Inter-Service Communication**
- gRPC endpoints for efficient communication with other services.

---

## **Technology Stack**
- **Backend:** Go (Golang)
- **Database:** MongoDB, Redis
- **Communication:** gRPC, WebSocket
- **External APIs:** Google Calendar API, Jitsi Meet

---

## **How to Run**

### Clone the Repository
```bash
git clone https://github.com/your-username/user-service.git
cd user-service
