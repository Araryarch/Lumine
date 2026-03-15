# Requirements Document

## Introduction

Lumine Modern DevStack adalah aplikasi Docker management tool berbasis TUI (Terminal User Interface) yang dibangun dengan Go. Aplikasi ini dirancang sebagai alternatif modern dan cross-platform dari Laragon (Windows) dan Laravel Herd untuk mengelola environment development web secara lokal. Lumine menyediakan orchestration service Docker, manajemen proyek otomatis, database management, dan interface TUI yang intuitif untuk meningkatkan produktivitas developer.

## Glossary

- **Lumine**: Aplikasi Docker management tool berbasis TUI yang menjadi subjek dari dokumen ini
- **Service**: Container Docker yang menjalankan komponen infrastruktur (Nginx, MySQL, Redis, dll)
- **Service_Manager**: Komponen yang mengelola lifecycle service Docker
- **Version_Manager**: Komponen yang mengelola versi runtime (PHP, Node.js, Python)
- **Project_Manager**: Komponen yang mengelola proyek development web
- **TUI**: Terminal User Interface, interface berbasis teks di terminal
- **Virtual_Host**: Konfigurasi domain lokal untuk mengakses proyek (misal: project.test)
- **Hosts_File**: File sistem operasi yang memetakan domain ke IP address (/etc/hosts atau C:\Windows\System32\drivers\etc\hosts)
- **Port_Conflict**: Kondisi dimana dua service mencoba menggunakan port yang sama
- **Daemon_Mode**: Mode operasi dimana aplikasi berjalan di background
- **SSL_Certificate**: Sertifikat digital untuk mengaktifkan HTTPS pada domain lokal
- **Tunnel**: Koneksi yang mengekspos localhost ke internet publik
- **Dependency_Tool**: Tool eksternal yang dibutuhkan proyek (Composer, NPM, PNPM, Yarn)
- **Database_Connection**: Konfigurasi koneksi ke database (MySQL, PostgreSQL, SQLite)
- **Query_Log**: Log yang merekam semua query database yang dieksekusi
- **Toast_Notification**: Notifikasi sementara yang muncul di interface TUI
- **Admin_Privilege**: Hak akses administrator/root yang diperlukan untuk operasi sistem tertentu

## Requirements

### Requirement 1: Service Lifecycle Management

**User Story:** Sebagai developer, saya ingin mengelola lifecycle service Docker (start/stop/restart), sehingga saya dapat mengontrol service yang berjalan sesuai kebutuhan development.

#### Acceptance Criteria

1. WHEN pengguna memilih action "start" pada service, THE Service_Manager SHALL memulai container Docker untuk service tersebut dalam waktu maksimal 10 detik
2. WHEN pengguna memilih action "stop" pada service yang sedang berjalan, THE Service_Manager SHALL menghentikan container Docker dalam waktu maksimal 5 detik
3. WHEN pengguna memilih action "restart" pada service, THE Service_Manager SHALL menghentikan kemudian memulai ulang container Docker dalam waktu maksimal 15 detik
4. THE Service_Manager SHALL mendukung service berikut: Nginx, Apache, Caddy, PHP-FPM, MySQL, PostgreSQL, Redis, MailHog, MongoDB, Elasticsearch
5. WHEN service gagal dijalankan, THE Service_Manager SHALL menampilkan pesan error yang deskriptif dengan informasi penyebab kegagalan
6. THE Lumine SHALL menyimpan status terakhir setiap service (running atau stopped) dan menampilkannya di TUI

### Requirement 2: Multi-Version Runtime Support

**User Story:** Sebagai developer, saya ingin mengganti versi PHP atau Node.js secara on-the-fly, sehingga saya dapat menjalankan proyek dengan versi runtime yang berbeda tanpa conflict.

#### Acceptance Criteria

1. THE Version_Manager SHALL mendukung minimal 5 versi PHP (7.4, 8.0, 8.1, 8.2, 8.3)
2. THE Version_Manager SHALL mendukung minimal 4 versi Node.js (16, 18, 20, 21)
3. WHEN pengguna mengganti versi PHP, THE Version_Manager SHALL menghentikan container PHP lama dan memulai container PHP baru dengan versi yang dipilih dalam waktu maksimal 15 detik
4. WHEN pengguna mengganti versi Node.js, THE Version_Manager SHALL mengupdate symlink atau environment variable untuk menggunakan versi yang dipilih dalam waktu maksimal 5 detik
5. THE Version_Manager SHALL menampilkan versi runtime yang sedang aktif di dashboard TUI
6. WHEN pergantian versi gagal, THE Version_Manager SHALL mengembalikan ke versi sebelumnya dan menampilkan pesan error

### Requirement 3: Automatic Port Conflict Resolution

**User Story:** Sebagai developer, saya ingin sistem mendeteksi dan menyelesaikan port conflict secara otomatis, sehingga saya tidak perlu manual mencari port yang tersedia.

#### Acceptance Criteria

1. WHEN service akan dijalankan, THE Service_Manager SHALL memeriksa apakah port yang dibutuhkan sudah digunakan oleh proses lain
2. IF port sudah digunakan, THEN THE Service_Manager SHALL mencari port alternatif yang tersedia dalam range 8000-9000
3. WHEN port alternatif ditemukan, THE Service_Manager SHALL menggunakan port tersebut dan menampilkan notifikasi perubahan port kepada pengguna
4. THE Service_Manager SHALL menyimpan mapping port yang digunakan setiap service dan menampilkannya di TUI
5. IF tidak ada port tersedia dalam range 8000-9000, THEN THE Service_Manager SHALL menampilkan error dan membatalkan operasi start service

### Requirement 4: Daemon Mode with Auto-Restart

**User Story:** Sebagai developer, saya ingin aplikasi berjalan di background dan memastikan service tetap hidup, sehingga saya tidak perlu khawatir service mati saat bekerja.

#### Acceptance Criteria

1. THE Lumine SHALL menyediakan command "lumine daemon start" untuk menjalankan aplikasi dalam Daemon_Mode
2. WHILE Daemon_Mode aktif, THE Lumine SHALL memonitor status semua service yang seharusnya running setiap 30 detik
3. WHEN service terdeteksi crash atau stopped secara tidak sengaja, THE Lumine SHALL otomatis me-restart service tersebut dalam waktu maksimal 10 detik
4. THE Lumine SHALL mencatat setiap auto-restart ke log file dengan timestamp dan alasan restart
5. THE Lumine SHALL menyediakan command "lumine daemon stop" untuk menghentikan Daemon_Mode
6. WHEN Daemon_Mode dihentikan, THE Lumine SHALL menghentikan monitoring tetapi tidak menghentikan service yang sedang berjalan

### Requirement 5: Project Scaffolding via TUI

**User Story:** Sebagai developer, saya ingin membuat proyek baru melalui menu TUI, sehingga saya dapat dengan cepat setup proyek tanpa menjalankan command manual.

#### Acceptance Criteria

1. THE Project_Manager SHALL menyediakan menu TUI untuk membuat proyek baru dengan pilihan framework: Laravel, WordPress, Symfony, React, Vue, Next.js, Static HTML
2. WHEN pengguna memilih framework, THE Project_Manager SHALL menampilkan form input untuk nama proyek dan konfigurasi tambahan (versi PHP/Node.js)
3. WHEN pengguna submit form, THE Project_Manager SHALL membuat direktori proyek dan menjalankan scaffolding command yang sesuai (composer create-project, npx create-react-app, dll)
4. THE Project_Manager SHALL menampilkan progress bar atau loading indicator selama proses scaffolding
5. WHEN scaffolding selesai, THE Project_Manager SHALL menampilkan notifikasi sukses dengan informasi path proyek dan URL akses
6. IF scaffolding gagal, THEN THE Project_Manager SHALL menampilkan error message dan membersihkan direktori proyek yang sudah dibuat

### Requirement 6: Automatic Virtual Host Configuration

**User Story:** Sebagai developer, saya ingin sistem otomatis membuat domain lokal untuk proyek baru, sehingga saya dapat langsung mengakses proyek melalui browser tanpa konfigurasi manual.

#### Acceptance Criteria

1. WHEN proyek baru dibuat, THE Project_Manager SHALL otomatis generate Virtual_Host configuration dengan domain pattern {project_name}.test
2. THE Project_Manager SHALL menambahkan entry "127.0.0.1 {project_name}.test" ke Hosts_File sistem operasi
3. THE Project_Manager SHALL membuat konfigurasi Nginx atau Apache untuk Virtual_Host tersebut dan me-reload web server
4. WHEN proyek dihapus, THE Project_Manager SHALL menghapus entry dari Hosts_File dan konfigurasi Virtual_Host
5. THE Project_Manager SHALL menampilkan URL akses proyek (http://{project_name}.test) di dashboard TUI
6. IF pengguna tidak memiliki Admin_Privilege, THEN THE Project_Manager SHALL menampilkan instruksi untuk menjalankan command dengan sudo/admin

### Requirement 7: Automatic SSL Certificate Generation

**User Story:** Sebagai developer, saya ingin sistem otomatis menerbitkan SSL certificate untuk domain lokal, sehingga saya dapat mengakses proyek dengan HTTPS tanpa warning browser.

#### Acceptance Criteria

1. WHEN Virtual_Host dibuat, THE Project_Manager SHALL otomatis generate self-signed SSL_Certificate untuk domain tersebut menggunakan mkcert atau openssl
2. THE Project_Manager SHALL mengkonfigurasi web server untuk menggunakan SSL_Certificate pada port 443
3. THE Project_Manager SHALL menambahkan SSL_Certificate ke trusted certificate store sistem operasi
4. WHEN proyek diakses via HTTPS, THE browser SHALL tidak menampilkan warning "Not Secure"
5. THE Project_Manager SHALL menampilkan URL HTTPS (https://{project_name}.test) di dashboard TUI
6. IF mkcert tidak terinstall, THEN THE Project_Manager SHALL menampilkan instruksi instalasi mkcert kepada pengguna

### Requirement 8: Localhost Tunneling (Expose)

**User Story:** Sebagai developer, saya ingin membagikan localhost ke internet publik, sehingga saya dapat mendemonstrasikan proyek kepada client atau testing dari device lain.

#### Acceptance Criteria

1. THE Project_Manager SHALL menyediakan action "expose" pada setiap proyek di TUI
2. WHEN pengguna memilih action "expose", THE Project_Manager SHALL membuat Tunnel menggunakan service tunneling (ngrok, cloudflared, atau localtunnel)
3. THE Project_Manager SHALL menampilkan public URL yang dapat diakses dari internet dalam waktu maksimal 10 detik
4. THE Project_Manager SHALL menampilkan status Tunnel (active/inactive) dan traffic statistics di TUI
5. THE Project_Manager SHALL menyediakan action "stop expose" untuk menghentikan Tunnel
6. WHEN Tunnel dihentikan, THE Project_Manager SHALL menutup koneksi dan menghapus public URL

### Requirement 9: Dependency Tool Detection

**User Story:** Sebagai developer, saya ingin sistem mendeteksi apakah dependency tools terinstall, sehingga saya tahu tools apa yang perlu diinstall sebelum membuat proyek.

#### Acceptance Criteria

1. WHEN Lumine dijalankan pertama kali, THE Lumine SHALL memeriksa keberadaan Dependency_Tool berikut: Composer, NPM, PNPM, Yarn, Git
2. THE Lumine SHALL menampilkan status setiap Dependency_Tool (installed/not installed) dengan versi yang terdeteksi di dashboard TUI
3. WHEN Dependency_Tool tidak terinstall, THE Lumine SHALL menampilkan warning icon dan link ke dokumentasi instalasi
4. THE Lumine SHALL menyediakan action "refresh" untuk memeriksa ulang status Dependency_Tool
5. WHEN pengguna mencoba membuat proyek yang membutuhkan Dependency_Tool yang tidak terinstall, THE Project_Manager SHALL menampilkan error dan membatalkan operasi

### Requirement 10: Database Quick Actions

**User Story:** Sebagai developer, saya ingin melakukan operasi database dasar dari TUI, sehingga saya tidak perlu membuka tool database terpisah untuk operasi sederhana.

#### Acceptance Criteria

1. THE Lumine SHALL menyediakan menu "Database" di TUI dengan actions: Create Database, Drop Database, Backup Database
2. WHEN pengguna memilih "Create Database", THE Lumine SHALL menampilkan form input nama database dan membuat database baru pada Database_Connection yang aktif
3. WHEN pengguna memilih "Drop Database", THE Lumine SHALL menampilkan konfirmasi dan menghapus database yang dipilih
4. WHEN pengguna memilih "Backup Database", THE Lumine SHALL membuat SQL dump file di direktori backups dengan format {database_name}_{timestamp}.sql
5. THE Lumine SHALL menampilkan list database yang tersedia di Database_Connection yang aktif
6. IF operasi database gagal, THEN THE Lumine SHALL menampilkan error message dengan detail error dari database server

### Requirement 11: Database Connection Switcher

**User Story:** Sebagai developer, saya ingin mengganti koneksi database antara MySQL, PostgreSQL, atau SQLite, sehingga saya dapat bekerja dengan database yang berbeda sesuai kebutuhan proyek.

#### Acceptance Criteria

1. THE Lumine SHALL menyediakan menu "Switch Database Connection" di TUI dengan pilihan: MySQL, PostgreSQL, SQLite
2. WHEN pengguna memilih database type, THE Lumine SHALL menampilkan form input untuk connection details (host, port, username, password)
3. WHEN pengguna submit connection details, THE Lumine SHALL mencoba koneksi ke database dan menampilkan status koneksi (success/failed)
4. IF koneksi berhasil, THEN THE Lumine SHALL menyimpan connection details dan menggunakan koneksi tersebut untuk operasi database selanjutnya
5. THE Lumine SHALL menampilkan Database_Connection yang sedang aktif di dashboard TUI
6. THE Lumine SHALL menyimpan multiple connection profiles dan memungkinkan pengguna switch antar profiles

### Requirement 12: Real-time Database Log Monitoring

**User Story:** Sebagai developer, saya ingin melihat query log atau error log database secara real-time, sehingga saya dapat debugging masalah database dengan cepat.

#### Acceptance Criteria

1. THE Lumine SHALL menyediakan panel "Database Logs" di TUI yang menampilkan Query_Log dan error log secara real-time
2. WHILE panel "Database Logs" aktif, THE Lumine SHALL menampilkan setiap query yang dieksekusi dengan timestamp dan execution time
3. THE Lumine SHALL menyediakan filter untuk menampilkan hanya slow queries (execution time > 1 detik)
4. THE Lumine SHALL menyediakan filter untuk menampilkan hanya error logs
5. THE Lumine SHALL menyediakan action "clear logs" untuk membersihkan tampilan log
6. THE Lumine SHALL menyimpan log ke file untuk keperluan debugging lebih lanjut

### Requirement 13: TUI Dashboard Layout

**User Story:** Sebagai developer, saya ingin interface TUI yang terorganisir dengan baik, sehingga saya dapat dengan mudah melihat status service dan logs secara bersamaan.

#### Acceptance Criteria

1. THE Lumine SHALL menampilkan dashboard TUI dengan layout split panel: panel kiri untuk daftar service/proyek, panel kanan untuk log viewer
2. THE Lumine SHALL menampilkan header bar yang berisi informasi: nama aplikasi, versi, dan status Daemon_Mode
3. THE Lumine SHALL menampilkan footer bar yang berisi keyboard shortcuts yang tersedia
4. THE Lumine SHALL menyediakan tab navigation untuk berpindah antara view: Services, Projects, Database, Logs
5. THE Lumine SHALL menampilkan status indicator (icon atau warna) untuk setiap service: running (hijau), stopped (merah), error (kuning)
6. THE Lumine SHALL menyediakan search box untuk mencari service atau proyek dengan cepat

### Requirement 14: Keyboard Navigation

**User Story:** Sebagai developer, saya ingin navigasi TUI menggunakan keyboard, sehingga saya dapat bekerja lebih cepat tanpa menggunakan mouse.

#### Acceptance Criteria

1. THE Lumine SHALL mendukung Vim-style navigation keys: j (down), k (up), h (left), l (right)
2. THE Lumine SHALL mendukung arrow keys untuk navigasi alternatif
3. THE Lumine SHALL mendukung Tab key untuk berpindah antar panel atau form field
4. THE Lumine SHALL mendukung Enter key untuk memilih item atau submit form
5. THE Lumine SHALL mendukung Escape key untuk cancel action atau kembali ke menu sebelumnya
6. THE Lumine SHALL mendukung function keys: F1 (help), F2 (rename), F3 (search), F5 (refresh), F10 (quit)
7. THE Lumine SHALL menampilkan keyboard shortcuts di footer bar atau help screen (F1)

### Requirement 15: Toast Notification System

**User Story:** Sebagai developer, saya ingin melihat notifikasi sukses atau error untuk setiap action, sehingga saya mendapat feedback langsung dari operasi yang dilakukan.

#### Acceptance Criteria

1. WHEN operasi berhasil, THE Lumine SHALL menampilkan Toast_Notification dengan warna hijau dan icon checkmark selama 3 detik
2. WHEN operasi gagal, THE Lumine SHALL menampilkan Toast_Notification dengan warna merah dan icon error selama 5 detik
3. THE Lumine SHALL menampilkan Toast_Notification di posisi top-right corner dari TUI
4. THE Lumine SHALL mendukung multiple Toast_Notification yang ditampilkan secara stack (maksimal 3 notifikasi)
5. THE Lumine SHALL otomatis menghilangkan Toast_Notification setelah durasi yang ditentukan
6. THE Lumine SHALL menyediakan action untuk dismiss Toast_Notification secara manual dengan menekan 'd' key

### Requirement 16: Cross-Platform Permission Handling

**User Story:** Sebagai developer, saya ingin sistem menangani permintaan Admin_Privilege secara otomatis, sehingga saya tidak perlu manual menjalankan command dengan sudo setiap kali.

#### Acceptance Criteria

1. WHEN Lumine dijalankan pertama kali, THE Lumine SHALL memeriksa apakah memiliki Admin_Privilege untuk mengedit Hosts_File dan bind port 80/443
2. IF Lumine tidak memiliki Admin_Privilege, THEN THE Lumine SHALL menampilkan dialog yang meminta pengguna untuk restart aplikasi dengan sudo (Linux/Mac) atau Run as Administrator (Windows)
3. THE Lumine SHALL menyimpan flag bahwa Admin_Privilege sudah diberikan untuk session tersebut
4. WHEN operasi membutuhkan Admin_Privilege (edit Hosts_File, bind privileged port), THE Lumine SHALL menggunakan privilege yang sudah diberikan tanpa meminta ulang
5. IF pengguna menolak memberikan Admin_Privilege, THEN THE Lumine SHALL tetap berjalan dengan fitur terbatas (tanpa Virtual_Host otomatis dan SSL)
6. THE Lumine SHALL menampilkan warning di dashboard jika berjalan tanpa Admin_Privilege

### Requirement 17: Project Auto-Detection

**User Story:** Sebagai developer, saya ingin sistem otomatis mendeteksi proyek yang sudah ada di direktori projects, sehingga saya tidak perlu manual mendaftarkan proyek yang sudah dibuat sebelumnya.

#### Acceptance Criteria

1. WHEN Lumine dijalankan, THE Project_Manager SHALL melakukan scan pada direktori projects untuk mendeteksi proyek yang sudah ada
2. THE Project_Manager SHALL mendeteksi tipe proyek berdasarkan file marker: artisan (Laravel), wp-config.php (WordPress), package.json (Node.js), composer.json (PHP)
3. THE Project_Manager SHALL otomatis membuat Virtual_Host untuk proyek yang terdeteksi jika belum ada
4. THE Project_Manager SHALL menampilkan list proyek yang terdeteksi di dashboard TUI dengan icon yang sesuai dengan tipe proyek
5. THE Project_Manager SHALL menyediakan action "rescan projects" untuk mendeteksi ulang proyek yang baru ditambahkan secara manual
6. WHEN proyek baru terdeteksi, THE Project_Manager SHALL menampilkan Toast_Notification dengan informasi proyek yang ditemukan

### Requirement 18: Service Health Check

**User Story:** Sebagai developer, saya ingin sistem memeriksa kesehatan service secara berkala, sehingga saya dapat mengetahui jika ada service yang bermasalah sebelum mempengaruhi development.

#### Acceptance Criteria

1. WHILE service running, THE Service_Manager SHALL melakukan health check setiap 60 detik dengan mengirim request ke endpoint health check service
2. WHEN health check gagal 3 kali berturut-turut, THE Service_Manager SHALL menandai service sebagai unhealthy dan menampilkan warning di dashboard
3. THE Service_Manager SHALL menampilkan status health check di detail view service: healthy (hijau), unhealthy (kuning), down (merah)
4. THE Service_Manager SHALL mencatat hasil health check ke log file dengan timestamp
5. WHEN service unhealthy, THE Service_Manager SHALL menampilkan Toast_Notification dengan rekomendasi action (restart service, check logs)
6. THE Service_Manager SHALL menyediakan action "force health check" untuk memeriksa kesehatan service secara manual

### Requirement 19: Configuration Management

**User Story:** Sebagai developer, saya ingin menyimpan dan mengelola konfigurasi Lumine, sehingga saya dapat dengan mudah backup atau share konfigurasi ke device lain.

#### Acceptance Criteria

1. THE Lumine SHALL menyimpan konfigurasi aplikasi di file ~/.lumine/config.yaml
2. THE Lumine SHALL menyimpan konfigurasi berikut: default PHP version, default Node.js version, projects directory path, default database connection, preferred web server
3. THE Lumine SHALL menyediakan menu "Settings" di TUI untuk mengedit konfigurasi
4. WHEN konfigurasi diubah, THE Lumine SHALL memvalidasi nilai konfigurasi dan menampilkan error jika tidak valid
5. THE Lumine SHALL menyediakan action "export config" untuk membuat backup file konfigurasi
6. THE Lumine SHALL menyediakan action "import config" untuk restore konfigurasi dari backup file

### Requirement 20: Service Template Management

**User Story:** Sebagai developer, saya ingin membuat template service custom, sehingga saya dapat dengan cepat menjalankan service yang sering digunakan dengan konfigurasi yang sudah ditentukan.

#### Acceptance Criteria

1. THE Service_Manager SHALL menyediakan menu "Service Templates" di TUI untuk mengelola template service
2. THE Service_Manager SHALL menyediakan action "create template" untuk membuat template baru dengan konfigurasi: image, port, environment variables, volumes
3. THE Service_Manager SHALL menyimpan template di file ~/.lumine/templates.yaml
4. WHEN pengguna memilih "start from template", THE Service_Manager SHALL menampilkan list template yang tersedia dan membuat service baru berdasarkan template yang dipilih
5. THE Service_Manager SHALL menyediakan action "edit template" dan "delete template" untuk mengelola template yang sudah ada
6. THE Service_Manager SHALL menyediakan default templates untuk service umum: LAMP stack, LEMP stack, MEAN stack, JAMstack

