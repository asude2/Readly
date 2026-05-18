<template>
  <div class="min-h-screen w-full flex items-center justify-center p-4 bg-gray-50 dark:bg-slate-900 transition-colors duration-300">
    <div class="w-full max-w-xl p-10 border border-gray-200 dark:border-slate-700 rounded-xl shadow-2xl bg-white dark:bg-slate-800 transition-colors duration-300">
      <h1 class="text-black dark:text-white text-3xl font-bold mb-8 flex items-center">
        <i class="fa-solid fa-arrow-left pr-4 cursor-pointer hover:text-red-600 transition-colors" @click="goToLogin"></i> Kayıt Ol
      </h1>

      <form @submit.prevent="submitRegistration" class="flex flex-col space-y-4">
        <!-- Ad ve Soyad Yan Yana -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="firstname">Ad</label>
            <input v-model="form.firstname" id="firstname" type="text" placeholder="Adınız" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
          </div>

          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="lastname">Soyad</label>
            <input v-model="form.lastname" id="lastname" type="text" placeholder="Soyadınız" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
          </div>
        </div>

        <div class="flex flex-col">
          <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="email">E-posta</label>
          <input v-model="form.email" id="email" type="email" placeholder="E-posta" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
        </div>

        <div class="flex flex-col">
          <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="username">Kullanıcı Adı</label>
          <input v-model="form.username" id="username" type="text" placeholder="Kullanıcı Adınız" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
        </div>

        <!-- Şifre ve Şifre Tekrar Yan Yana -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="password">Şifre</label>
            <input v-model="form.password" id="password" type="password" placeholder="Şifre" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
          </div>

          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="password2">Şifre Tekrar</label>
            <input v-model="form.password2" id="password2" type="password" placeholder="Şifre Tekrar" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
          </div>
        </div>

        <!-- Doğum Tarihi ve Cinsiyet Yan Yana -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="dob">Doğum Tarihi</label>
            <input v-model="form.dob" id="dob" type="date" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
          </div>

          <div class="flex flex-col">
            <label class="block text-sm font-medium text-gray-800 dark:text-slate-300 mb-1.5" for="gender">Cinsiyet</label>
            <select v-model="form.gender" id="gender" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-4 py-2 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-all duration-200">
              <option value="">Seçiniz</option>
              <option value="male">Erkek</option>
              <option value="female">Kadın</option>
              <option value="other">Diğer</option>
            </select>
          </div>
        </div>

        <div v-if="message" class="p-4 mb-2 text-sm rounded-lg border" :class="isError ? 'bg-red-50 border-red-200 text-red-700 dark:bg-red-900/20 dark:border-red-800 dark:text-red-400' : 'bg-green-50 border-green-200 text-green-700 dark:bg-green-900/20 dark:border-green-800 dark:text-green-400'">
          {{ message }}
        </div>

        <div class="flex justify-end pt-2">
          <button type="submit" :disabled="isLoading" class="bg-red-600 hover:bg-red-700 text-white font-medium py-2.5 px-8 rounded-lg shadow-md hover:shadow-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed">
            {{ isLoading ? 'Kaydediliyor...' : 'Kayıt Ol' }}
          </button>
        </div>

        <div class="text-center mt-6 pt-4 border-t border-gray-200 dark:border-slate-700">
          <p class="text-sm font-semibold text-gray-600 dark:text-slate-400 cursor-pointer hover:text-red-600 dark:hover:text-red-500 transition-colors duration-300" @click.prevent="goToLogin">
            Zaten Hesabınız Var mı? Giriş Yap
          </p>
        </div>
      </form>
    </div>
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
