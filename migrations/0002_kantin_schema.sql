-- Kantin POS System Database Migration
-- Drop old tables first if they exist
DROP TABLE IF EXISTS transaction_items CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS menu_diskon CASCADE;
DROP TABLE IF EXISTS detail_transaksi CASCADE;
DROP TABLE IF EXISTS diskon CASCADE;
DROP TABLE IF EXISTS transaksi CASCADE;
DROP TABLE IF EXISTS menu CASCADE;
DROP TABLE IF EXISTS siswa CASCADE;
DROP TABLE IF EXISTS stan CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin_stan', 'siswa')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create stan table
CREATE TABLE stan (
    id SERIAL PRIMARY KEY,
    nama_stan VARCHAR(100) NOT NULL,
    nama_pemilik VARCHAR(100) NOT NULL,
    telp VARCHAR(20),
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create siswa table
CREATE TABLE siswa (
    id SERIAL PRIMARY KEY,
    nama_siswa VARCHAR(100) NOT NULL,
    alamat TEXT,
    telp VARCHAR(20),
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    foto VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create menu table
CREATE TABLE menu (
    id SERIAL PRIMARY KEY,
    nama_makanan VARCHAR(100) NOT NULL,
    harga DOUBLE PRECISION NOT NULL,
    jenis VARCHAR(20) NOT NULL CHECK (jenis IN ('makanan', 'minuman')),
    foto VARCHAR(255),
    deskripsi TEXT,
    id_stan INTEGER NOT NULL REFERENCES stan(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create transaksi table
CREATE TABLE transaksi (
    id SERIAL PRIMARY KEY,
    tanggal TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    id_stan INTEGER NOT NULL REFERENCES stan(id) ON DELETE CASCADE,
    id_siswa INTEGER NOT NULL REFERENCES siswa(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'belum dikonfirm' CHECK (status IN ('belum dikonfirm', 'dimasak', 'diantar', 'sampai')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create detail_transaksi table
CREATE TABLE detail_transaksi (
    id SERIAL PRIMARY KEY,
    id_transaksi INTEGER NOT NULL REFERENCES transaksi(id) ON DELETE CASCADE,
    id_menu INTEGER NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    qty INTEGER NOT NULL,
    harga_beli DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create diskon table
CREATE TABLE diskon (
    id SERIAL PRIMARY KEY,
    nama_diskon VARCHAR(100) NOT NULL,
    persentase_diskon DOUBLE PRECISION NOT NULL,
    tanggal_awal TIMESTAMP NOT NULL,
    tanggal_akhir TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create menu_diskon junction table
CREATE TABLE menu_diskon (
    id SERIAL PRIMARY KEY,
    id_menu INTEGER NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    id_diskon INTEGER NOT NULL REFERENCES diskon(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(id_menu, id_diskon)
);

-- Create indexes for better query performance
CREATE INDEX idx_stan_user ON stan(id_user);
CREATE INDEX idx_siswa_user ON siswa(id_user);
CREATE INDEX idx_menu_stan ON menu(id_stan);
CREATE INDEX idx_transaksi_stan ON transaksi(id_stan);
CREATE INDEX idx_transaksi_siswa ON transaksi(id_siswa);
CREATE INDEX idx_transaksi_status ON transaksi(status);
CREATE INDEX idx_detail_transaksi ON detail_transaksi(id_transaksi);
CREATE INDEX idx_detail_menu ON detail_transaksi(id_menu);
CREATE INDEX idx_menu_diskon_menu ON menu_diskon(id_menu);
CREATE INDEX idx_menu_diskon_diskon ON menu_diskon(id_diskon);

-- Insert sample data
-- Sample users
INSERT INTO users (username, password, role) VALUES
('admin_stan1', '$2a$10$dummyhash1', 'admin_stan'),
('admin_stan2', '$2a$10$dummyhash2', 'admin_stan'),
('siswa1', '$2a$10$dummyhash3', 'siswa'),
('siswa2', '$2a$10$dummyhash4', 'siswa');

-- Sample stan
INSERT INTO stan (nama_stan, nama_pemilik, telp, id_user) VALUES
('Warung Nasi Bu Yani', 'Bu Yani', '081234567890', 1),
('Kedai Minuman Segar', 'Pak Budi', '081234567891', 2);

-- Sample siswa
INSERT INTO siswa (nama_siswa, alamat, telp, id_user, foto) VALUES
('Ahmad Fauzi', 'Jl. Merdeka No. 1', '081234567892', 3, 'ahmad.jpg'),
('Siti Nurhaliza', 'Jl. Sudirman No. 2', '081234567893', 4, 'siti.jpg');

-- Sample menu
INSERT INTO menu (nama_makanan, harga, jenis, foto, deskripsi, id_stan) VALUES
('Nasi Goreng', 15000, 'makanan', 'nasi-goreng.jpg', 'Nasi goreng spesial dengan telur', 1),
('Mie Ayam', 12000, 'makanan', 'mie-ayam.jpg', 'Mie ayam dengan kuah gurih', 1),
('Ayam Geprek', 18000, 'makanan', 'ayam-geprek.jpg', 'Ayam geprek pedas level 5', 1),
('Es Teh Manis', 5000, 'minuman', 'es-teh.jpg', 'Es teh manis segar', 2),
('Jus Jeruk', 8000, 'minuman', 'jus-jeruk.jpg', 'Jus jeruk segar tanpa gula', 2),
('Kopi Susu', 10000, 'minuman', 'kopi-susu.jpg', 'Kopi susu dingin', 2);

-- Sample diskon
INSERT INTO diskon (nama_diskon, persentase_diskon, tanggal_awal, tanggal_akhir) VALUES
('Diskon Akhir Tahun', 20, '2025-12-01', '2025-12-31'),
('Promo Pagi', 10, '2026-01-01', '2026-06-30');

-- Sample menu_diskon
INSERT INTO menu_diskon (id_menu, id_diskon) VALUES
(1, 2),
(2, 2);
