# Fitur Diskon 2 Level

## 📋 Overview

Sistem diskon mendukung 2 level pengaturan:

### 1. **Diskon Global** (Superadmin)
- Diatur oleh superadmin
- Berlaku untuk **semua stan**
- `tipe_diskon: "global"`
- `id_stan: null`

### 2. **Diskon Per Stan** (Admin Stan)
- Diatur oleh admin stan masing-masing
- Berlaku **hanya untuk stan tertentu**
- `tipe_diskon: "stan"`
- `id_stan: <ID_STAN>`

## 🔑 Field Model Diskon

```json
{
  "id": 1,
  "nama_diskon": "Promo Ramadhan",
  "persentase_diskon": 15,
  "tanggal_awal": "2025-03-01T00:00:00Z",
  "tanggal_akhir": "2025-03-31T23:59:59Z",
  "tipe_diskon": "global",  // "global" atau "stan"
  "id_stan": null,           // null untuk global, ID untuk stan
  "stan": null               // relasi ke stan (jika ada)
}
```

## 🚀 API Endpoints

### Create Diskon Global
```http
POST /api/diskon
Content-Type: application/json

{
  "nama_diskon": "Diskon Ramadhan",
  "persentase_diskon": 10,
  "tanggal_awal": "2025-03-01T00:00:00Z",
  "tanggal_akhir": "2025-03-31T23:59:59Z",
  "tipe_diskon": "global"
}
```

### Create Diskon Per Stan
```http
POST /api/diskon
Content-Type: application/json

{
  "nama_diskon": "Promo Stan A",
  "persentase_diskon": 15,
  "tanggal_awal": "2025-02-01T00:00:00Z",
  "tanggal_akhir": "2025-02-28T23:59:59Z",
  "tipe_diskon": "stan",
  "id_stan": 1
}
```

### Get Global Diskon
```http
GET /api/diskon/global
```

### Get Diskon by Stan
```http
GET /api/diskon/by-stan?stan_id=1
```
Mengembalikan **hanya** diskon yang dibuat oleh stan tersebut.

### Get Active Diskon by Stan
```http
GET /api/diskon/active-by-stan?stan_id=1
```
Mengembalikan:
- Diskon **global** yang sedang aktif
- Diskon **stan** yang sedang aktif untuk stan tersebut

## 📊 Use Cases

### 1. Superadmin membuat diskon global
```json
{
  "nama_diskon": "Diskon Hari Kemerdekaan",
  "persentase_diskon": 17,
  "tanggal_awal": "2025-08-17T00:00:00Z",
  "tanggal_akhir": "2025-08-17T23:59:59Z",
  "tipe_diskon": "global"
}
```
→ Berlaku untuk **semua stan**

### 2. Admin Stan A membuat diskon internal
```json
{
  "nama_diskon": "Promo Ulang Tahun Stan A",
  "persentase_diskon": 20,
  "tanggal_awal": "2025-06-01T00:00:00Z",
  "tanggal_akhir": "2025-06-30T23:59:59Z",
  "tipe_diskon": "stan",
  "id_stan": 1
}
```
→ Berlaku **hanya untuk Stan A**

### 3. Siswa melihat diskon untuk Stan A
`GET /api/diskon/active-by-stan?stan_id=1`

Akan mendapat:
- Diskon global (jika ada dan aktif)
- Diskon internal Stan A (jika ada dan aktif)

## ⚠️ Validasi

1. **Diskon global**: `id_stan` harus `null`
2. **Diskon stan**: `id_stan` harus diisi dan valid
3. **Format tanggal**: Gunakan RFC3339 (`YYYY-MM-DDTHH:MM:SSZ`)
4. **Persentase**: 0-100

## 🔄 Migration

Jalankan migration untuk menambahkan field baru:
```bash
psql -U <user> -d <database> -f migrations/0003_add_diskon_level.sql
```

Atau restart aplikasi (GORM auto-migrate akan menambahkan column baru).

## 📝 Notes

- Diskon lama akan otomatis menjadi diskon global
- Stan tidak bisa edit/delete diskon global
- Superadmin bisa edit/delete semua diskon
- Diskon bisa di-assign ke menu via endpoint assign
