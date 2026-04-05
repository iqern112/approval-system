# Approval System

ระบบจัดการคำขออนุมัติ โดยแบ่งผู้ใช้งานออกเป็น 2 บทบาท คือ **Admin** และ **User**

---

## โครงสร้างโปรเจค

```
approval-system/
├── approval-back/       # Backend API (Go)
├── approval-front/      # Frontend (Angular)
├── init/                # SQL สำหรับตั้งค่าฐานข้อมูลเริ่มต้น
└── docker-compose.yml   # Docker สำหรับรัน PostgreSQL ในเครื่อง
```

---

## ฟีเจอร์หลัก

### User
- ล็อกอินเข้าสู่ระบบด้วยบัญชีที่ Admin สร้างให้
- สร้างรายการคำขออนุมัติ (Approval Request)
- ติดตามสถานะคำขอของตัวเอง

### Admin
- ล็อกอินเข้าสู่ระบบ
- เพิ่ม / จัดการบัญชีผู้ใช้ (User)
- ดูรายการคำขออนุมัติทั้งหมด
- **อนุมัติ** หรือ **ไม่อนุมัติ** คำขอ พร้อมระบุเหตุผล

---

## Tech Stack

| ส่วน | เทคโนโลยี |
|---|---|
| Backend | Go |
| Frontend | Angular |
| Database | PostgreSQL 15 |
| Container | Docker / Docker Compose |

---

## วิธีรันในเครื่อง (Local Development)

### สิ่งที่ต้องมี
- [Docker](https://www.docker.com/) และ Docker Compose
- [Go](https://go.dev/) 1.20+
- [Node.js](https://nodejs.org/) 18+ และ Angular CLI

### 1. รัน Database

```bash
docker-compose up -d
```

PostgreSQL จะรันที่ `localhost:5433`
ข้อมูลเริ่มต้นจะถูก seed อัตโนมัติจากไฟล์ใน `init/`

### 2. รัน Backend

```bash
cd approval-back
go mod tidy
go run main.go
```

Backend จะรันที่ `http://localhost:8080` (หรือ port ที่กำหนดใน .env)

### 3. รัน Frontend

```bash
cd approval-front
npm install
ng serve
```

Frontend จะรันที่ `http://localhost:4200`

---

## ⚙️ Environment Variables (Backend)

สร้างไฟล์ `.env` ใน `approval-back/` โดยอ้างอิงจากตัวอย่างนี้:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=approval_db
```

---

## Database

| Config | ค่า |
|---|---|
| Host | localhost |
| Port | 5433 |
| User | myuser |
| Password | mypassword |
| Database | approval_db |

---

## 📸 UI Screenshots

### หน้า Login
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125458" src="https://github.com/user-attachments/assets/e104f41c-0727-4293-8d87-a761db5275cd" />


### หน้า Dashboard (User)
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125726" src="https://github.com/user-attachments/assets/edc009d2-0a4a-41e8-badd-b6a970cfbac6" />


### หน้าต่างสร้างคำขอใหม่
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125733" src="https://github.com/user-attachments/assets/1f4dcd43-ed34-42df-9b4b-70a29a811de3" />


### หน้า Admin — รายการคำขอทั้งหมด
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125618" src="https://github.com/user-attachments/assets/e09a81ca-79b1-4713-9f2d-cc3aa102a8b9" />


### หน้าสร้างคำขออนุมัติ
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125639" src="https://github.com/user-attachments/assets/02042900-9770-4ad2-a498-6515b69cea50" />


### หน้า Admin — อนุมัติ / ไม่อนุมัติ
<img width="1920" height="1032" alt="Screenshot 2026-04-05 125658" src="https://github.com/user-attachments/assets/e4ad8a92-7ec2-4bb3-ab4c-5a77b27d16e9" />

---
