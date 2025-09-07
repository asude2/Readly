<template>
    <div class="max-w-md mx-auto mt-10 p-6 border rounded-lg shadow-lg">
        <h1 class="text-2xl p-5"><i class="fa-solid fa-arrow-left pr-2"></i>Giriş Yap</h1>

        <div v-if="message" class="p-4 mb-4 text-sm rounded-lg" :class="isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'">
            {{ message }}
        </div>

        <form @submit.prevent="handleLogin" class="flex flex-col">
            <div class="p-4 flex flex-col">
                <p>Kullanıcı Adınız</p>
                <input v-model="form.username" type="text" placeholder="kullanıcı adı" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
            </div>
            <div class="p-4 flex flex-col">
                <p>Şifreniz</p>
                <input v-model="form.password" type="password" placeholder="şifre" class=" border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-400" required>
            </div>
            <div class=" flex justify-end">
                <button type="button" class=" text-blue-600 p-2 mr-3 rounded-[20px] hover:bg-gray-100 transition">Şifremi Unuttum</button>
                <button type="submit" :disabled="isLoading" class="bg-gray-200 font-medium p-2 pl-6 pr-6 rounded-[20px] hover:bg-gray-100 transition disabled:bg-gray-400">
                  {{ isLoading ? 'Giriş Yapılıyor...' : 'Giriş Yap' }}
                </button>
            </div>
            <div>
                <button @click.prevent="goToRegister" class="m-6 px-3 py-2 text-blue-600 font-bold rounded-[10px] hover:bg-gray-100 transition">HESABINIZ YOK MU? HEMEN KAYDOLUN</button>
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


