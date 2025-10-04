<template>
  <div class="pl-7 pt-5">
    <h1 class="text-red-500 text-3xl font-medium italic font-serif flex  cursor-default">
        Readly
        <i class="fa-solid fa-book-open text-2xl px-2 hover:scale-110 transition-transform transition-colors duration-700 border-2 border-transparent rounded-full"></i>
    </h1>
    <i @click="goMainPg" class="fa-regular fa-circle-left text-4xl  cursor-pointer"></i>
  </div>

  <div class="flex justify-center gap-20 relative">

    <!-- Sol taraf: Profil resmi ve kitap ekleme -->
    <div>
      <div>
        <img v-if="photo" :src="photo" @click="showModal=true" class="w-[11rem] h-[11rem] rounded-full object-cover mt-2 cursor-pointer border border-gray-400"  alt="profilePhoto" />
      </div>
      <div v-if="username === localStorage.getItem('username')" class="flex flex-col cursor-pointer mt-5">
        <i @click="addBook" class="fa-solid fa-plus bg-gray-300 text-gray-800 py-5 pl-9 rounded-full hover:bg-gray-400 transition">Add Book</i>
      </div>
    </div>

    <!-- Sağ taraf: Profil bilgileri -->
    <div class="w-72">
      <div class="flex gap-10 items-center">
        <h1 class="username text-2xl font-semibold">{{ username }}</h1>
        <div v-if="username === localStorage.getItem('username')" class="flex items-center gap-2">
          <button class="bg-black text-white font-medium rounded-[15px] px-4 py-1 hover:bg-gray-800" @click="showEditModal = true">Edit Profile</button>
          <i class="fa-solid fa-gear text-xl cursor-pointer"></i>
        </div>
      </div>
      <div class="flex gap-12 pt-2">
        <div @click="showFollowers=true" class="flex cursor-pointer">
            <p class="font-bold text-[18px]">{{ followersCount }}&nbsp;</p>
            <p> followers</p>
        </div>
        <div @click="showFollowing=true" class="flex cursor-pointer">
            <p class="font-bold text-[18px]">{{ followingCount }}&nbsp;</p>
            <p> following</p>
        </div>
      </div>
      <div class="pt-5">
        <p class="isim font-semibold text-xl">{{ firstname }} {{ lastname }}</p>
        <p class="biography pt-1 text-gray-500 font-light leading-relaxed">{{ bio }}</p>
      </div>
    </div>

    <!-- Profil resmi büyültme  -->
    <transition name="zoom">
      <div v-if="showModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showModal=false">
        <img :src="photo" class="w-96 h-96 rounded-full shadow-lg cursor-pointer" @click="showModal = false">
      </div>
    </transition>

    <!-- Profil düzenleme (edit profile kısmını büyültme) -->
    <transition name="zoom">
      <div v-if="showEditModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showEditModal = false">
        <div class="bg-white p-6 rounded-lg w-80" @click.stop>
          <h2 class="text-xl font-semibold mb-4">Edit Profile</h2>
          <div class="flex  block font-medium">
              <i class="fa-regular fa-image pt-1"></i>
              <p class="pl-1">Profile Photo</p>
          </div>
          <input type="file" accept="image/*" @change="addProfilePhoto" class="pb-4"> <!--accept:sadece resim formatındaki dosyaların seçilmesine izin verir.-->
          <button v-if="username===localStorage.getItem('username')" @click="removePP" class="bg-gray-200 border border-gray-600 px-1 py-[2px] rounded-[2px] mb-4">remove photo</button>
          <label class="block mb-2 font-medium">Biography</label>
          <textarea v-model="bio" class="w-full border border-gray-300 rounded px-2 py-1 mb-4"></textarea>
          <label class="block mb-2 font-medium">First Name</label>
          <input v-model="firstname" type="text" class="w-full border border-gray-300 rounded px-2 py-1 mb-4" />
          <label class="block mb-2 font-medium">Last Name</label>
          <input v-model="lastname" type="text" class="w-full border border-gray-300 rounded px-2 py-1 mb-4" />
          <div class="flex justify-end gap-2">
            <button class="px-4 py-1 rounded bg-black text-white" @click="saveProfile">Save</button>
          </div>
        </div>
      </div>
    </transition>

    <!--takipçileri görme-->
    <transition name="zoom">
      <div v-if="showFollowers" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 flex flex-col" @click="showFollowers=false">
        <div v-for="f in followers" :key="f"  class="bg-white p-3 pl-7 rounded-lg w-80 flex flex-col m-[1px]" @click.stop>
          <div>{{ f }}</div>
        </div>
      </div>
    </transition>

    <!--takip etiklerimi görme-->
    <transition name="zoom">
      <div v-if="showFollowing" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 flex flex-col" @click="showFollowing=false">
        <div v-for="f in following" :key="f" class="bg-white p-3 pl-7 rounded-lg w-80 flex flex-col m-[1px]" @click.stop>
          <div>{{ f }}</div>
        </div>
      </div>
    </transition>

  </div>

  <!-- Kitap kartları -->
  <div class="mt-10 grid grid-cols-3 gap-20 w-2/3 m-auto">
    <div v-for="book in books" :key="book.id" class="relative rounded-2xl shadow-lg bg-white p-6 transition-all duration-300">
      <img :src="book.image" class="w-full h-40 object-cover mt-2 rounded-[5px]"  alt="book image" />
      <div class="px-4 py-4 overflow-hidden">
        <h2 class="text-xl font-bold mb-2">{{ book.title }}</h2>
        <p class="text-gray-600 text-sm mb-4" :style="book.expanded ? 'max-height:none' : 'max-height:4rem; overflow:hidden;'">{{ book.description }}</p>
        <button class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg" @click="book.expanded = !book.expanded">{{ book.expanded ? 'Show Less' : 'Read More' }}</button>
        <button @click="openEditBook(book)" class="bg-gray-300 hover:bg-gray-400 font-semibold py-2 px-4 rounded-lg ml-2">Edit</button>
        <button @click="deleteBook(book.id)" class="border border-black rounded-[5px] mt-3 px-[70px] py-1 font-semibold hover:bg-gray-100">Delete</button>
      </div>
    </div>
  </div>

  <!-- Kitap düzenleme modal -->
  <transition name="zoom">
    <div v-if="showEditBook" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showEditBook = false">
      <div class="bg-white p-6 rounded-lg w-80" @click.stop>
        <div class="flex  block mb-2 font-medium">
            <i class="fa-regular fa-image pt-1"></i>
            <p class="pl-1">Book Image</p>
        </div>
        <input type="file"   accept="image/*"   @change="e=>addBookImage(e,editingBook)"  class="pb-4">
        <label class="block mb-2 font-medium">Book Name</label>
        <input v-model="editingBook.title" type="text" class="w-full border border-gray-300 rounded px-2 py-1 mb-4" />
        <label class="block mb-2 font-medium">Book Description</label>
        <textarea v-model="editingBook.description" class="w-full border border-gray-300 rounded px-2 py-1 mb-4"></textarea>
        <div class="flex justify-end gap-2">
          <button class="px-4 py-1 rounded bg-black text-white" @click="saveEditedBook">Save</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted, watch} from 'vue'
import { useRouter} from 'vue-router'
import defaultPhoto from "@/assets/defaultPhoto.png"

const router=useRouter()

// route parametreleri props ile geliyor
const props = defineProps({
  username: { type: String, required: false } // route paramı gelmezse localStorage fallback'i olacak
})
// route paramı yoksa localStorage'tan al, yoksa boş string
const username = ref(props.username || localStorage.getItem('username') || '')



const showModal = ref(false)
const showEditModal = ref(false)
const showEditBook = ref(false)
const editingBook = ref({})

const firstname = ref('')
const lastname = ref('')
const bio = ref('')
const books = ref([])
const photo = ref(defaultPhoto)

const followersCount=ref(0)
const followingCount=ref(0)
const showFollowers=ref(false)
const showFollowing=ref(false)
const followers = ref([])
const following = ref([])




//profil fetchleme
const fetchProfile = async () => {
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
      bio.value = data.bio || ''
    } else { console.error('Profil Bulunamadı:',data.error) }
  } catch (err) { console.error('Profil fetch hatası:', err) }
}

//profili düzenle ve kaydet
const saveProfile = async () => {
  try{
    const updatedProfile ={
      username: username.value,
      firstname: firstname.value,
      lastname: lastname.value,
      bio: bio.value,
    };
    await fetch("http://localhost:8000/api/updateProfile", {
      method: "PUT",
      headers:{
          "Content-Type": "application/json",
      },
      body: JSON.stringify(updatedProfile)
    });
    showEditModal.value=false;
    await fetchProfile();
  } catch(err){
    console.log("Profile edit error:", err)
  }
}

const addProfilePhoto = (event) => {
  const file = event.target.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = e => {
    photo.value = e.target.result
    localStorage.setItem(`profilePhoto_${username.value}`, e.target.result)
  }
  reader.readAsDataURL(file)
}

const removePP=()=>{
  const onay=confirm("Profil Fotoğrafınızı silmek istediğinizden emin misiniz?")
  if(!onay){return}
  photo.value=defaultPhoto
  localStorage.removeItem(`profilePhoto_${username.value}`)
}
const savedPhoto = localStorage.getItem(`profilePhoto_${localStorage.getItem("username")}`)
if (savedPhoto) {
  photo.value = savedPhoto
}
else{
  photo.value=defaultPhoto
}


//kitap fetchleme
const fetchBooks = async () => {
  try {
    const res = await fetch(`http://localhost:8000/api/getBooks?username=${username.value}`)
    const data = await res.json()
     if (res.ok) {
      books.value=data
    } else { console.error('Kitap Bulunamadı:',data.error) }
  } catch (err) { console.error('Kitap fetch hatası:', err) }
}

// Kitap ekleme
const addBook = async () => {
  if (username.value!== localStorage.getItem('username')){
    console.error('sadece kendi profiline kitap ekleyebilirsin.')
    return
  }

  const newBook = { title: 'New Book', description: 'Book Description', username: username.value }
  try {
    const res = await fetch('http://localhost:8000/api/addBook', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newBook)
    })
    if (!res.ok) console.error('Kitap eklenemedi')
    await fetchBooks()
  } catch(err) { console.error(err) }
}


// Kitap silme
const deleteBook = async (id) => {
  if(username!==localStorage.getItem('username')){
    console.error('sadece kendi profiline ait kitapları silebilirsin')
    return
  }
  const onay = confirm("Bu kitabı silmek istediğine emin misin?")
  if (!onay) return
  await fetch('http://localhost:8000/api/deleteBook', {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id, username: username.value })
  })
  fetchBooks()
}

// Kitap editleme
const openEditBook = (book) => {
    if (username.value !== localStorage.getItem('username')) {
    console.error("Sadece kendi profiline ait kitapları düzenleyebilirsin")
    return
  }
  editingBook.value = { ...book}
  showEditBook.value = true
}

// Kitap düzenlemeyi kaydet
const saveEditedBook = async () => {
  try {
    await fetch("http://localhost:8000/api/updateBook", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({...editingBook.value, username: username.value }),
    });
    showEditBook.value = false;
    await fetchBooks();
  } catch (err) {
    console.log("Book edit error:", err);
  }
};


const addBookImage = (event, book) => {
  const file = event.target.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = e => {
    book.image = e.target.result
  }
  reader.readAsDataURL(file)
}


const fetchFollowers=async()=>{
    try {
    const res = await fetch(`http://localhost:8000/api/getFollowers?username=${username.value}`)
    const data = await res.json()
    if (res.ok) {
      followersCount.value = data.followers
      followingCount.value = data.following
    }
  } catch (err) {
    console.error('Takipçi bilgisi çekilemedi:', err)
  }
}

const fetchFollowersList=async()=>{
    try {
    const res = await fetch(`http://localhost:8000/api/getFollowersList?username=${username.value}`)
    const data = await res.json()
    if(res.ok){
      followers.value=data
    }
  } catch (err) {
    console.error('Takipçi isimleri çekilemedi:', err)
  }
}
const fetchFollowingList=async()=>{
    try {
    const res = await fetch(`http://localhost:8000/api/getFollowingList?username=${username.value}`)
    const data = await res.json()
    if(res.ok){
      following.value=data
    }
  } catch (err) {
    console.error('Takip ettiklerimin isimleri çekilemedi:', err)
  }
}


//bu kod her route param (props.username) değiştiğinde tetikleniyor.
//yeni kullanıcı adını username değişkenine atıyor
//bütün fetch fonksiyonlarını yeniden çalıştırıyor (profil bilgisi, kitaplar, takipçiler vs güncelleniyor)
watch(() => props.username, (newUser) => {
  username.value = newUser || localStorage.getItem('username') || ''
  fetchProfile()
  fetchBooks()
  fetchFollowers()
  fetchFollowersList()
  fetchFollowingList()
})






onMounted(() => {
  fetchProfile()
  fetchBooks()
  fetchFollowers()
  fetchFollowersList()
  fetchFollowingList()
})



 const goMainPg=()=>{
    router.push("/anasayfa")
  }

</script>

<style scoped>
.zoom-enter-active, .zoom-leave-active { transition: all 0.3s ease; }
.zoom-enter-from, .zoom-leave-to { transform: scale(0.5); opacity: 0; }
.zoom-enter-to, .zoom-leave-from { transform: scale(1); opacity: 1; }
</style>
