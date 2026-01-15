## Soal Tes Teknis Developer: Aplikasi Manajemen Stok ATK (Full-Stack) ##
                            # Rois Afif Anggoro #
# 📦 Aplikasi Manajemen Inventori ATK - Full Stack

Sistem manajemen inventori Alat Tulis Kantor (ATK) dengan fitur CRUD lengkap.

## 🎯 Overview

Aplikasi full-stack untuk mengelola inventori ATK di perusahaan, mencakup:
- **Backend:** REST API dengan Go/Golang
- **Frontend:** React dengan Vite
- **Database:** MySQL
- **Env:** Alamat Server pribadi
- **Table:** mst_atk

Sesuai dengan requirements dari dokumen tes teknis "Soal Tes Teknis Developer IT Corp".

---

## 🏗️ Arsitektur

```
Manajamen ATK/
├── ATK-Backend/           # Go REST API Server
│   ├── main.go           # Entry point & DB connection
│   ├── models/           # Data models & DB operations
│   │   └── atk.go       # ATK model with CRUD
│   ├── routes/           # HTTP handlers
│   │   ├── index.go     # GET & POST endpoints
│   │   └── id.go        # PUT & DELETE endpoints
│   ├── go.mod           # Dependencies
│   └── .env             # Database configuration
│
├── ATK-Frontend/          # React Dashboard
│   ├── src/
│   │   ├── components/
│   │   │   └── ATKDashboard.jsx   # Main CRUD component
│   │   ├── App.jsx       # Root component
│   │   ├── main.jsx      # Entry point
│   │   └── index.css     # Tailwind styles
│   ├── package.json      # Dependencies
│   ├── vite.config.js    # Vite & proxy config
│   └── README.md         # Frontend docs
│
└── TESTING.md            # Testing guide
```

---

## 🚀 Quick Start

### 1. Setup Backend

```bash
# Masuk ke folder backend
cd ATK-Backend

# Pastikan .env sudah dikonfigurasi
# DATABASE=user:password@tcp(host:port)/database
# PORT=5200

# Jalankan server
go run .
```

**Expected output:**
```
Connected to DB and initialized table mst_atk
Server started on :5200
```

### 2. Setup Frontend

```bash
# Buka terminal baru
cd ATK-Frontend

# Install dependencies (first time only)
npm install

# Jalankan dev server
npm run dev
```

**Expected output:**
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### 3. Akses Aplikasi

Buka browser (dengan klik pada terminal): **http://localhost:5173**

---

## 💾 Database Schema

### Tabel: `mst_atk`

| Column | Type | Constraint | Deskripsi |
|--------|------|------------|-----------|
| `id` | INT | PRIMARY KEY, AUTO_INCREMENT | ID unik item |
| `nama` | VARCHAR(255) | NOT NULL | Nama ATK |
| `jenis` | VARCHAR(255) | NOT NULL | Kategori/jenis ATK |
| `qty` | INT | NOT NULL | Jumlah stok |
| `is_deleted` | BOOLEAN | NULL | Soft delete | (opsional/improvements)

**Contoh insert data:**
```sql
INSERT INTO mst_atk (nama, jenis, qty) VALUES
('Pensil 2B', 'Alat Tulis', 50),
('Kertas HVS A4', 'Kertas', 500),
('Stapler', 'Alat Kantor', 15);
```

---

## 🔌 API Endpoints

Base URL: `http://localhost:5200` (I recomend to use postman to backend testing)

### GET /api/atk (http://localhost:5200/api/atk)
Ambil semua data ATK

**Response:**
```json
[
  {
    "id": 1,
    "nama": "Pensil 2B",
    "jenis": "Alat Tulis",
    "qty": 50
  }
]
```

### POST /api/atk (http://localhost:5200/api/atk)
Tambah ATK baru

**Request Body:**
```json
{
  "nama": "Pulpen Biru",
  "jenis": "Alat Tulis",
  "qty": 100
}
```

**Response:**
```json
{
  "id": 2,
  "nama": "Pulpen Biru",
  "jenis": "Alat Tulis",
  "qty": 100
}
```

### PUT /api/atk/{id} (http://localhost:5200/api/atk{id})
Update ATK by ID

**Request Body:**
```json
{
  "nama": "Pulpen Biru Tebal",
  "jenis": "Alat Tulis",
  "qty": 80
}
```

**Response:**
```json
{
  "id": 2,
  "nama": "Pulpen Biru Tebal",
  "jenis": "Alat Tulis",
  "qty": 80
}
```

### DELETE /api/atk/{id} (http://localhost:5200/api/atk{id})
Hapus ATK by ID

eror responses:
**Response:** `204 No Content`
**Response:** `404 no api`
**Response:** `400 wrong api`
**Response:** `500 wrong value`

---

## ✨ Fitur Frontend

### 📊 Dashboard Statistics
- **Total Item ATK** - Jumlah jenis item berbeda
- **Total Stok** - Total kuantitas semua item
- **Jenis ATK** - Jumlah kategori berbeda

### 📋 Tabel Inventori
- Tampilan data lengkap dalam tabel
- Color-coded stock levels:
  - 🔴 Merah: qty < 5 (stok kritis)
  - 🟡 Kuning: qty 5-19 (stok rendah)
  - 🟢 Hijau: qty >= 20 (stok aman)

### ➕ Create/Edit Modal
- Form validation (required fields, qty >= 0)
- Shared modal untuk Create & Edit
- Real-time data update

### 🗑️ Delete Confirmation
- Konfirmasi sebelum hapus (swal fire)
- Immediate UI update

### 🔄 Real-time Features
- Auto-refresh setelah operasi CRUD
- Loading states
- Error handling dengan pesan jelas

---

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.25.5
- **Router:** gorilla/mux v1.8.0
- **Database:** MySQL (go-sql-driver/mysql v1.6.0)
- **Env:** godotenv v1.5.1

### Frontend
- **Framework:** React 18
- **Build Tool:** Vite 5
- **Styling:** Tailwind CSS 3
- **HTTP Client:** Axios 1.10.0
- **Icons:** Lucide React
- **Router:** React Router DOM 7

---