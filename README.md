# Tes server Go di Render

Server HTTP sederhana tanpa dependency eksternal atau database.
GET / dan GET /health mengembalikan JSON status ok dan waktu UTC.
Ini menguji akses ke aplikasi Go yang di-host Render setelah deploy.

## Jalankan lokal (perlu Go 1.22 atau lebih baru)

go run .

Buka http://localhost:10000/health atau jalankan:
Invoke-RestMethod http://localhost:10000/health

## Deploy manual melalui dashboard Render

1. Upload isi folder ini ke repository GitHub kamu.
2. Di Render pilih New > Web Service, lalu hubungkan repository tersebut.
3. Language/Runtime: Go.
4. Root Directory: kosong jika go.mod berada di root repository.
5. Build Command: go build -o bin/server .
6. Start Command: ./bin/server
7. Instance Type: Free. Tidak perlu database atau disk.
8. Health Check Path: /health.
9. Klik Deploy Web Service, tunggu status Live.
10. Buka https://NAMA-SERVICE.onrender.com/health.

Tes melalui PowerShell (ganti nama service):
Invoke-RestMethod https://NAMA-SERVICE.onrender.com/health -TimeoutSec 120

Hasil sukses memiliki status: ok dan message: Koneksi ke server Go berhasil!
Server membaca PORT dari Render dan bind ke 0.0.0.0.

Dokumentasi gratis: https://render.com/docs/free
Render Free tidur setelah 15 menit tanpa trafik. Akses berikutnya bisa
memerlukan sekitar satu menit untuk menghidupkan server kembali.
Jika halaman meminta kartu, catat pesan persisnya dan pastikan instance
Free yang dipilih; masuk dashboard belum membuktikan deploy selesai.

## Status validasi

Kode belum dijalankan lokal karena Go belum tersedia di lingkungan pembuatannya.
Belum di-upload ke GitHub atau di-deploy ke akun Render.
