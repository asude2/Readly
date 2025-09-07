<template>
  <div class="max-w-md mx-auto mt-10 p-6 border rounded-lg shadow-lg">
    <h1 class="text-2xl p-5 flex items-center">
      <i class="fa-solid fa-arrow-left pr-2"></i> Kayıt Ol
    </h1>

    <form @submit.prevent="submitRegistration" class="flex flex-col space-y-4">
      <div class="flex flex-col">
        <label for="firstname">Ad</label>
        <input v-model="form.firstname" id="firstname" type="text" placeholder="Adınız" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="lastname">Soyad</label>
        <input v-model="form.lastname" id="lastname" type="text" placeholder="Soyadınız" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="email">E-posta</label>
        <input v-model="form.email" id="email" type="email" placeholder="E-posta" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="username">Kullanıcı Adı</label>
        <input v-model="form.username" id="username" type="text" placeholder="Kullanıcı Adınız" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="password">Şifre</label>
        <input v-model="form.password" id="password" type="password" placeholder="Şifre" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="password2">Şifre Tekrar</label>
        <input v-model="form.password2" id="password2" type="password" placeholder="Şifre Tekrar" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="dob">Doğum Tarihi</label>
        <input v-model="form.dob" id="dob" type="date" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
      </div>

      <div class="flex flex-col">
        <label for="gender">Cinsiyet</label>
        <select v-model="form.gender" id="gender" class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400">
          <option value="">Seçiniz</option>
          <option value="male">Erkek</option>
          <option value="female">Kadın</option>
          <option value="other">Diğer</option>
        </select>
      </div>

      <div v-if="message" class="p-4 mb-4 text-sm rounded-lg" :class="isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'">
        {{ message }}
      </div>

      <div class="flex justify-end gap-2">
        <button type="submit" :disabled="isLoading" class="bg-gray-700 text-white font-medium px-6 py-2 rounded-[20px] hover:bg-gray-600 transition-colors disabled:bg-gray-400">
          {{ isLoading ? 'Kaydediliyor...' : 'Kayıt Ol' }}
        </button>
      </div>

      <p class="text-blue-600 font-bold text-center mt-4 cursor-pointer hover:bg-gray-100 transition-colors px-3 py-2 rounded-[10px]" @click.prevent="goToLogin">
        Zaten Hesabınız Var mı? Giriş Yap
      </p>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      form: {
        firstname: '',
        lastname: '',
        email: '',
        username: '',
        password: '',
        password2: '',
        dob: '',
        gender: ''
      },
      isLoading: false,
      message: '',
      isError: false
    };
  },
  methods: {
    goToLogin() {
      this.$router.push('/');
    },
    async submitRegistration() {
      if (this.form.password !== this.form.password2) {
        this.isError=true;
        this.message='Şifreler uyuşmuyor!';
        return;
      }
      this.isLoading = true;
      this.message = '';

      try {
        const payload = {
          firstname: this.form.firstname,
          lastname: this.form.lastname,
          email: this.form.email,
          username: this.form.username,
          password: this.form.password,
          dob: this.form.dob,
          gender: this.form.gender,
        };

        const response = await axios.post('http://localhost:8000/api/register', payload);

        // Başarılı olursa
        this.isError = false;
        this.message = response.data.message;
        // Formu sıfırla
        Object.keys(this.form).forEach(key => this.form[key] = '');

      } catch (error) {
        // Hata olursa
        this.isError = true;
        if (error.response && error.response.data.error) {
          this.message = error.response.data.error;
        } else {
          this.message = 'Bir sunucu hatası oluştu.';
        }
      } finally {
        this.isLoading = false;
      }
    }
  }
}
</script>
