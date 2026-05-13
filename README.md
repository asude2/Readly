<br/>
<div align="center">
  <h1 align="center">Readly 📚</h1>
  <p align="center">
    Kitapseverler için geliştirilmiş sosyal ağ platformu.
  </p>
</div>

## 📖 Proje Hakkında

**Readly**, kitapseverlerin bir araya geldiği, okudukları kitapları paylaştığı, incelemeler yaptığı ve diğer okurlarla etkileşime girdiği dinamik bir sosyal platformdur. Kullanıcılar kendilerine özel bir kütüphane oluşturabilir, kitapları beğenebilir, yorum yapabilir ve ilgi alanlarına uygun diğer kullanıcıları takip edebilirler.

Ayrıca platform, kapsamlı bir **Gizlilik Mimarisi** (Açık/Gizli Profil seçenekleri) ve **Gerçek Zamanlı Mesajlaşma** (Özel ve Grup sohbetleri) özelliklerine sahiptir.

---

## ✨ Temel Özellikler

- **🔐 Kimlik Doğrulama & Güvenlik:** Güvenli kayıt olma ve giriş yapma. Şifreler `bcrypt` ile şifrelenerek veritabanında saklanır.
- **👤 Profil Yönetimi:** Profil fotoğrafı ekleme, biyografi düzenleme. Hesap gizliliği ayarları ("Herkes" veya "Gizli Profil"). Şifre ve kullanıcı adı değiştirme (Tüm sistemde otomatik senkronizasyon).
- **👥 Sosyal Etkileşim:** 
  - Açık profilleri direkt takip edebilme.
  - Gizli profillere takip isteği gönderme ve istekleri Onaylama/Reddetme sistemi.
  - Kitap paylaşımlarını beğenme ve yorum yapma.
- **💬 Gelişmiş Mesajlaşma (DM & Gruplar):** 
  - Kullanıcılar arası birebir (özel) mesajlaşma.
  - Çoklu kullanıcı ekleyerek ve çıkararak grup sohbetleri oluşturma.
  - Okunmayan mesajlar için anlık bildirim rozetleri (Polling mantığı ile gerçek zamanlı hissi).
  - Sohbet içerisinden anında kullanıcı profillerine yönlendirme.
- **🔍 Arama Sistemi:** Platformdaki diğer okurları ve kitap incelemelerini arayıp bulma.

---

## 🛠️ Kullanılan Teknolojiler

### Frontend
- **[Vue 3](https://vuejs.org/)** - Kullanıcı arayüzü kütüphanesi (Composition API).
- **[Vite](https://vitejs.dev/)** - Hızlı ve modern frontend derleyici.
- **[Tailwind CSS](https://tailwindcss.com/)** - Modern ve hızlı stil yazımı.
- **[Vue Router](https://router.vuejs.org/)** - Sayfalar arası yönlendirme (Routing).
- **FontAwesome** - İkon kütüphanesi.

### Backend
- **[Go (Golang)](https://go.dev/)** - Güçlü ve performanslı sunucu dili (`net/http` ile saf REST API mimarisi).
- **[SQLite](https://sqlite.org/)** - Hızlı ve taşınabilir ilişkisel veritabanı.
- **Bcrypt** - Kullanıcı şifreleme ve güvenlik işlemleri.

---

## 📂 Dosya ve Klasör Yapısı

```text
readly-project/
├── readly-backend/          # Go Sunucu Klasörü
│   ├── main.go              # Ana sunucu dosyası, API endpointleri ve DB konfigürasyonu
│   ├── go.mod               # Go bağımlılıkları
│   └── data.db              # SQLite Veritabanı (Uygulama çalışınca oluşur)
│
├── readly-frontend/         # Vue.js Klasörü
│   ├── index.html
│   ├── package.json         # Node.js bağımlılıkları
│   ├── tailwind.config.js   # Tailwind ayarları
│   └── src/
│       ├── main.js          # Vue uygulama başlatıcısı
│       ├── style.css        # Global CSS ve Tailwind direktifleri
│       ├── router/          # Vue Router ayarları (index.js)
│       └── components/      # Arayüz Bileşenleri
│           ├── anasayfa.vue # Ana Akış (Feed), Bildirimler
│           ├── profil.vue   # Profil, Ayarlar, Gizlilik
│           ├── mesajlar.vue # Özel/Grup Sohbetleri, Üye Yönetimi
│           ├── loginPage.vue
│           └── signupPage.vue
```

---

## 🚀 Kurulum ve Çalıştırma

Projeyi bilgisayarınızda çalıştırmak için aşağıdaki adımları izleyin. Sisteminizde **Node.js** ve **Go**'nun yüklü olduğundan emin olun.

### 1. Repoyu Klonlayın
```bash
git clone https://github.com/asude2/Readly.git
cd Readly
```

### 2. Backend'i Başlatın (Go)
Yeni bir terminal açın ve backend klasörüne gidin:
```bash
cd readly-backend
```
Bağımlılıkları yükleyin ve sunucuyu başlatın:
```bash
go mod tidy
go run main.go
```
*Sunucu varsayılan olarak `http://localhost:8000` portunda çalışacaktır. İlk çalıştırmada `data.db` adlı SQLite veritabanı otomatik olarak oluşturulur.*

### 3. Frontend'i Başlatın (Vue.js)
Farklı bir terminal sekmesi açın ve frontend klasörüne gidin:
```bash
cd readly-frontend
```
Node paketlerini yükleyin:
```bash
npm install
```
Geliştirme sunucusunu başlatın:
```bash
npm run dev
```
*Terminalde beliren adrese (genellikle `http://localhost:5173`) tarayıcınızdan giderek uygulamayı kullanmaya başlayabilirsiniz!*

---

## 💡 Ekran Görüntüleri ve Kullanım

*(Gelecekte buraya uygulamanızın anasayfa, profil ve mesajlar kısmından ekran görüntüleri ekleyebilirsiniz: `![Anasayfa](link-to-image)`)*

---

*Bu proje, kitap okuma alışkanlığını sosyal bir deneyimle birleştirmek amacıyla geliştirilmiştir. Katkı sağlamak isterseniz Pull Request göndermekten çekinmeyin!* 📚✨
