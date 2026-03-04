<template>
  <div >
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

          <!-- ------------arama çubukları kısmı-------------- -->
          <div class="pt-5">
            <nav class="w-full">
                <ul class="text-red-600 flex items-center tracking-widest justify-center space-x-12">
                  <li>
                    <a @click="kitapAra" href="#" class="relative after:block after:h-px after:w-0 after:bg-red-500 after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      kitap yorumlamaları
                    </a>
                  </li>
                  <li>
                    <a @click="okurAra" href="#" class="relative after:block after:h-px after:w-0 after:bg-red-500 after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      okurlar
                    </a>
                  </li>
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-red-500 after:transition-all after:duration-300 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      kitaplar
                    </a>
                  </li>
                </ul>
            </nav>
            <!--okur arama çubuğu-->
            <div v-if="okurGör" class="flex items-center border border-red-500 rounded-full px-3 my-4 w-full">
                <i class="fa-solid fa-magnifying-glass text-gray-500"></i>
                <input @keyup.enter="searchUsers" v-model="searchQuery" type="text" placeholder="kullanıcı ara.." class="w-full px-2 py-1 outline-none rounded-full"/>
                <button @click="searchUsers" class="bg-gray-200 rounded-[20px] hover:bg-gray-300 px-3 cursor-pointer">ARA</button>
            </div>
            <div v-if="okurGör" v-for="user in users.filter(u=>u.username!==username)" :key="user.username" class="flex bg-gray-200 justify-between items-center border rounded-[20px] py-2 m-2 cursor-pointer border-gray-300 ">
                <div @click="goUsersProfile(user.username)" class="flex items-center font-semibold mx-6">
                    <img v-if="user.photo" :src="user.photo" alt="profile photo" class="w-8 h-8 rounded-full mr-2 object-cover border border-gray-400">
                    <i v-else class="fa-solid fa-circle-user text-2xl mr-2 text-gray-500"></i>
                    <span>{{ user.username }}</span>
                </div>
                <button @click="follow(user.username, user.isFollowing)" class="bg-black text-white border rounded-[10px] px-2 py-[2px] font-light text-[15px] hover:bg-gray-700 mx-6">{{ user.isFollowing ? 'Takipten Çık' : 'Takip Et' }}</button>
            </div>
            <!--kitap yorumu arama çubuğu-->
            <div v-if="kitapGör" class="flex items-center border border-red-500 rounded-full px-3 my-4 w-full">
                <i class="fa-solid fa-magnifying-glass text-red-600 pr-1"></i>
                <input v-model="searchBookQuery" type="text" placeholder="kitap ara.." class="w-full px-2 py-1 outline-none rounded-full"/>
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
                    <p class="font-bold text-red-600">{{ c.username }}</p>
                    <p class="text-gray-700">{{ c.content }}</p>
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
import { ref, onMounted, computed} from 'vue'
import { useRouter } from 'vue-router'

const router=useRouter()


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



onMounted(() => {
  fetchProfile()
  fetchAllBooks()
})

</script>

