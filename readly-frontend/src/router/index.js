import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/loginPage.vue'
import Register from '../components/signupPage.vue'
import Anasayfa from '../components/anasayfa.vue'
import Profil from '../components/profil.vue'



const routes = [
  {path: '/', name: 'Login', component: Login },
  {path: '/register', name: 'Register', component: Register },
  {path: '/anasayfa', name: 'Anasayfa', component:Anasayfa},
  {path: '/profil/:username', name:'profil', component: Profil, props:true},
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router


