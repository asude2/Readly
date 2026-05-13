<template>
  <div>
      <!--------------------------------------------HEADER KISMI------------------------------------------------------------------------>
        <div class="flex justify-between">
          <!-- ------------logo kısmı-------------- -->
          <div class="pl-7 pt-5">
            <h1 class="text-[#b91c1c] text-3xl font-medium italic font-serif flex  cursor-default">
                Readly
                <i class="fa-solid fa-book-open text-2xl px-2 hover:scale-110 transition-transform transition-colors duration-700 border-2 border-transparent rounded-full"></i>
            </h1>
            <i @click="goLoginPg" class="fa-regular fa-circle-left text-4xl  cursor-pointer"></i>
          </div>










          <div class="pt-5 max-w-2xl mx-auto">
            <nav class="w-full mb-6">
              <ul class="text-red-600 flex items-center tracking-widest justify-center space-x-8 md:space-x-12 font-medium text-sm">
                <li>
                  <button @click="kitapAra" class="relative pb-1 after:block after:h-[2px] after:w-0 after:bg-red-600 after:transition-all after:duration-300 after:absolute after:bottom-0 after:left-0 hover:after:w-full" :class="{'after:w-full': kitapGör}">
                    kitap incelemeleri
                  </button>
                </li>
                <li>
                  <button @click="okurAra" class="relative pb-1 after:block after:h-[2px] after:w-0 after:bg-red-600 after:transition-all after:duration-300 after:absolute after:bottom-0 after:left-0 hover:after:w-full" :class="{'after:w-full': okurGör}">
                    okurlar
                  </button>
                </li>
                <li>
                  <button class="relative pb-1 after:block after:h-[2px] after:w-0 after:bg-red-600 after:transition-all after:duration-300 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                    kitaplar
                  </button>
                </li>
              </ul>
            </nav>

            <!-- Okur Arama Çubuğu -->
            <div v-if="okurGör" class="transition-all duration-300">
              <div class="flex items-center bg-white/10 border border-red-500/50 rounded-full px-4 py-1.5 my-4 w-full focus-within:border-red-500 focus-within:ring-1 focus-within:ring-red-500/30">
                <i class="fa-solid fa-magnifying-glass text-gray-400 mr-2"></i>
                <input
                  @keyup.enter="searchUsers"
                  v-model="searchQuery"
                  type="text"
                  placeholder="kullanıcı ara.."
                  class="w-full bg-transparent text-black placeholder-gray-500 outline-none text-sm"
                />
                <button @click="searchUsers" class="bg-red-600 text-white rounded-full hover:bg-red-700 px-4 py-1 text-xs font-bold transition-colors">
                  ARA
                </button>
              </div>

              <!-- Kullanıcı Listesi -->
              <div v-for="user in users.filter(u=>u.username!==username)" :key="user.username"
                  class="flex bg-white/5 hover:bg-white/10 justify-between items-center border border-gray-700 rounded-xl py-3 px-4 mb-2 cursor-pointer transition-all">
                <div @click="goUsersProfile(user.username)" class="flex items-center font-semibold text-white">
                  <img v-if="user.photo" :src="user.photo" alt="profile photo" class="w-10 h-10 rounded-full mr-3 object-cover border border-gray-600">
                  <i v-else class="fa-solid fa-circle-user text-3xl mr-3 text-gray-400"></i>
                  <span>{{ user.username }}</span>
                </div>
                <button
                  @click="follow(user.username, user.isFollowing, user.isPending)"
                  class="border rounded-full px-4 py-1 text-xs font-semibold transition-all"
                  :class="user.isFollowing ? 'bg-transparent border-gray-500 text-gray-300 hover:bg-red-900/20' : user.isPending ? 'bg-gray-200 text-gray-700 hover:bg-gray-300' : 'bg-white text-black hover:bg-gray-200'"
                >
                  {{ user.isFollowing ? 'Takipten Çık' : user.isPending ? 'İstek Gönderildi' : 'Takip Et' }}
                </button>
              </div>
            </div>

            <!-- Kitap Yorumu Arama Çubuğu -->
            <div v-if="kitapGör" class="flex items-center bg-white/10 border border-red-500/50 rounded-full px-4 py-2 my-4 w-full focus-within:border-red-500 focus-within:ring-1 focus-within:ring-red-500/30 transition-all">
              <i class="fa-solid fa-magnifying-glass text-red-600 mr-3"></i>
              <input
                v-model="searchBookQuery"
                type="text"
                placeholder="kitap incelemesi ara.."
                class="w-full bg-transparent text-black placeholder-gray-500 outline-none text-sm"
              />
            </div>
          </div>



          <div class="flex items-center mt-11 pr-5 gap-2">
            <!-- ! messages -->
            <button @click="$router.push('/mesajlar')" class="relative p-2 text-gray-600 hover:bg-gray-100 rounded-full transition-colors mr-2">
              <i class="fa-regular fa-paper-plane text-xl"></i>
              <span v-if="unreadMessageCount > 0" class="absolute top-0 right-0 flex items-center justify-center h-4 w-4 bg-red-600 text-white text-[10px] rounded-full border-2 border-white font-bold">
                {{ unreadMessageCount > 9 ? '9+' : unreadMessageCount }}
              </span>
            </button>

            <!-- ! notifications -->
            <div class="relative notifications-wrapper">
              <button @click="toggleNotifications" class="relative p-2 text-gray-600 hover:bg-gray-100 rounded-full transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
              </svg>
              <span v-if="pendingRequests.length" class="absolute top-1 right-1 flex h-3 w-3">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-3 w-3 bg-red-600"></span>
              </span>
            </button>

            <div v-if="isNotificationsOpen"
                class="absolute right-0 mt-2 w-80 z-50 rounded-2xl bg-white p-2 shadow-xl border border-gray-100 ring-1 ring-black ring-opacity-5">

              <div class="p-3 border-b border-gray-50 mb-2">
                <p class="text-sm font-bold text-black">Bildirimler</p>
              </div>

              <div class="max-h-20 overflow-y-auto">
                <div v-if="pendingRequests.length || acceptedFollows.length">
                  <div v-for="item in pendingRequests" :key="`pending-${item.follower}`"
                      class="flex items-center justify-between gap-3 rounded-xl hover:bg-gray-50 p-3 mb-1 transition-colors">
                    <div class="text-xs text-gray-800">
                      <span class="font-bold">{{ item.follower }}</span> seni takip etmek istiyor.
                    </div>
                    <div class="flex gap-2">
                      <button @click="approveRequest(item.follower)"
                              class="px-3 py-1.5 rounded-lg bg-red-600 text-white hover:bg-red-700 text-[10px] font-semibold">
                        Onayla
                      </button>
                      <button @click="rejectRequest(item.follower)"
                              class="px-3 py-1.5 rounded-lg bg-gray-200 text-gray-800 hover:bg-gray-300 text-[10px] font-semibold transition-colors">
                        Sil
                      </button>
                    </div>
                  </div>

                  <div v-for="item in acceptedFollows" :key="`accepted-${item.follower}`"
                      class="rounded-xl hover:bg-gray-50 p-3 mb-1 text-xs text-gray-600 transition-colors">
                    <span class="font-bold text-black">{{ item.follower }}</span> seni takip etti.
                  </div>
                </div>

                <div v-else class="py-8 text-center">
                  <p class="text-sm text-gray-400">Henüz yeni bir bildirim yok.</p>
                </div>
              </div>
            </div>
          </div>
          </div>




          <!-- -------------profile gitme kısmı------------ -->
          <div @click.prevent="goToProfile" class="m-6 cursor-pointer flex items-center space-x-2">
            <p class="text-xl font-semibold">Profile</p>
            <img v-if="profile" :src="profile" alt="Profil" class="w-10 h-10 rounded-full object-cover border-2 border-red-500 hover:scale-110 transition transition-transform duration-300">
          </div>


        </div>

        <!--------------------------------------------CONTENT KISMI-------------------------------------------------------------------->
        <div class="flex-column justify-center items-start min-h-screen m-40 mt-20">
          <div v-for="book in filteredBooks" :key="book.id" class="rounded-2xl shadow-lg bg-red-50 p-6 transition-all duration-300 w-[90%] relative justify-center m-10 ">
            <p @click="goUsersProfile(book.username)" class="flex items-center cursor-pointer font-semibold text-l ml-5 mb-8 mt-2  px-4 py-1 inline-flex  bg-red-50 hover:bg-[#b91c1c] hover:text-white rounded-[10px] transition-colors">
              <img v-if="book.user_photo" :src="book.user_photo" class="w-8 h-8 rounded-full mr-2 object-cover border border-white">
              <i v-else class="fa-solid fa-circle-user text-sm mr-2"></i>{{ book.username }}
            </p>
            <img class="w-[20%] h-40 m-5 object-cover  rounded-[5px]" alt="book image" :src="book.image" />
            <div class="px-4 py-4 overflow-hidden">
              <h2 class="text-xl font-bold mb-2">{{book.title}}</h2>
              <p class="text-gray-600 text-sm mb-4">{{book.description}}</p>

              <div class="flex items-center gap-2 mt-2">
                <i @click="toggleLike(book)" class="cursor-pointer text-xl transition-all duration-200" :class="book.is_liked ? 'fa-solid fa-heart text-red-500' : 'fa-regular fa-heart text-black hover:text-red-500'"></i>
                <span class="text-sm font-medium">{{ book.like_count }}</span>
              </div>

              <div class="mt-4">
                <div @click="toggleComments(book)" class="flex items-center cursor-pointer hover:text-red-500 transition-colors">
                  <i class="fa-regular fa-comment mr-2"></i>
                  <span>Yorumlar ({{ book.comment_count || 0 }})</span>
                </div>

                <div v-if="book.showComments" class="mt-4 bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
                  <div v-for="c in book.comments" :key="c.id" class="mb-3 text-sm border-b pb-2">
                    <div class="flex justify-between items-start">
                      <div class="flex-1">
                        <p class="font-bold text-red-600">{{ c.username }}</p>
                        <p class="text-gray-700">{{ c.content }}</p>
                      </div>
                      <button
                        v-if="c.username === username"
                        @click="deleteComment(book, c.id)"
                        class="ml-2 text-red-500 hover:text-red-700 text-xl mt-1"
                      >
                        <i class="fa-solid fa-trash"></i>
                      </button>
                    </div>
                  </div>
                </div>

                <div class="flex mt-3">
                  <input v-model="book.newComment" @keyup.enter="submitComment(book)" placeholder="Yorumunu yaz..." class="flex-1 border rounded-lg px-2 py-1 text-sm outline-none">
                  <button @click="submitComment(book)" class="ml-2 bg-red-500 text-white px-3 py-1 rounded-lg text-sm hover:bg-red-600">Gönder</button>
                </div>
              </div>

            </div>
          </div>
        </div>


  </div>

</template>

<script setup>
import { comment } from 'postcss'
import { ref, onMounted, onUnmounted, computed} from 'vue'
import { useRouter } from 'vue-router'

const router=useRouter()
const isNotificationsOpen = ref(false);


const goLoginPg=()=>{
  router.push('/')
}

const goToProfile=()=>{
  router.push(`/profil/${username.value}`)
}

//bunu profil.vue ya yönlendirdik gibi bir şey.. route'a /profil/asuderk gibi gidiyor
const goUsersProfile=(user) => {
  router.push({ name: 'profil', params: { username: user } })
}




const firstname = ref('')
const lastname = ref('')
const username = ref('')
const users = ref([])
const profile=ref('')
const searchQuery=ref('')
const okurGör=ref(false)
const kitapGör=ref(true)
const books=ref([])
const searchBookQuery = ref("")
const pendingRequests = ref([])
const acceptedFollows = ref([])

const unreadMessageCount = ref(0)
let messagePollingInterval = null

const fetchUnreadMessageCount = async () => {
  if (!username.value) return
  try {
    const res = await fetch(`http://localhost:8000/api/messages/unreadCount?currentUser=${username.value}`)
    if (res.ok) {
      const data = await res.json()
      unreadMessageCount.value = data.count || 0
    }
  } catch (err) {
    console.error('Mesaj bildirimleri alınamadı:', err)
  }
}

const filteredBooks = computed(() => {
  if (!searchBookQuery.value) return books.value
  return books.value.filter(book =>
    book.title.toLowerCase().includes(searchBookQuery.value.toLowerCase())
  )
})



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
      profile.value=data.photo || ""
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

const fetchNotifications = async () => {
  try {
    const res = await fetch(`http://localhost:8000/api/getNotifications?username=${username.value}`)
    const data = await res.json()
    if (res.ok) {
      pendingRequests.value = data.pendingRequests || []
      acceptedFollows.value = data.acceptedFollows || []
    } else {
      console.error('Bildirimler alınamadı:', data.error)
    }
  } catch (err) {
    console.error('Bildirim fetch hatası:', err)
  }
}

const approveRequest = async (follower) => {
  try {
    const res = await fetch('http://localhost:8000/api/followUser', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        follower,
        following: username.value,
        action: 'approve'
      })
    })
    const data = await res.json()
    if (res.ok) {
      await fetchNotifications()
      searchUsers()
    } else {
      console.error('Onaylama hatası:', data.error)
    }
  } catch (err) {
    console.error('Onaylama hatası:', err)
  }
}

const rejectRequest = async (follower) => {
  try {
    const res = await fetch('http://localhost:8000/api/followUser', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        follower,
        following: username.value,
        action: 'reject'
      })
    })
    const data = await res.json()
    if (res.ok) {
      await fetchNotifications()
      searchUsers()
    } else {
      console.error('Reddetme hatası:', data.error)
    }
  } catch (err) {
    console.error('Reddetme hatası:', err)
  }
}

// takip etme / takipten çıkma
const follow = async (targetUsername, isFollowing, isPending) => {
  const action = isFollowing || isPending ? 'unfollow' : 'follow'
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
      await searchUsers()
      await fetchNotifications()
    } else console.error('Takip hatası:', data.error)
  } catch (err) { console.error('Takip hatası:', err) }
}


const okurAra=()=>{
  kitapGör.value=false
  okurGör.value=true
}
const kitapAra=()=>{
  okurGör.value=false
  kitapGör.value=true
}


const fetchAllBooks= async () => {
  try{
    const res= await fetch(`http://localhost:8000/api/getAllBooks?currentUser=${username.value}`);
    const data=await res.json()
    if(res.ok){
      books.value=data.map(book=>({
        ...book,
        showComments: false,
        comments: [],
        newComment: ''
      }));
    }
    else { console.error('Kitap Bulunamadı:',data.error) }
  }
  catch (err) { console.error('Kitap fetch hatası:', err) }
}


const toggleLike = async (book) => {
  const originalStatus = book.is_liked;
  const originalCount = book.like_count;

  book.is_liked = !book.is_liked;
  book.like_count += book.is_liked ? 1 : -1;

  try {
    const res = await fetch('http://localhost:8000/api/toggleLike', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value, // Giriş yapan kullanıcı
        book_id: book.id
      })
    });

    if (!res.ok) throw new Error('İşlem başarısız');
  } catch (err) {
    // 3. Hata olursa eski haline döndür
    console.error('Like hatası:', err);
    book.is_liked = originalStatus;
    book.like_count = originalCount;
    alert("İşlem başarısız oldu.");
  }
}


const fetchComments = async (bookId) => {
    const res = await fetch(`http://localhost:8000/api/getComments?book_id=${bookId}`);
    const data = await res.json();
};


//yorumları açıp kapama ve yorumları fetchleme
const toggleComments = async (book) => {
  book.showComments = !book.showComments;

  //eğer daha önce yorumlar yüklenmemişse backendden çek
  if (book.showComments && book.comments.length === 0) {
    try {
      const res = await fetch(`http://localhost:8000/api/getComments?book_id=${book.id}`);
      const data = await res.json();
      book.comments = data || [];
    } catch (err) { console.error('Yorum yükleme hatası:', err); }
  }
};

const submitComment = async (book) => {
  if (!book.newComment.trim()) return;

  try {
    const res = await fetch('http://localhost:8000/api/addComment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value,
        book_id: book.id,
        content: book.newComment
      })
    });
    if (res.ok) {
      const saveComment = await res.json();
      book.comments.push(saveComment); //listeye yeni yorumu ekle
      book.newComment = ''; //inputu temizle
    }
  } catch (err) {
    console.error('Yorum gönderme hatası:', err);
  }
};

const deleteComment = async (book, commentId) => {
  // Kullanıcıdan onay al
  if (!confirm('Bu yorumu silmek istediğinizden emin misiniz?')) {
    return;
  }

  try {
    const res = await fetch('http://localhost:8000/api/deleteComment', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        comment_id: commentId,
        username: username.value
      })
    });
    if (res.ok) {
      // Silinen yorumu listeden çıkar
      book.comments = book.comments.filter(c => c.id !== commentId);
      book.comment_count -= 1;
    } else {
      const error = await res.json();
      alert(error.error || 'Yorum silinemedi');
    }
  } catch (err) {
    console.error('Yorum silme hatası:', err);
    alert('Yorum silinemedi');
  }
};


// Menüyü açıp kapatan fonksiyon
const toggleNotifications = () => {
  isNotificationsOpen.value = !isNotificationsOpen.value;
};

const closeNotifications = (e) => {
  if (!e.target.closest('.notifications-wrapper')) {
    isNotificationsOpen.value = false;
  }
};


onMounted(async () => {
  username.value = localStorage.getItem('username') || ''
  if (!username.value) {
    router.push('/')
    return
  }
  
  await fetchProfile()
  await fetchAllBooks()
  await fetchNotifications()
  
  fetchUnreadMessageCount()
  messagePollingInterval = setInterval(fetchUnreadMessageCount, 3000)

  window.addEventListener('click', closeNotifications);
})

onUnmounted(() => {
  window.removeEventListener('click', closeNotifications);
  if (messagePollingInterval) clearInterval(messagePollingInterval)
})

</script>

<style scoped>
/* ==================== DARK MODE STİLLERİ ==================== */

/* Ana Arka Plan */
:global(.dark) {
  background-color: #0f1419;
  color: #e3e3e3;
}

/* Kitap Kartları */
:global(.dark) .bg-red-50 {
  background-color: #1e2535;
  color: #e3e3e3;
}

/* Beyaz Alanlar - Yorum Bölümleri */
:global(.dark) .bg-white {
  background-color: #252d3d;
  color: #e3e3e3;
  border-color: #3a4557;
}

/* Gri Background Elements */
:global(.dark) .bg-gray-200 {
  background-color: #252d3d;
  color: #e3e3e3;
  border-color: #3a4557;
}

/* Text Colors */
:global(.dark) .text-red-600 {
  color: #ff7a7a;
}

:global(.dark) .text-gray-600 {
  color: #b0b0b0;
}

:global(.dark) .text-gray-700 {
  color: #a0a0a0;
}

:global(.dark) .text-gray-500 {
  color: #808080;
}

:global(.dark) .text-black {
  color: #e3e3e3;
}

:global(.dark) .text-xl {
  color: inherit;
}

/* Borders */
:global(.dark) .border-red-500 {
  border-color: #ff6b6b;
}

:global(.dark) .border-gray-300,
:global(.dark) .border-gray-400 {
  border-color: #3a4557;
}

:global(.dark) .border-gray-100 {
  border-color: #2d3748;
}

/* Form Elements */
:global(.dark) input {
  background-color: #252d3d;
  color: #e3e3e3;
  border-color: #3a4557;
}

:global(.dark) input::placeholder {
  color: #808080;
}

:global(.dark) button {
  transition: all 0.3s ease;
}

:global(.dark) button.bg-gray-200 {
  background-color: #3a4557;
  color: #e3e3e3;
}

:global(.dark) button.bg-gray-200:hover {
  background-color: #4a5568;
}

:global(.dark) button.bg-red-500 {
  background-color: #c33333;
}

:global(.dark) button.bg-red-500:hover {
  background-color: #a52a2a;
}

:global(.dark) button.bg-black {
  background-color: #2d3748;
  color: #e3e3e3;
}

:global(.dark) button.bg-black:hover {
  background-color: #3a4557;
}

/* Heart Icon Colors */
:global(.dark) .fa-heart {
  color: #ff7a7a;
}

:global(.dark) .text-red-500 {
  color: #ff7a7a;
}

:global(.dark) .hover\:text-red-500:hover {
  color: #ffa0a0;
}

/* Shadow Effects */
:global(.dark) .shadow-lg {
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
}

:global(.dark) .shadow-sm {
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

/* Smooth Transitions */
:global(.dark) * {
  transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease;
}
</style>
