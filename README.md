📚 Readly - Kitap Paylaşım ve Sosyal Okur Ağı
Readly, kitapseverlerin okudukları kitaplar hakkında yorum yapabildiği, diğer okurları takip edebildiği ve yeni kitap keşiflerinde bulunduğu sosyal bir platformdur.


✨ Özellikler

Kullanıcı Sistemi: Kayıt olma, giriş yapma ve profil özelleştirme

Kitap Yönetimi: Kitap incelemesi ekleme, güncelleme ve silme

Sosyal Etkileşim: Diğer okurları arama, takip etme ve mesajlaşma 

Güvenlik: bcrypt ile şifre hash'leme ve güvenli kimlik doğrulama.

Görsel Profil: Base64 formatında profil fotoğrafı yükleme desteği


🚀 Teknolojiler

Bu proje modern ve yüksek performanslı bir teknoloji yığını ile geliştirilmiştir:
Katman          Teknoloji
Frontend        Vue.js 3 (Composition API)
Backend         Go (Golang)
Veritabanı      SQLite
Stil            Tailwind CSS



🛠️ Kurulum ve Çalıştırma

1. Backend Hazırlığı
# Backend dizinine gidin
cd readly-backend

# Gerekli bağımlılıkları indirin
go mod tidy

# Sunucuyu başlatın
go run main.go
# API http://localhost:8000 adresinde çalışacaktır.


2. Frontend Hazırlığı
# Ana dizine (veya frontend dizinine) gidin
npm install

# Geliştirme sunucusunu başlatın
npm run dev


## 📂 Proje Yapısı

```text
.
├── readly-backend/          # Golang & SQLite Backend
│   ├── main.go              # API Endpoint'leri ve Veritabanı Mantığı
│   └── users.db             # Uygulama Veritabanı
├── readly-frontend/         # Vue.js Frontend
│   ├── public/              # Statik Varlıklar
│   ├── src/
│   │   ├── assets/          # CSS ve Görseller
│   │   ├── components/      # Uygulama Sayfaları ve Bileşenleri
│   │   │   ├── anasayfa.vue
│   │   │   ├── loginPage.vue
│   │   │   ├── profil.vue
│   │   │   └── signupPage.vue
│   │   ├── router/          # Sayfa Yönlendirmeleri
│   │   ├── App.vue          # Kök Bileşen
│   │   └── main.js          # Başlangıç Dosyası
│   ├── index.html           # Ana HTML Şablonu
│   └── package.json         # Bağımlılıklar ve Scriptler
└── README.md                # Proje Dokümantasyonu