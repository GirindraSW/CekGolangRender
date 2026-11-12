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

## Deploy ke Vercel Hobby (project yang sama)

Project menggunakan net/http dan main.go di root, dengan go.mod Go 1.22.
Vercel mendukung backend Go langsung sejak 2 April 2026. Konfigurasi
vercel.json memilih framework go. Tidak diperlukan folder api, handler
serverless terpisah, rewrite, dependency baru, atau perubahan endpoint.
Kode server dan panduan Render di atas tetap berlaku.

Dokumentasi resmi:
- https://vercel.com/changelog/zero-configuration-go-backend-support
- https://vercel.com/docs/project-configuration/vercel-json
- https://vercel.com/docs/plans/hobby
- https://vercel.com/docs/git

Langkah deploy manual:
1. Commit dan push file project ini ke repository GitHub milikmu.
2. Login Vercel dengan GitHub. Pilih akun/team dengan paket Hobby gratis
   untuk tes pribadi nonkomersial. Jangan pilih trial Pro atau upgrade.
3. Pilih Add New > Project > Import Git Repository.
4. Jika repository belum muncul, Configure GitHub App dan berikan akses
   ke repository ini, lalu Import.
5. Root Directory: folder yang berisi go.mod, main.go, dan vercel.json.
   Gunakan root repository jika ketiga file itu berada di root.
6. Pastikan Framework Preset: Go. Gunakan build/install/output default
   dari preset. Jangan salin start/build command Render ke Vercel.
7. Environment Variables: tidak ada yang perlu ditambahkan. PORT dibaca
   dari environment runtime; jangan mengunci nilainya di dashboard.
   Fallback 10000 hanya dipakai saat PORT tidak disediakan, misalnya lokal.
8. Klik Deploy. Tunggu status Ready dan periksa Build Logs jika gagal.
9. Buka https://NAMA-PROJECT.vercel.app/ dan /health.
10. Untuk tes dari PowerShell:
    Invoke-RestMethod https://NAMA-PROJECT.vercel.app/health -TimeoutSec 120
    Hasil sukses: status ok, message Koneksi ke server Go berhasil!, time UTC.
11. Path lain harus 404 dan POST /health harus 405. Runtime Logs akan
    menampilkan request GET/HEAD yang berhasil diproses server.

Pemeriksaan lokal setelah Go tersedia:
    go test ./...
    go build ./...
    go run .
    Invoke-RestMethod http://localhost:10000/health

Status persiapan Vercel:
- main.go dan go.mod tidak diubah; endpoint dan logic tetap sama.
- vercel.json ditambahkan; .vercel/ diabaikan oleh Git.
- Tidak ada secret/credential atau environment variable wajib tambahan.
- Go belum ditemukan di PATH/lokasi instalasi umum pada saat persiapan.
  go test ./... dan go build ./... belum dapat dijalankan sampai Go tersedia.
- Belum melakukan push, deploy, atau perubahan pengaturan akun Vercel.
- Hobby gratis memiliki kuota dan hanya untuk penggunaan pribadi
  nonkomersial. Konfigurasi ini tidak mengaktifkan produk berbayar.

## Tes Create: POST /items

Endpoint lama / dan /health tetap tersedia. Endpoint baru menerima JSON:
{"name":"Kursi"}

Tes setelah commit, push, dan deployment Vercel terbaru Ready:

```powershell
$body = @{ name = 'Kursi' } | ConvertTo-Json
Invoke-RestMethod -Uri 'https://cek-golang-render.vercel.app/items' -Method Post -ContentType 'application/json' -Body $body
```

Hasil: HTTP 201 dengan data.id, data.name, data.created_at, dan storage memory.
name wajib diisi dan maksimal 100 karakter. Body maksimal 4096 byte.
GET /items mengembalikan 405; gunakan PowerShell untuk mengirim POST.

Penyimpanan hanya di memori instance, bukan database. Data bisa hilang saat
restart/redeploy dan tidak dibagikan antar-instance Vercel. Ini hanya tes
Create; belum menyediakan Read, Update, atau Delete. Untuk data permanen,
perlu database. Tidak membutuhkan secret atau environment variable baru.

Jalankan pemeriksaan jika Go tersedia:
go test ./...
go build ./...

Untuk menjalankan lokal gunakan go run . agar items.go ikut dikompilasi.
