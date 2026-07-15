# Data Acquisition Engine

Sistem intelijen data perusahaan yang dibangun dengan **Go (Golang)** untuk mengumpulkan dan mengagregasi informasi dari tiga sumber utama — **Metadata Website**, **Informasi Domain (RDAP)**, dan **Lokasi Geografis (OpenStreetMap Nominatim)** — secara cepat dan efisien melalui pemrosesan paralel dengan goroutine.

Proyek ini merupakan bagian dari **Technical Challenge Program Praktik Kerja Lapangan (PKL)** di **PT Berani Digital Indonesia**.

---

## Desain & Arsitektur Software

Proyek ini menerapkan prinsip **Clean Architecture** dengan pemisahan tanggung jawab yang jelas antar lapisan:

```
data-acquisition-engine/
├── cmd/api/              # Entry point aplikasi (main.go)
├── internal/
│   ├── handler/          # HTTP Handler (Controller layer)
│   ├── service/          # Business logic & integrasi API (Service layer)
│   └── response/         # Helper standarisasi response JSON
```

- **Cmd Layer** (`cmd/api/main.go`) — Inisialisasi aplikasi, dependency injection, dan registrasi routing.
- **Handler Layer** (`internal/handler/`) — Menerima request HTTP, memvalidasi input, memanggil service, dan mengembalikan response.
- **Service Layer** (`internal/service/`) — Seluruh logika bisnis dan komunikasi dengan API eksternal.

### Optimasi Concurrency

Endpoint integrasi final `GET /company-information` dirancang untuk meminimalkan latency dengan memanfaatkan **goroutine** dan **sync.WaitGroup**:

1. **Fase Paralel Awal** — Pemanggilan `WebsiteService.Extract()` dan `DomainService.Extract()` dieksekusi secara bersamaan dalam dua goroutine terpisah. Kedua operasi ini tidak memiliki dependensi satu sama lain, sehingga dapat berjalan paralel penuh.

2. **Fase Paralel Bergantung** — Setelah hasil website tersedia, `LocationService.Find()` dijalankan secara paralel menggunakan properti `title` dari metadata website sebagai query pencarian lokasi di OpenStreetMap Nominatim.

3. **Mekanisme Fallback** — Jika pencarian lokasi menggunakan `title` website menghasilkan `null` atau error (misalnya karena title terlalu panjang atau tidak relevan), sistem secara otomatis mengekstrak nama domain sebelum ekstensi (contoh: `"paper"` dari `"paper.id"`) dan melakukan pemanggilan ulang ke LocationService dengan query tersebut.

```go
// Alur orkestrasi:
//
// ┌─────────────────┐     ┌─────────────────┐
// │ WebsiteService  │     │  DomainService  │    ← Paralel (wg.Add(2))
// └────────┬────────┘     └─────────────────┘
//          │ title
//          ▼
// ┌─────────────────┐
// │ LocationService │                           ← Setelah title tersedia
// └────────┬────────┘
//          │ (jika null/error)
//          ▼
// ┌─────────────────┐
// │  Fallback:      │
// │  extractNama()  │                           ← Retry dengan nama domain
// │  → LocationSvc  │
// └─────────────────┘
```

---

## Instalasi & Menjalankan Aplikasi

### Prasyarat

- **Go** versi 1.21 atau lebih baru
- Koneksi internet (untuk memanggil API eksternal)

### Langkah Instalasi

```bash
# 1. Clone repository
git clone https://github.com/<username>/data-acquisition-engine.git
cd data-acquisition-engine

# 2. Download dependencies
go mod download

# 3. Jalankan aplikasi
go run cmd/api/main.go
```

Server akan berjalan pada `http://localhost:8080`.

---

## Dokumentasi Endpoint API

### Standarisasi Response

Seluruh endpoint mengembalikan response JSON dengan format konsisten:

```json
{
  "success": true,
  "data": { ... }
}
```

```json
{
  "success": false,
  "error": "pesan error"
}
```

---

### 1. Extract Website Metadata

Mengekstrak metadata dari sebuah URL website, meliputi title, deskripsi, Open Graph, email, nomor telepon, dan tautan media sosial.

- **Endpoint:** `POST /extract/website`
- **Content-Type:** `application/json`

**Request Body:**

```json
{
  "url": "https://paper.id"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "url": "https://paper.id",
    "title": "Paper.id - Platform Invoice & Pembayaran Bisnis",
    "description": "...",
    "canonical": "https://paper.id",
    "favicon": "/favicon.ico",
    "emails": ["support@paper.id"],
    "phones": ["+6281234567890"],
    "social_media": ["https://linkedin.com/company/paperid"],
    "open_graph": {
      "title": "Paper.id",
      "description": "...",
      "image": "https://paper.id/og-image.png"
    }
  }
}
```

---

### 2. Extract Domain Intelligence

Mengambil informasi domain melalui protokol RDAP, termipu registrar, tanggal registrasi, tanggal kedaluwarsa, status, dan nameserver.

- **Endpoint:** `POST /extract/domain`
- **Content-Type:** `application/json`

**Request Body:**

```json
{
  "domain": "paper.id"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "domain": "paper.id",
    "registrar": "Digital Registra",
    "registered_at": "2018-03-15 00:00:00",
    "expired_at": "2026-03-15 00:00:00",
    "last_updated": "2024-06-01 00:00:00",
    "status": ["serverTransferProhibited"],
    "nameservers": ["ns1.example.com", "ns2.example.com"]
  }
}
```

---

### 3. Find Location

Mencari informasi lokasi geografis berdasarkan query teks menggunakan OpenStreetMap Nominatim.

- **Endpoint:** `POST /extract/location`
- **Content-Type:** `application/json`

**Request Body:**

```json
{
  "query": "Paper.id Jakarta"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "display_name": "Jakarta, Indonesia",
    "latitude": "-6.2088",
    "longitude": "106.8456",
    "importance": 0.95,
    "osm_type": "relation",
    "address": {
      "city": "Jakarta",
      "country": "Indonesia"
    }
  }
}
```

---

### 4. Company Information (Integrated Endpoint)

Endpoint integrasi final yang menggabungkan ketiga service secara paralel. Menerima satu parameter query `domain` dan mengembalikan seluruh data intelijen perusahaan dalam satu response.

- **Endpoint:** `GET /company-information`
- **Method:** `GET`
- **Query Parameter:** `domain` (wajib)

**Request:**

```
GET /company-information?domain=paper.id
```

**Response:**

```json
{
  "success": true,
  "data": {
    "website": {
      "url": "https://paper.id",
      "title": "Paper.id - Platform Invoice & Pembayaran Bisnis",
      "description": "...",
      "canonical": "https://paper.id",
      "favicon": "/favicon.ico",
      "emails": ["support@paper.id"],
      "phones": ["+6281234567890"],
      "social_media": ["https://linkedin.com/company/paperid"],
      "open_graph": {
        "title": "Paper.id",
        "description": "...",
        "image": "https://paper.id/og-image.png"
      }
    },
    "domain": {
      "domain": "paper.id",
      "registrar": "Digital Registra",
      "registered_at": "2018-03-15 00:00:00",
      "expired_at": "2026-03-15 00:00:00",
      "last_updated": "2024-06-01 00:00:00",
      "status": ["serverTransferProhibited"],
      "nameservers": ["ns1.example.com", "ns2.example.com"]
    },
    "location": {
      "display_name": "Jakarta, Indonesia",
      "latitude": "-6.2088",
      "longitude": "106.8456",
      "importance": 0.95,
      "osm_type": "relation",
      "address": {
        "city": "Jakarta",
        "country": "Indonesia"
      }
    }
  }
}
```

**Error Handling:**

| Kondisi | HTTP Status | Response |
|---------|-------------|----------|
| Parameter `domain` kosong | `400 Bad Request` | `{"success": false, "error": "domain query param is required"}` |
| Seluruh service gagal | `500 Internal Server Error` | `{"success": false, "error": "all services failed to return data"}` |
| Sebagian service gagal | `200 OK` | Field yang gagal bernilai `null` pada response |

---

## Video Presentasi

[Video Presentasi (YouTube Unlisted)](https://youtu.be/<placeholder-link>)
