# Approval System

ระบบจัดการคำขออนุมัติ โดยแบ่งผู้ใช้งานออกเป็น 2 บทบาท คือ **Admin** และ **User**

---

## 🗂️ โครงสร้างโปรเจค

```
approval-system/
├── approval-back/       # Backend API (Go)
├── approval-front/      # Frontend (Angular)
├── init/                # SQL สำหรับตั้งค่าฐานข้อมูลเริ่มต้น
└── docker-compose.yml   # Docker สำหรับรัน PostgreSQL ในเครื่อง
```

---

## ✨ ฟีเจอร์หลัก

### 👤 User
- ล็อกอินเข้าสู่ระบบด้วยบัญชีที่ Admin สร้างให้
- สร้างรายการคำขออนุมัติ (Approval Request)
- ติดตามสถานะคำขอของตัวเอง

### 🛡️ Admin
- ล็อกอินเข้าสู่ระบบ
- เพิ่ม / จัดการบัญชีผู้ใช้ (User)
- ดูรายการคำขออนุมัติทั้งหมด
- **อนุมัติ** หรือ **ไม่อนุมัติ** คำขอ พร้อมระบุเหตุผล

---

## 🛠️ Tech Stack

| ส่วน | เทคโนโลยี |
|---|---|
| Backend | Go |
| Frontend | Angular |
| Database | PostgreSQL 15 |
| Container | Docker / Docker Compose |

---

## 🚀 วิธีรันในเครื่อง (Local Development)

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

## 🗄️ Database

| Config | ค่า |
|---|---|
| Host | localhost |
| Port | 5433 |
| User | myuser |
| Password | mypassword |
| Database | approval_db |

---

## 📦 Deploy

| ส่วน | Platform |
|---|---|
| Frontend (Angular) | [Vercel](https://vercel.com) |
| Backend (Go) | [Render](https://render.com) |
| Database (PostgreSQL) | [Render](https://render.com) |
