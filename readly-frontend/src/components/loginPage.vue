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

      <div class="book-content max-w-md mx-auto p-10 border rounded-[2px] shadow-lg bg-stone-200">
        <h1 class="text-black text-2xl p-5">
          <i class="fa-solid fa-arrow-left pr-5 cursor-pointer" @click="logoutApp"></i>
          Giriş Yap
        </h1>

        <div v-if="message" class="p-4 mb-4 text-sm rounded-lg" :class="isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'">
          {{ message }}
        </div>

        <form @submit.prevent="handleLogin" class="flex flex-col">
          <div class="p-4 flex flex-col">
            <p class="text-black">Kullanıcı Adınız</p>
            <input v-model="form.username" type="text" placeholder="kullanıcı adı" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
          </div>
          <div class="p-4 flex flex-col">
            <p class="text-black">Şifreniz</p>
            <input v-model="form.password" type="password" placeholder="şifre" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
          </div>
          <div class=" flex justify-end gap-2">
            <button type="button" class=" text-black p-2 rounded-[20px] border border-black hover:bg-white hover:border-white transition-colors duration-300 ease-in-out text-sm">Şifremi Unuttum</button>
            <button type="submit" :disabled="isLoading" class="text-red-600 border border-red-600 bg-red-200 font-medium p-2 pl-6 pr-6 rounded-[20px] bg-transparent hover:bg-[#b91c1c] hover:text-white transition-colors duration-300 ease-in-out">
              {{ isLoading ? 'Giriş Yapılıyor...' : 'Giriş Yap' }}
            </button>
          </div>
          <div class="text-center">
            <button @click.prevent="goToRegister" class="m-5 pt-5 text-black font-bold rounded-[10px] hover:text-[#b91c1c] transition transition-colors duration-300 ease-in-out text-sm">HESABINIZ YOK MU? HEMEN KAYDOLUN</button>
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
  align-items: start; /* Biraz yukarıda başlasın, mt-10'u korumak için */
  perspective: 2000px; /* Derinlik hissi (çok önemli) */
  overflow: hidden; /* Taşmaları gizle */
}

/* Kitabın kendisi ve kapakların sarmalayıcısı */
.book-wrap {
  position: relative;
  width: 100%;
  max-width: 448px; /* max-w-md ile aynı */
  margin-top: 2.5rem; /* mt-10 ile aynı */
  transition: transform 1s ease-in-out;
  transform-style: preserve-3d; /* Çocukların 3D boşlukta kalmasını sağlar */
}

/* Kapakların ortak özellikleri */
[class^="book-cover-"] {
  position: absolute;
  top: 0;
  width: 50%; /* Her kapak formun yarısını kaplar */
  height: 100%;
  background: #b91c1c; /* Daha koyu bir kırmızı (kitap kapağı rengi) */
  transition: transform 1.5s cubic-bezier(0.4, 0, 0.2, 1);
  transform-style: preserve-3d;
  z-index: 10; /* Formun üzerinde durmaları için */
  border: 1px solid #ffffff;
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
  background: #d8d8d8; /* Formun rengine yakın, açılan iç kısım */
  transform: rotateY(180deg); /* Arka yüzü 180 derece döndür */
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


