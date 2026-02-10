<template>
    <div class="max-w-md mx-auto mt-10 p-10 border rounded-lg shadow-lg bg-red-200">
        <h1 class="text-white text-2xl p-5"><i class="fa-solid fa-arrow-left pr-5 cursor-pointer" @click="logoutApp"></i>Giriş Yap</h1>

        <div v-if="message" class="p-4 mb-4 text-sm rounded-lg" :class="isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'">
            {{ message }}
        </div>

        <form @submit.prevent="handleLogin" class="flex flex-col">
            <div class="p-4 flex flex-col">
                <p class="text-white">Kullanıcı Adınız</p>
                <input v-model="form.username" type="text" placeholder="kullanıcı adı" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
            </div>
            <div class="p-4 flex flex-col">
                <p class="text-white">Şifreniz</p>
                <input v-model="form.password" type="password" placeholder="şifre" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
            </div>
            <div class=" flex justify-end">
                <button type="button" class=" text-white p-2 mr-3 rounded-[20px] border border-white hover:bg-white hover:text-red-600 transition-colors duration-300 ease-in-out">Şifremi Unuttum</button>
                <button type="submit" :disabled="isLoading" class="text-red-600 border border-red-600 bg-red-200 font-medium p-2 pl-6 pr-6 rounded-[20px] bg-transparent hover:bg-red-600 hover:text-white transition-colors duration-300 ease-in-out">
                  {{ isLoading ? 'Giriş Yapılıyor...' : 'Giriş Yap' }}
                </button>
            </div>
            <div>
                <button @click.prevent="goToRegister" class="m-5 pt-5  text-white font-bold rounded-[10px] hover:text-red-500 transition transition-colors duration-300 ease-in-out">HESABINIZ YOK MU? HEMEN KAYDOLUN</button>
            </div>
        </form>
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
            isError: false
        }
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


