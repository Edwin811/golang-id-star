Fullstack Microservices - IDStar Technical Project
Proyek ini dibangun sebagai bagian dari persiapan teknis untuk posisi Full Stack Developer di IDStar Group. Sistem ini mengimplementasikan arsitektur microservices menggunakan Golang di sisi backend dan React.js di sisi frontend, dengan seluruh infrastruktur berjalan di atas Docker.

🚀 Fitur Utama
Microservices Architecture: Pemisahan tanggung jawab antara service-employee dan service-customer.

Containerization: Menggunakan Docker dan Docker Compose untuk orkestrasi layanan dan standarisasi lingkungan.

Inter-service Communication: Komunikasi antar-layanan melalui jaringan internal Docker.

CI/CD Pipeline: Otomatisasi pengujian dan build menggunakan GitHub Actions.

Database Integration: Integrasi dengan MySQL menggunakan Raw SQL untuk performa tinggi.

🛠️ Tech Stack
Backend: Golang 1.22+

Frontend: React.js

Database: MySQL

DevOps: Docker, Docker Compose, GitHub Actions

📂 Struktur Proyek
.github/workflows/: Konfigurasi CI/CD.

service-employee/: Layanan manajemen data karyawan (Port 8080).

service-customer/: Layanan manajemen data pelanggan (Port 8081).

frontend-idstar/: Aplikasi dashboard React.

⚙️ Cara Menjalankan (Lokal)
Pastikan Anda memiliki Docker Desktop yang berjalan di mesin Anda (Dioptimalkan untuk Mac M1).

Clone repository ini.

Jalankan perintah orkestrasi:

Bash
docker-compose up --build
Akses Frontend di http://localhost:3000.