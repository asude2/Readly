<template>
  <div class="book-scene">

    <div class="book-wrap" :class="{ 'is-open': bookOpened }">

      <div class="book-cover-left">
        <div class="cover-front text-white text-5xl text-right p-5 mt-20 tracking-widest font-['Playfair_Display']">REA</div>
        <div class="cover-back"></div>
      </div>

      <div class="book-cover-right">
        <div class="cover-front text-white text-left text-5xl p-5 mt-20 tracking-widest font-['Playfair_Display']">DLY</div>
        <div class="cover-back"></div>
      </div>

      <div class="book-content w-full max-w-lg mx-auto p-10 border border-gray-200 dark:border-slate-700 rounded-xl shadow-2xl bg-white dark:bg-slate-800 transition-colors duration-300">
        <h1 class="text-black dark:text-white text-3xl font-bold mb-8 flex items-center">
          <i class="fa-solid fa-arrow-left pr-4 cursor-pointer hover:text-red-600 transition-colors" @click="logoutApp"></i>
          Giriş Yap
        </h1>

        <div v-if="message" class="p-4 mb-6 text-sm rounded-lg border" :class="isError ? 'bg-red-50 border-red-200 text-red-700 dark:bg-red-900/20 dark:border-red-800 dark:text-red-400' : 'bg-green-50 border-green-200 text-green-700 dark:bg-green-900/20 dark:border-green-800 dark:text-green-400'">
          {{ message }}
        </div>

        <form @submit.prevent="handleLogin" class="flex flex-col space-y-5">
          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5">Kullanıcı Adınız</label>
            <input v-model="form.username" type="text" placeholder="Kullanıcı adı" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2.5 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200" required>
          </div>
          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5">Şifreniz</label>
            <input v-model="form.password" type="password" placeholder="Şifre" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2.5 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200" required>
          </div>
          <div class="flex items-center justify-between pt-2">
            <button type="button" class="text-sm font-medium text-gray-600 dark:text-slate-400 hover:text-red-600 dark:hover:text-red-400 transition-colors duration-200">Şifremi Unuttum</button>
            <button type="submit" :disabled="isLoading" class="bg-red-600 hover:bg-red-700 text-white font-medium py-2.5 px-8 rounded-lg shadow-md hover:shadow-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed">
              {{ isLoading ? 'Giriş Yapılıyor...' : 'Giriş Yap' }}
            </button>
          </div>
          <div class="text-center mt-8 pt-4 border-t border-gray-200 dark:border-slate-700">
            <button @click.prevent="goToRegister" class="text-sm font-semibold text-gray-600 dark:text-slate-400 hover:text-red-600 dark:hover:text-red-500 transition-colors duration-300">Hesabınız yok mu? Hemen kaydolun</button>
          </div>
        </form>
      </div>

    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            form: {
                username: '',
                password: ''
            },
            isLoading: false,
            message: '',
            isError: false,
            bookOpened: false //animasyonu başlatmak için state
        }
    },

    mounted() {
        // sayfa yüklendikten kısa bir süre sonra kapağı açmak için
        setTimeout(() => {
            this.bookOpened = true;
        }, 500); // Yarım saniye bekle ve aç
    },

    methods: {
        goToRegister() {
            this.$router.push('/register')
        },
        logoutApp(){
            window.location.href = "https://www.google.com"
        },
        async handleLogin() {
            this.isLoading = true;
            this.message = '';

            try {
                // Yeni /api/login endpoint'ine istek atıyoruz
                const response = await axios.post('http://localhost:8000/api/login', this.form);

                this.isError = false;
                this.message = response.data.message; // "Giriş başarılı!"

                localStorage.setItem('username', this.form.username);

                // Giriş başarılı olduğunda kullanıcıyı ana sayfaya yönlendiriyorum.
                this.$router.push('/anasayfa');

            } catch (error) {
                this.isError = true;
                if (error.response && error.response.data.error) {
                    this.message = error.response.data.error; // "Kullanıcı adı veya şifre hatalı"
                } else {
                    this.message = 'Sunucuya bağlanılamadı.';
                }
            } finally {
                this.isLoading = false;
            }
        }
    }
}

</script>

<style scoped>
/* --- CSS 3D Kitap Animasyonu Stilleri --- */

/* Tüm sahneyi kapsayan alan */
.book-scene {
  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center; /* Dikey ve yatay olarak ortala */
  perspective: 2000px; /* Derinlik hissi (çok önemli) */
  overflow: hidden; /* Taşmaları gizle */
}

/* Kitabın kendisi ve kapakların sarmalayıcısı */
.book-wrap {
  position: relative;
  width: 100%;
  max-width: 512px; /* max-w-lg ile aynı (büyütüldü) */
  margin-top: 0; /* Ortaladığımız için margin-top kaldırıldı */
  transition: transform 1s ease-in-out;
  transform-style: preserve-3d; /* Çocukların 3D boşlukta kalmasını sağlar */
}

/* Kapakların ortak özellikleri */
[class^="book-cover-"] {
  position: absolute;
  top: 0;
  width: 50%; /* Her kapak formun yarısını kaplar */
  height: 100%;
  background: #b91c1c; /* Koyu kırmızı (kitap kapağı rengi) */
  transition: transform 1.5s cubic-bezier(0.4, 0, 0.2, 1);
  transform-style: preserve-3d;
  z-index: 10; /* Formun üzerinde durmaları için */
  border: 1px solid #7f1d1d;
  box-shadow: inset 0 0 20px rgba(0,0,0,0.3);
}

/* Sol Kapak */
.book-cover-left {
  left: 0;
  transform-origin: left center; /* Soldan dışarı doğru dönecek */
  border-radius: 8px 0 0 8px; /* Sadece sol köşeler yuvarlak */
}

/* Sağ Kapak */
.book-cover-right {
  right: 0;
  transform-origin: right center; /* Sağdan dışarı doğru dönecek */
  border-radius: 0 8px 8px 0; /* Sadece sağ köşeler yuvarlak */
}

/* Kapakların ön ve arka yüzleri (3D efekt için gerekli) */
.cover-front, .cover-back {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden; /* Arka yüzü görmeyelim */
}

.cover-back {
  background: #e2e8f0; /* Formun rengine yakın, açılan iç kısım */
  transform: rotateY(180deg); /* Arka yüzü 180 derece döndür */
}
:global(.dark) .cover-back {
  background: #1e293b;
}

/* --- Animasyon Tetikleyicileri (Vue state'ine bağlı) --- */

/* Kitap açıldığında sol kapağı 120 derece sola döndür */
.is-open .book-cover-left {
  transform: rotateY(-140deg);
}

/* Kitap açıldığında sağ kapağı 120 derece sağa döndür */
.is-open .book-cover-right {
  transform: rotateY(140deg);
}

/* Kapaklar açılırken formu hafifçe öne getir */
.is-open .book-content {
  transform: translateZ(0px); /* 3D derinliğini sıfırla */
  opacity: 1;
}

/* Kitap İçeriği (Login Formunuz) */
.book-content {
  position: relative;
  z-index: 1; /* Kapakların arkasında kalması için */
  opacity: 0; /* Başlangıçta gizli */
  transform: translateZ(-20px); /* Başlangıçta hafifçe geride */
  transition: opacity 1s ease 0.5s, transform 1s ease 0.5s; /* Kapaklar açılmaya başladıktan sonra belirsin */
}
</style>


