<template>
  <div class="flex justify-between">
        <div class="pl-7 pt-5">
          <h1 class="text-red-500 text-3xl font-medium italic font-serif flex  cursor-default">
              Readly
              <i class="fa-solid fa-book-open text-2xl px-2 hover:scale-110 transition-transform transition-colors duration-700 border-2 border-transparent rounded-full"></i>
          </h1>
          <i @click="goLoginPg" class="fa-regular fa-circle-left text-4xl  cursor-pointer"></i>
        </div>

        <div class="pt-5">
            <nav class="w-full">
                <ul class="flex items-center tracking-widest justify-center space-x-12">
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      kitaplar
                    </a>
                  </li>
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      okurlar
                    </a>
                  </li>
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      yorumlamalar
                    </a>
                  </li>
                </ul>
            </nav>
            <div class="flex items-center border border-gray-300 rounded-full px-3 my-4 w-full">
                <i class="fa-solid fa-magnifying-glass text-gray-500"></i>
                <input v-model="searchQuery" type="text" placeholder="arama yap.." class="w-full px-2 py-1 outline-none rounded-full"/>
                <button @click="searchUsers" class="bg-gray-200 rounded-[20px] hover:bg-gray-300 px-3 cursor-pointer">ARA</button>
            </div>
            <div v-for="user in users" :key="user.username" class="flex bg-gray-200 justify-between items-center border rounded-[20px] py-2 m-2 cursor-pointer border-gray-300 ">
                <p class="font-semibold mx-6">{{user.username}} {{ user.firstName }}{{ user.lastName }}</p>
                <button @click="follow(user.username, user.isFollowing)" class="bg-black text-white border rounded-[10px] px-2 py-[2px] font-light text-[15px] hover:bg-gray-700 mx-6">{{ user.isFollowing ? 'Takipten Çık' : 'Takip Et' }}</button>
            </div>
        </div>

        <div @click.prevent="goToProfile" class="pt-3">
          <i class="fa-solid fa-circle-user p-3 text-3xl flex cursor-pointer"><p class="px-1 text-2xl">profil</p></i>
        </div>

  </div>

</template>

<script setup>
import { ref, onMounted} from 'vue'
import { useRouter } from 'vue-router'

const router=useRouter()
const goToProfile=()=>{
  router.push('/profil')
}
const goLoginPg=()=>{
  router.push('/')
}

const firstname = ref('')
const lastname = ref('')
const username = ref('')
const users = ref([])
const searchQuery=ref('')




//profil fetchleme
const fetchProfile = async () => {
  username.value = localStorage.getItem('username') || ''
  if (!username.value) {
  console.error('Username bulunamadı')
  return
  }
  try {
    const res = await fetch(`http://localhost:8000/api/profile?username=${username.value}`) //backenddeki profile endpointine get isteği atılıyor
    const data = await res.json() //backend bu username'e göre kullanıcı bilgilerini bulup jsona döndürüyor. Datada kullanıcı bilgileri var.
    if (res.ok) {
      firstname.value=data.firstname || ""
      lastname.value=data.lastname || ""
    } else { console.error('Profil Bulunamadı:',data.error) }
  } catch (err) { console.error('Profil fetch hatası:', err) }
}

//kullanıcı arama
const searchUsers=async()=>{
  if(!searchQuery.value)return
  try{
    const res = await fetch(`http://localhost:8000/api/findUsers?q=${searchQuery.value}&current=${username.value}`)
    const data=await res.json()
    if(res.ok){
      users.value=data
    }
    else console.error('Kullanıcı bulunamadı:',data.error)
  }
  catch (err){console.error('Kullanıcı arama hatası:', err)}
}

// takip etme / takipten çıkma
const follow = async (targetUsername, isFollowing) => {
  const action = isFollowing ? 'unfollow' : 'follow'
  try {
    const res = await fetch('http://localhost:8000/api/followUser', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        follower: username.value,
        following: targetUsername,
        action
      })
    })
    const data = await res.json()
    if (res.ok) {
      searchUsers()
    } else console.error('Takip hatası:', data.error)
  } catch (err) { console.error('Takip hatası:', err) }
}














onMounted(() => {
  fetchProfile()
})

</script>

