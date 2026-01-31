# Image Upload dengan Base64

## Overview
Backend Swipeup sekarang mendukung upload gambar menggunakan format Base64. Gambar yang diupload akan otomatis:
- Di-decode dari Base64
- Disimpan ke filesystem di folder `uploads/images/`
- Diberi nama unique menggunakan UUID
- Path-nya disimpan di database

## Format yang Didukung
- PNG
- JPEG/JPG
- GIF
- WebP

## Cara Penggunaan

### 1. Format Base64
Anda bisa mengirim gambar dalam 2 format:

#### Dengan Data URI (Recommended)
```json
{
  "foto": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
}
```

#### Tanpa Data URI
```json
{
  "foto": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
}
```

#### URL (Backward Compatible)
```json
{
  "foto": "https://example.com/photo.jpg"
}
```

### 2. Entities yang Mendukung Base64 Image
- **Siswa** (`foto` field)
- **Menu** (`foto` field)

### 3. API Endpoints

#### Create Siswa
```http
POST /api/siswa
Content-Type: application/json

{
  "nama_siswa": "Budi Santoso",
  "alamat": "Jl. Merdeka No. 123",
  "telp": "081234567890",
  "id_user": 1,
  "foto": "data:image/png;base64,..."
}
```

#### Update Siswa
```http
PUT /api/siswa/1
Content-Type: application/json

{
  "foto": "data:image/jpeg;base64,..."
}
```

**Note:** Saat update dengan foto baru, foto lama akan otomatis dihapus.

#### Create Menu
```http
POST /api/menu
Content-Type: application/json

{
  "nama_makanan": "Nasi Goreng",
  "harga": 15000,
  "jenis": "makanan",
  "deskripsi": "Nasi goreng spesial",
  "foto": "data:image/jpeg;base64,...",
  "id_stan": 1
}
```

#### Update Menu
```http
PUT /api/menu/1
Content-Type: application/json

{
  "foto": "data:image/png;base64,..."
}
```

### 4. Akses Gambar
Setelah gambar diupload, Anda bisa mengaksesnya melalui:

```
http://localhost:8080/uploads/images/<uuid>.<ext>
```

Contoh response dari API:
```json
{
  "success": true,
  "message": "Siswa created successfully",
  "data": {
    "id": 1,
    "nama_siswa": "Budi Santoso",
    "foto": "uploads/images/123e4567-e89b-12d3-a456-426614174000.png",
    ...
  }
}
```

Untuk menampilkan gambar:
```
http://localhost:8080/uploads/images/123e4567-e89b-12d3-a456-426614174000.png
```

## Image Lifecycle

### Create
1. Client mengirim Base64 string
2. Backend mendeteksi format Base64
3. Backend decode dan deteksi tipe image (magic bytes)
4. Backend generate UUID untuk filename
5. Backend save ke `uploads/images/<uuid>.<ext>`
6. Backend simpan path ke database

### Update
1. Client mengirim Base64 string baru
2. Backend fetch data existing
3. Backend delete foto lama (jika ada)
4. Backend save foto baru
5. Backend update path di database

### Delete
1. Backend fetch data
2. Backend delete record dari database
3. Backend delete file foto dari filesystem

## Technical Details

### File Structure
```
uploads/
  └── images/
      ├── 123e4567-e89b-12d3-a456-426614174000.png
      ├── 987fcdeb-51a2-43d7-b890-123456789abc.jpg
      └── ...
```

### Image Detection
Backend menggunakan "magic bytes" untuk mendeteksi tipe image:
- PNG: `\x89PNG`
- JPEG: `\xFF\xD8\xFF`
- GIF: `GIF87a` atau `GIF89a`
- WebP: `RIFF....WEBP`

### UUID Generation
Menggunakan `github.com/google/uuid` untuk generate unique filename.

## Error Handling

### Invalid Base64
```json
{
  "success": false,
  "message": "Failed to process image",
  "error": "illegal base64 data at input byte 123"
}
```

### Unsupported Format
```json
{
  "success": false,
  "message": "Failed to process image",
  "error": "unsupported image format"
}
```

### File Save Error
```json
{
  "success": false,
  "message": "Failed to process image",
  "error": "unable to create image file: ..."
}
```

## Best Practices

1. **Client Side:**
   - Compress image sebelum encode ke Base64
   - Resize image ke ukuran yang wajar (max 1920x1080)
   - Gunakan format JPEG untuk foto (lebih kecil)
   - Gunakan format PNG untuk logo/icon (lebih tajam)

2. **Security:**
   - Backend sudah validate image format via magic bytes
   - Only accept supported formats (PNG, JPEG, GIF, WebP)
   - Generated filename menggunakan UUID (prevent path traversal)

3. **Performance:**
   - Untuk list view, gunakan thumbnail atau lazy loading
   - Cache image di client side
   - Consider CDN untuk production

## Migration Guide

### Dari URL ke Base64

**Before:**
```javascript
// Client code
const siswa = {
  nama_siswa: "Budi",
  foto: "https://example.com/photo.jpg"
};
```

**After:**
```javascript
// Convert image to base64
const file = document.getElementById('fileInput').files[0];
const reader = new FileReader();
reader.onloadend = function() {
  const base64String = reader.result; // "data:image/png;base64,..."
  
  const siswa = {
    nama_siswa: "Budi",
    foto: base64String
  };
  
  // Send to API
  fetch('/api/siswa', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(siswa)
  });
};
reader.readAsDataURL(file);
```

### JavaScript Helper
```javascript
function imageToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

// Usage
const file = document.getElementById('photo').files[0];
const base64 = await imageToBase64(file);
```

## Testing

Lihat contoh di Bruno documentation:
- `docs/2-Siswa/Create Siswa.bru`
- `docs/2-Siswa/Update Siswa.bru`
- `docs/4-Menu/Create Menu.bru`
- `docs/4-Menu/Update Menu.bru`
