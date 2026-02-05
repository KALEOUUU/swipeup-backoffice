# Swipeup BE API - Dokumentasi Bruno

Dokumentasi lengkap API Swipeup Backend menggunakan **Bruno** untuk testing dan dokumentasi API.

## 📋 Daftar Fitur

### 1. Users (Pengguna)
Mengelola pengguna sistem (Admin Stan dan Siswa).
- ✅ [Create User](1-Users/Create%20User.bru) - Membuat pengguna baru
- ✅ [Get All Users](1-Users/Get%20All%20Users.bru) - Mendapatkan semua pengguna
- ✅ [Get User by ID](1-Users/Get%20User%20by%20ID.bru) - Mendapatkan pengguna berdasarkan ID

### 2. Siswa (Pelajar)
Mengelola data siswa yang akan bertransaksi di kantin.
- ✅ [Create Siswa](2-Siswa/Create%20Siswa.bru) - Membuat data siswa baru
- ✅ [Get All Siswa](2-Siswa/Get%20All%20Siswa.bru) - Mendapatkan semua siswa
- ✅ [Get Siswa by ID](2-Siswa/Get%20Siswa%20by%20ID.bru) - Mendapatkan siswa berdasarkan ID
- ✅ [Get Siswa by User ID](2-Siswa/Get%20Siswa%20by%20User%20ID.bru) - Mendapatkan siswa berdasarkan User ID
- ✅ [Update Siswa](2-Siswa/Update%20Siswa.bru) - Mengupdate data siswa
- ✅ [Delete Siswa](2-Siswa/Delete%20Siswa.bru) - Menghapus data siswa

### 3. Stan (Penjual/Toko)
Mengelola data stan/toko yang menjual di kantin.
- ✅ [Create Stan](3-Stan/Create%20Stan.bru) - Membuat stan baru
- ✅ [Get All Stan](3-Stan/Get%20All%20Stan.bru) - Mendapatkan semua stan
- ✅ [Get Stan by ID](3-Stan/Get%20Stan%20by%20ID.bru) - Mendapatkan stan berdasarkan ID
- ✅ [Get Stan by User ID](3-Stan/Get%20Stan%20by%20User%20ID.bru) - Mendapatkan stan berdasarkan User ID
- ✅ [Update Stan](3-Stan/Update%20Stan.bru) - Mengupdate data stan
- ✅ [Delete Stan](3-Stan/Delete%20Stan.bru) - Menghapus data stan

### 4. Menu (Makanan/Minuman)
Mengelola menu makanan dan minuman yang dijual di setiap stan.
- ✅ [Create Menu](4-Menu/Create%20Menu.bru) - Membuat menu baru
- ✅ [Get All Menu](4-Menu/Get%20All%20Menu.bru) - Mendapatkan semua menu
- ✅ [Get Menu by ID](4-Menu/Get%20Menu%20by%20ID.bru) - Mendapatkan menu berdasarkan ID
- ✅ [Get Menu by Stan ID](4-Menu/Get%20Menu%20by%20Stan%20ID.bru) - Mendapatkan menu berdasarkan Stan ID
- ✅ [Search Menu by Name](4-Menu/Search%20Menu%20by%20Name.bru) - Mencari menu berdasarkan nama
- ✅ [Update Menu](4-Menu/Update%20Menu.bru) - Mengupdate data menu
- ✅ [Delete Menu](4-Menu/Delete%20Menu.bru) - Menghapus data menu

### 5. Transaksi (Pembelian)
Mengelola transaksi pembelian siswa.
- ✅ [Create Transaksi](5-Transaksi/Create%20Transaksi.bru) - Membuat transaksi pembelian baru
- ✅ [Get All Transaksi](5-Transaksi/Get%20All%20Transaksi.bru) - Mendapatkan semua transaksi
- ✅ [Get Transaksi by ID](5-Transaksi/Get%20Transaksi%20by%20ID.bru) - Mendapatkan transaksi berdasarkan ID
- ✅ [Get Transaksi by Siswa ID](5-Transaksi/Get%20Transaksi%20by%20Siswa%20ID.bru) - Mendapatkan transaksi siswa
- ✅ [Update Transaksi Status](5-Transaksi/Update%20Transaksi%20Status.bru) - Mengupdate status transaksi

### 6. Diskon (Potongan Harga)
Mengelola diskon dengan dukungan 2 level (Global & Per Stan).
- ✅ [Create Diskon](6-Diskon/Create%20Diskon.bru) - Membuat diskon baru
- ✅ [Create Diskon Stan](6-Diskon/Create%20Diskon%20Stan.bru) - Membuat diskon per stan
- ✅ [Create Diskon Menu](6-Diskon/Create%20Diskon%20Menu.bru) - Membuat diskon per menu
- ✅ [Get All Diskon](6-Diskon/Get%20All%20Diskon.bru) - Mendapatkan semua diskon
- ✅ [Get Diskon by ID](6-Diskon/Get%20Diskon%20by%20ID.bru) - Mendapatkan diskon berdasarkan ID
- ✅ [Get Global Diskon](6-Diskon/Get%20Global%20Diskon.bru) - Mendapatkan diskon global
- ✅ [Get Diskon by Stan](6-Diskon/Get%20Diskon%20by%20Stan.bru) - Mendapatkan diskon berdasarkan stan
- ✅ [Get Active Diskon by Stan](6-Diskon/Get%20Active%20Diskon%20by%20Stan.bru) - Mendapatkan diskon aktif untuk stan
- ✅ [Update Diskon](6-Diskon/Update%20Diskon.bru) - Mengupdate data diskon
- ✅ [Update Diskon Stan](6-Diskon/Update%20Diskon%20Stan.bru) - Mengupdate diskon stan
- ✅ [Update Diskon Menu](6-Diskon/Update%20Diskon%20Menu.bru) - Mengupdate diskon menu
- ✅ [Delete Diskon](6-Diskon/Delete%20Diskon.bru) - Menghapus data diskon
- ✅ [Delete Diskon Stan](6-Diskon/Delete%20Diskon%20Stan.bru) - Menghapus diskon stan
- ✅ [Delete Diskon Menu](6-Diskon/Delete%20Diskon%20Menu.bru) - Menghapus diskon menu

### 7. Cart (Keranjang Belanja)
Mengelola keranjang belanja siswa sebelum checkout.
- ✅ [Add to Cart](7-Cart/Add%20to%20Cart.bru) - Menambah item ke keranjang
- ✅ [Get Cart](7-Cart/Get%20Cart.bru) - Mendapatkan isi keranjang siswa
- ✅ [Update Cart Item](7-Cart/Update%20Cart%20Item.bru) - Mengupdate jumlah item di keranjang
- ✅ [Remove from Cart](7-Cart/Remove%20from%20Cart.bru) - Menghapus item dari keranjang
- ✅ [Clear Cart](7-Cart/Clear%20Cart.bru) - Mengosongkan seluruh keranjang
- ✅ [Checkout Cart](7-Cart/Checkout%20Cart.bru) - Checkout keranjang ke transaksi

### 8. Auth (Autentikasi)
Mengelola autentikasi dan otorisasi pengguna.
- ✅ [Register](8-Auth/Register.bru) - Registrasi pengguna baru
- ✅ [Register Admin Stan](8-Auth/Register%20Admin%20Stan.bru) - Registrasi admin stan dengan stan
- ✅ [Login](8-Auth/Login.bru) - Login dan dapatkan token JWT
- ✅ [Get Profile](8-Auth/Profile.bru) - Mendapatkan profil pengguna yang terautentikasi

### 9. Activity Logs (Log Aktivitas)
Mengelola dan memonitor log aktivitas sistem.
- ✅ [Get User Activities](9-Activity-Logs/Get-User-Activities.bru) - Mendapatkan aktivitas user tertentu
- ✅ [Get All Activities](9-Activity-Logs/Get-All-Activities.bru) - Mendapatkan semua aktivitas dengan filter
- ✅ [Get Activities by Date Range](9-Activity-Logs/Get-Activities-by-Date-Range.bru) - Mendapatkan aktivitas berdasarkan rentang tanggal
- ✅ [Get Activity Stats](9-Activity-Logs/Get-Activity-Stats.bru) - Mendapatkan statistik aktivitas
- ✅ [Clean Old Logs](9-Activity-Logs/Clean-Old-Logs.bru) - Membersihkan log aktivitas lama

## 🚀 Cara Menggunakan

1. **Install Bruno**
   - Download dari [bruno.app](https://www.usebruno.com/)
   - Install sesuai OS Anda

2. **Import Koleksi**
   - Buka Bruno
   - Klik "Open Collection"
   - Pilih folder `docs/` dari project ini
   - Koleksi akan dimuat otomatis

3. **Konfigurasi Base URL**
   - Base URL default: `http://localhost:8080`
   - Sesuaikan dengan server Anda jika berbeda

4. **Jalankan Request**
   - Pilih request yang ingin di-test
   - Klik tombol "Send" atau tekan Cmd+Enter (Mac) / Ctrl+Enter (Windows/Linux)
   - Lihat response di panel sebelah kanan

## 📊 Struktur Request

Setiap request memiliki struktur standar:
- **Method**: GET, POST, PUT, DELETE
- **URL**: `http://localhost:8080/api/{endpoint}`
- **Body** (untuk POST/PUT): JSON dengan data yang diperlukan
- **Tests**: Validasi otomatis untuk setiap response

## ⚠️ Catatan Penting

- Pastikan backend server sudah berjalan sebelum menjalankan request
- Default port adalah `8080`, sesuaikan jika berbeda
- Beberapa endpoint memerlukan data yang sudah ada (misal: ID yang valid)
- Gunakan ID yang real saat testing untuk menghindari error "Not Found"

## 🔄 Workflow Testing Rekomendasi

1. **Setup Autentikasi**
   - Register User (untuk Admin Stan dan Siswa)
   - Login untuk mendapatkan JWT Token
   - Get Profile untuk verifikasi token

2. **Setup Data Awal**
   - Create User (untuk Admin Stan dan Siswa)
   - Create Stan (dengan user Admin Stan)
   - Create Siswa (dengan user Siswa)

3. **Setup Menu & Diskon**
   - Create Menu (untuk setiap Stan)
   - Create Diskon Global (opsional, oleh superadmin)
   - Create Diskon Stan (oleh admin stan)
   - Create Diskon Menu (untuk menu tertentu)

4. **Testing Cart & Transaksi**
   - Add to Cart (siswa menambah item ke keranjang)
   - Get Cart (melihat isi keranjang)
   - Update Cart Item (mengubah jumlah item)
   - Checkout Cart (konversi ke transaksi)
   - Create Transaksi (konfirmasi transaksi)
   - Get Transaksi (untuk verifikasi)
   - Update Transaksi Status

5. **Testing Activity Logs**
   - Get User Activities (melihat aktivitas user)
   - Get All Activities (monitor sistem)
   - Get Activity Stats (dashboard metrics)
   - Get Activities by Date Range (analisis historis)

6. **Testing Fitur Lainnya**
   - Update & Delete untuk setiap entitas
   - Get dengan filter (by ID, by Stan, etc)
   - Remove from Cart & Clear Cart
   - Clean Old Logs (maintenance)

## 📝 Tipe Role User

- **superadmin**: Pengguna dengan akses penuh ke semua fitur termasuk diskon global
- **admin_stan**: Pengguna yang mengelola stan/toko, menu, dan diskon per stan
- **siswa**: Pengguna yang bertransaksi (membeli), mengelola cart, dan melihat aktivitas

## 💾 Tips & Tricks

- Gunakan environment variables di Bruno untuk menyimpan ID yang sering digunakan
- Setiap test memiliki assertions untuk validasi response
- Lihat console Bruno untuk melihat error detail jika ada
- **Autentikasi**: Simpan JWT token di environment variable `{{token}}` untuk request yang memerlukan auth
- **Cart**: Gunakan `Add to Cart` untuk menambah item, lalu `Get Cart` untuk verifikasi sebelum checkout
- **Diskon**: Perhatikan perbedaan antara diskon global (superadmin) dan diskon per stan (admin_stan)
- **Activity Logs**: Gunakan endpoint stats untuk dashboard metrics dan monitoring sistem
- **Image Upload**: Gunakan format Base64 untuk upload foto siswa dan menu (lihat README_IMAGE_UPLOAD.md)

## 📚 Dokumentasi Tambahan

- **[Diskon 2 Level](DISKON_2_LEVEL.md)** - Dokumentasi lengkap sistem diskon dengan 2 level (Global & Per Stan)
- **[Image Upload](README_IMAGE_UPLOAD.md)** - Panduan upload gambar menggunakan format Base64

---

**Created for Swipeup - Kantin POS System**
