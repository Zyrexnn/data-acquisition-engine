# Product Requirement Document (PRD)
## Project: Data Acquisition Engine (Business Intelligence & Verification Core)

---

## 1. Project Overview
Aplikasi ini adalah sebuah *Core Data Acquisition Engine* berskala industri yang berfungsi untuk mengumpulkan, memproses, memetakan (*mapping*), dan mengintegrasikan data dari berbagai sumber eksternal menjadi informasi terstruktur berbasis JSON. Sistem ini dibangun menggunakan arsitektur modular (*Service-Connector Pattern*) untuk memastikan setiap komponen dapat diuji secara independen, mudah dirawat, dan memiliki resiliensi tinggi terhadap kegagalan jaringan eksternal.

---

## 2. Technical Stack & Architecture
- **Language & Framework:** Go (Golang) v1.22+ dengan **Fiber** atau **Gin** sebagai HTTP Router.
- **Architecture Pattern:** Monolith Modular / Service-Connector Pattern.
- **Key Features:** Concurrency (Parallel API Requests), In-Memory Caching (TTL-based), Structured Logging, Global Error Handling.
- **Deployment & DevOps:** Docker (`Dockerfile` multi-stage, `docker-compose.yml`), Ready for PaaS (Render/Railway).

---

## 3. Component & Endpoint Specifications

### 3.1. Challenge 1: Website Metadata Extractor

- **Endpoint:** `POST /extract/website`
- **Request Body (JSON):**

```json
{
  "url": "https://paper.id"
}
```

**Requirements & Logic:**

1. Melakukan HTTP GET ke URL target menggunakan Go native `net/http` client.
2. Membaca dan melakukan parsing dokumen HTML mentah (direkomendasikan menggunakan `goquery`).
3. Mengekstraksi informasi: Tag `<title>`, `<meta name="description">`, `<link rel="canonical">`, `<link rel="shortcut icon">` atau favicon.
4. Mengekstraksi meta tag Open Graph (`og:title`, `og:description`, `og:image`).
5. Menggunakan Regular Expressions (Regex) untuk menyaring informasi email, nomor telepon, dan tautan social media (LinkedIn, Instagram, Twitter/X, Facebook) yang tertulis secara eksplisit pada teks halaman.

**Minimum Output Response (JSON):**

```json
{
  "url": "https://paper.id",
  "title": "String",
  "description": "String",
  "canonical": "String",
  "favicon": "String",
  "emails": ["String"],
  "phones": ["String"],
  "social_media": ["String"],
  "open_graph": {
    "title": "String",
    "description": "String",
    "image": "String"
  }
}
```

---

### 3.2. Challenge 2: Domain Intelligence

- **Endpoint:** `POST /extract/domain`
- **Request Body (JSON):**

```json
{
  "domain": "paper.id"
}
```

**Requirements & Logic:**

1. Membersihkan input untuk mengambil nama domain murni (tanpa `http://`, `https://`, atau trailing slash).
2. Melakukan HTTP request ke RDAP Open API: `https://rdap.org/domain/{domain}`.
3. Memetakan (mapping) JSON response masif dari RDAP ke dalam Go struct internal.
4. Mengambil nama registrar, melacak array events untuk mendapatkan tanggal registration, expiration, dan last changed/updated, serta menyaring array status dan nameservers.

**Minimum Output Response (JSON):**

```json
{
  "domain": "paper.id",
  "registrar": "String",
  "registered_at": "String",
  "expired_at": "String",
  "last_updated": "String",
  "status": ["String"],
  "nameservers": ["String"]
}
```

---

### 3.3. Challenge 3: Company Location Finder

- **Endpoint:** `POST /extract/location`
- **Request Body (JSON):**

```json
{
  "query": "PT Telkom Indonesia"
}
```

**Requirements & Logic:**

1. Melakukan URL Encoding pada parameter query.
2. Menembak OpenStreetMap Nominatim API: `https://nominatim.openstreetmap.org/search?q={query}&format=jsonv2`.
3. **CRITICAL:** Wajib menyertakan header `User-Agent` yang unik pada HTTP Request untuk menghindari blokir 403 Forbidden.
4. Mengambil data spasial dan informasi alamat berstruktur.

**Minimum Output Response (JSON):**

```json
{
  "display_name": "String",
  "latitude": "String",
  "longitude": "String",
  "importance": 0.0,
  "osm_type": "String",
  "address": {}
}
```

---

### 3.4. Final Integration: Company Information Core Engine

- **Endpoint:** `GET /company-information?domain={domain}`
- **Query Param:** `domain` (Contoh: `paper.id`)

**Requirements & Logic (Advanced Engineering):**

1. **Concurrency:** Wajib menggunakan `goroutine` dan `sync.WaitGroup` (atau `errgroup`) untuk mengeksekusi fungsi Website Metadata Extractor, Domain Intelligence, dan Company Location Finder secara paralel/bersamaan guna meminimalkan latency.
2. Khusus untuk pencarian lokasi (Location Finder), query awal dapat diturunkan secara cerdas dari nama domain atau entitas organisasi yang didapatkan dari respons RDAP.
3. **Resiliency & Fault Tolerance:** Jika salah satu API eksternal mengalami gangguan (down atau timeout), endpoint integrasi utama **TIDAK BOLEH** menghasilkan HTTP Status Code 500. Blok data yang gagal harus mengembalikan objek kosong/null, sedangkan data dari service yang sukses harus tetap ditampilkan.

**Minimum Output Response (JSON):**

```json
{
  "website": {},
  "domain": {},
  "location": {}
}
```

---

## 4. Non-Functional Requirements & Value-Added Features

### 4.1. In-Memory Caching

- Aplikasi wajib menerapkan mekanisme caching internal (misal: menggunakan `go-cache` atau implementasi `map` + `sync.RWMutex` dengan Time-To-Live (TTL) selama 10 menit).
- Setiap request berulang pada domain yang sama dalam rentang TTL wajib langsung dilayani dari memori tanpa melakukan HTTP request ulang ke internet untuk menghindari rate limiting dan menghemat resource.

### 4.2. Consistent Error Handling & Logging

- Seluruh error response harus dikembalikan secara konsisten dalam format JSON yang seragam.
- Sistem pencatatan log harus menggunakan structured logging (seperti `slog` bawaan Go atau `zap` logger) untuk mencatat setiap incoming request serta error tracking eksternal.

### 4.3. Git Workflow Regulations

- Proyek dikembangkan menggunakan Git secara bertahap.
- Dilarang keras melakukan one-shot commit. Setiap penyelesaian satu modul fungsional wajib diikuti oleh commit yang deskriptif menggunakan format Conventional Commits (Contoh: `feat: implement rdap connector for domain intelligence`).
