<template>
  <div class="pl-7 pt-5">
    <h1 class="text-[#b91c1c] text-3xl font-medium italic font-serif flex  cursor-default">
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
      <div  class="flex flex-col cursor-pointer mt-5">
        <i v-show="isOwner" @click="addBook" class="fa-solid fa-plus bg-red-600 text-white py-4 pl-8 rounded-full hover:bg-red-700 transition-colors duration-200 ease-in-out">Add Book</i>
      </div>
    </div>

    <!-- Sağ taraf: Profil bilgileri -->
    <div class="w-72">
      <div class="flex gap-10 items-center">
        <h1 class="username text-2xl font-semibold">{{ username }}</h1>
        <div  class="flex items-center gap-2">
          <button v-show="isOwner" class="border border-black text-black font-small rounded-[15px] px-4 py-1 hover:bg-gray-100 hover:border-gray-100 hover:border-red-600 transition-colors duration-300 ease-in-out" @click="showEditModal = true">Edit Profile</button>
          <i v-show="isOwner" @click="showSettingsModal = true" class="fa-solid fa-gear text-red-700 text-xl cursor-pointer hover:scale-110 transition-transform"></i>
        </div>
      </div>
        <div>
            <p v-if="!isOwner" @click="toggleFollow" class="takipEt font-semibold text-sm border rounded-md px-5 py-1 cursor-pointer max-w-max transition-colors"
               :class="isFollowing ? 'bg-transparent border-gray-500 text-gray-500 hover:bg-red-900/10' : isPending ? 'bg-gray-200 text-gray-700 border-gray-200 hover:bg-gray-300' : 'bg-red-700 text-white border-red-700 hover:bg-red-800'">
              {{ isFollowing ? 'Takipten Çık' : isPending ? 'İstek Gönderildi' : 'Takip Et' }}
            </p>
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
        <p class="mt-4 text-sm text-gray-500">Profil Gizliliği: <span class="font-semibold text-gray-800">{{ profilePrivacy }}</span></p>
      </div>
    </div>
  </div>

  <div v-if="!profileVisible" class="mt-8 p-6 bg-red-50 rounded-xl text-center text-red-700 border border-red-200">
    {{ privacyMessage }}
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
          <button  @click="removePP" class="bg-gray-200 border border-gray-600 px-1 py-[2px] rounded-[2px] mb-4">remove photo</button>
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

    <!-- Ayarlar Modal -->
    <transition name="zoom">
      <div v-if="showSettingsModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showSettingsModal = false">
        <div class="bg-white p-6 rounded-lg w-80 shadow-lg" @click.stop>
          <h2 class="text-xl font-semibold mb-6 flex items-center">
            <i class="fa-solid fa-sliders text-red-700 mr-2"></i>
            Ayarlar
          </h2>

          <!-- Tema Seçimi -->
          <div class="mb-6">
            <label class="block text-sm font-semibold mb-3 text-gray-700">Tema Seçimi</label>
            <div class="flex gap-3">
              <button
                @click="setTheme('light')"
                :class="[
                  'flex-1 py-2 px-3 rounded-lg font-medium transition-all duration-200 flex items-center justify-center gap-2',
                  currentTheme === 'light'
                    ? 'bg-red-600 text-white shadow-lg'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                ]"
              >
                <i class="fa-solid fa-sun"></i>
                Açık Mod
              </button>
              <button
                @click="setTheme('dark')"
                :class="[
                  'flex-1 py-2 px-3 rounded-lg font-medium transition-all duration-200 flex items-center justify-center gap-2',
                  currentTheme === 'dark'
                    ? 'bg-red-600 text-white shadow-lg'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                ]"
              >
                <i class="fa-solid fa-moon"></i>
                Karanlık Mod
              </button>
            </div>
          </div>

          <!-- Profil Gizliliği -->
          <div class="mb-6">
            <label class="block text-sm font-semibold mb-3 text-gray-700">Profil Gizliliği</label>
            <div class="space-y-2">
              <label class="flex items-center gap-3 cursor-pointer">
                <input type="radio" value="Herkes" v-model="profilePrivacy" class="accent-red-600" />
                <span>Herkes</span>
              </label>
              <label class="flex items-center gap-3 cursor-pointer">
                <input type="radio" value="Gizli Profil" v-model="profilePrivacy" class="accent-red-600" />
                <span>Gizli Profil</span>
              </label>
            </div>
          </div>

          <!-- Hesap Ayarları -->
          <div class="mb-6">
            <label class="block text-sm font-semibold mb-3 text-gray-700">Hesap Ayarları</label>
            <div class="space-y-3">
              <input type="text" v-model="newUsername" placeholder="Yeni Kullanıcı Adı (İsteğe Bağlı)" 
                     class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-red-600 text-sm" />
              <input type="password" v-model="oldPassword" placeholder="Eski Şifre" 
                     class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-red-600 text-sm" />
              <input type="password" v-model="newPassword" placeholder="Yeni Şifre" 
                     class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-red-600 text-sm" />
              <p class="text-[10px] text-gray-500">* Şifre değiştirmek için eski şifrenizi girmelisiniz.</p>
            </div>
          </div>

          <!-- Kapatma ve Kaydet Butonları -->
          <div class="flex justify-end gap-3">
            <button
              @click="showSettingsModal = false"
              class="px-4 py-2 bg-gray-300 text-gray-800 rounded-lg hover:bg-gray-400 transition-colors"
            >
              Kapat
            </button>
            <button
              @click="saveProfile(); saveAccountSettings(); showSettingsModal = false"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
            >
              Kaydet
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!--takipçileri görme-->
    <transition name="zoom">
      <div v-if="showFollowers" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 flex flex-col" @click="showFollowers=false">
        <div class="bg-white p-4 rounded-lg w-80 max-h-[80vh] overflow-y-auto" @click.stop>
          <h2 class="text-xl font-bold mb-4 text-center border-b pb-2">Followers</h2>
          <div v-if="followers.length === 0" class="text-center text-gray-500 py-4">No followers yet.</div>
          <div v-for="f in followers" :key="f.username" class="flex items-center gap-3 p-2 hover:bg-gray-100 rounded-lg cursor-pointer transition-colors" @click="goToProfile(f.username)">
            <img :src="f.photo || defaultPhoto" class="w-10 h-10 rounded-full object-cover border border-gray-300" alt="profile" />
            <span class="font-medium">{{ f.username }}</span>
          </div>
        </div>
      </div>
    </transition>

    <!--takip etiklerimi görme-->
    <transition name="zoom">
      <div v-if="showFollowing" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 flex flex-col" @click="showFollowing=false">
        <div class="bg-white p-4 rounded-lg w-80 max-h-[80vh] overflow-y-auto" @click.stop>
          <h2 class="text-xl font-bold mb-4 text-center border-b pb-2">Following</h2>
          <div v-if="following.length === 0" class="text-center text-gray-500 py-4">Not following anyone yet.</div>
          <div v-for="f in following" :key="f.username" class="flex items-center gap-3 p-2 hover:bg-gray-100 rounded-lg cursor-pointer transition-colors" @click="goToProfile(f.username)">
            <img :src="f.photo || defaultPhoto" class="w-10 h-10 rounded-full object-cover border border-gray-300" alt="profile" />
            <span class="font-medium">{{ f.username }}</span>
          </div>
        </div>
      </div>
    </transition>

  <!-- Kitap kartları -->
  <div class="mt-20 grid grid-cols-3 gap-20 w-2/3 m-auto">
    <div v-for="book in books" :key="book.id" class="relative rounded-2xl shadow-lg bg-red-100 p-6 transition-all duration-300">
      <img :src="book.image" class="w-full h-40 object-cover mt-2 rounded-[5px]"  alt="book image" />
      <div class="px-4 py-4 overflow-hidden">
        <h2 class="text-xl font-bold mb-2">{{ book.title }}</h2>
        <p class="text-gray-600 text-sm mb-4" :style="book.expanded ? 'max-height:none' : 'max-height:4rem; overflow:hidden;'">{{ book.description }}</p>
        <button class="bg-red-600 hover:bg-red-700 text-white font-semibold py-2 px-4 rounded-lg" @click="book.expanded = !book.expanded">{{ book.expanded ? 'Show Less' : 'Read More' }}</button>
        <button v-show="isOwner" @click="openEditBook(book)" class="bg-black text-white font-semibold py-2 px-4 rounded-lg ml-2">Edit</button>
        <button v-show="isOwner" @click="deleteBook(book.id)" class="border border-black rounded-[5px] mt-3 px-[70px] py-1 font-semibold ">Delete</button>
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
import { ref, onMounted, watch, computed} from 'vue'
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
const showSettingsModal = ref(false)
const editingBook = ref({})
const currentTheme = ref('light')
const profilePrivacy = ref('Herkes')

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
const isFollowing = ref(false)
const isPending = ref(false)
const profileVisible = ref(true)
const privacyMessage = ref('')

const newUsername = ref('')
const oldPassword = ref('')
const newPassword = ref('')




//profil fetchleme
const fetchProfile = async () => {
  if (!username.value) {
    console.error('Username bulunamadı')
    return
  }
  try {
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/profile?username=${username.value}&currentUser=${currentUser}`)
    if (res.status === 403) {
      profileVisible.value = false
      profilePrivacy.value = 'Gizli Profil'
      privacyMessage.value = 'Bu profili görüntüleme yetkiniz yok'
      return
    }
    const data = await res.json()
    if (res.ok) {
      firstname.value=data.firstname || ""
      lastname.value=data.lastname || ""
      bio.value = data.bio || ''
      photo.value = data.photo || defaultPhoto;
      profilePrivacy.value = data.privacy || 'Herkes'
      profileVisible.value = true
      privacyMessage.value = ''
    } else {
      console.error('Profil Bulunamadı:',data.error)
    }
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
      photo: photo.value,
      privacy: profilePrivacy.value,
    };
    await fetch("http://localhost:8000/api/updateProfile", {
      method: "PUT",
      headers:{"Content-Type": "application/json",},
      body: JSON.stringify(updatedProfile)
    });
    showEditModal.value=false;
    await fetchProfile();
  } catch(err){
    console.log("Profile edit error:", err)
  }
}

//photo işlemleri
const addProfilePhoto = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  try {
    const formData = new FormData()
    formData.append('username', username.value)
    formData.append('photo', file)

    const res = await fetch('http://localhost:8000/api/uploadProfilePhoto', {
      method: 'POST',
      body: formData,
    })
    const data = await res.json()
    if (res.ok) {
      photo.value = data.photo || defaultPhoto
      await fetchProfile()
    } else {
      console.error('Photo upload failed:', data.error)
    }
  } catch (err) {
    console.error('Photo upload error:', err)
  }
}
const removePP = async () => {
  const onay = confirm("Profil Fotoğrafınızı silmek istediğinizden emin misiniz?")
  if (!onay) return
  try {
    const res = await fetch('http://localhost:8000/api/removeProfilePhoto', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value }),
    })
    if (res.ok) {
      photo.value = defaultPhoto
      await fetchProfile()
    } else {
      const data = await res.json()
      console.error('Remove photo failed:', data.error)
    }
  } catch (err) {
    console.error('Remove photo error:', err)
  }
}


const saveAccountSettings = async () => {
  if (!newUsername.value && !newPassword.value) return;

  try {
    const res = await fetch('http://localhost:8000/api/changeAccountSettings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        currentUsername: localStorage.getItem('username'),
        newUsername: newUsername.value,
        oldPassword: oldPassword.value,
        newPassword: newPassword.value
      })
    })
    
    if (res.ok) {
      if (newUsername.value && newUsername.value !== localStorage.getItem('username')) {
        localStorage.setItem('username', newUsername.value)
        router.push(`/profil/${newUsername.value}`)
      }
      newUsername.value = ''
      oldPassword.value = ''
      newPassword.value = ''
      alert('Hesap ayarları başarıyla güncellendi!')
    } else {
      const errData = await res.text()
      alert('Hata: ' + errData)
    }
  } catch (err) {
    console.error('Hesap ayarları kaydedilirken hata:', err)
  }
}
// Photo is loaded from backend in fetchProfile(); no localStorage fallback used anymore.


//kitap fetchleme
const fetchBooks = async () => {
  try {
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/getBooks?username=${username.value}&currentUser=${currentUser}`)
    if (res.status === 403) {
      books.value = []
      return
    }
    const data = await res.json()
    if (res.ok) {
      // If data is null because of DB handling, fallback to []
      books.value = data || []
    } else {
      console.error('Kitap Bulunamadı:', data.error)
    }
  } catch (err) {
    console.error('Kitap fetch hatası:', err)
  }
}

// Kitap ekleme
const addBook = async () => {
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
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/getFollowers?username=${username.value}&currentUser=${currentUser}`)
    if (res.status === 403) {
      followersCount.value = 0
      followingCount.value = 0
      return
    }
    const data = await res.json()
    if (res.ok) {
      followersCount.value = data.followers || 0
      followingCount.value = data.following || 0
    }
  } catch (err) {
    console.error('Takipçi bilgisi çekilemedi:', err)
  }
}

const fetchFollowersList=async()=>{
  try {
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/getFollowersList?username=${username.value}&currentUser=${currentUser}`)
    if (res.status === 403) {
      followers.value = []
      return
    }
    const data = await res.json()
    if(res.ok){
      followers.value=data || []
    }
  } catch (err) {
    console.error('Takipçi isimleri çekilemedi:', err)
  }
}
const fetchIsFollowing = async () => {
  try {
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/checkFollowStatus?follower=${currentUser}&following=${username.value}`)
    const data = await res.json()
    if (res.ok) {
      isFollowing.value = data.isFollowing
      isPending.value = data.isPending
    } else {
      isFollowing.value = false
      isPending.value = false
    }
  } catch (err) {
    console.error('isFollowing fetch hata:', err)
    isFollowing.value = false
    isPending.value = false
  }
}

const goToProfile = (uname) => {
  showFollowers.value = false;
  showFollowing.value = false;
  router.push(`/profil/${uname}`);
}

const toggleFollow = async () => {
  try {
    const me = localStorage.getItem('username') || ''
    if (!me) { alert('Önce giriş yapmalısınız'); return }
    const action = isFollowing.value || isPending.value ? 'unfollow' : 'follow'
    const res = await fetch('http://localhost:8000/api/followUser', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ follower: me, following: username.value, action })
    })
    if (res.ok) {
      // refresh follower counts and state
      await fetchFollowers()
      await fetchFollowersList()
      await fetchIsFollowing()
    } else {
      const d = await res.json().catch(()=>({}))
      console.error('follow request failed', d)
    }
  } catch (err) {
    console.error('toggleFollow error:', err)
  }
}
const fetchFollowingList=async()=>{
  try {
    const currentUser = localStorage.getItem('username') || ''
    const res = await fetch(`http://localhost:8000/api/getFollowingList?username=${username.value}&currentUser=${currentUser}`)
    if (res.status === 403) {
      following.value = []
      return
    }
    const data = await res.json()
    if(res.ok){
      following.value=data || []
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
  fetchIsFollowing()
})

// computed owner flag
const isOwner = computed(() => username.value && username.value === localStorage.getItem('username'))

// Tema İşlemleri
const loadTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'light'
  currentTheme.value = savedTheme
}

const setTheme = (theme) => {
  currentTheme.value = theme
  localStorage.setItem('theme', theme)
  // Global tema değiş fonksiyonunu çağır
  if (window.toggleTheme && theme !== (window.getTheme?.() || 'light')) {
    window.toggleTheme()
  }
}


onMounted(() => {
  loadTheme()
  fetchProfile()
  fetchBooks()
  fetchFollowers()
  fetchFollowersList()
  fetchFollowingList()
  fetchIsFollowing()
})



 const goMainPg=()=>{
    router.push("/anasayfa")
  }

</script>

<style scoped>
.zoom-enter-active, .zoom-leave-active { transition: all 0.3s ease; }
.zoom-enter-from, .zoom-leave-to { transform: scale(0.5); opacity: 0; }
.zoom-enter-to, .zoom-leave-from { transform: scale(1); opacity: 1; }

/* Dark Mode Stilleri */
:global(.dark) {
  background-color: #1a1a1a;
  color: #f5f5f5;
}

:global(.dark) .bg-white {
  background-color: #2d2d2d;
  color: #f5f5f5;
}

:global(.dark) .text-black,
:global(.dark) .text-gray-700,
:global(.dark) .text-gray-600 {
  color: #e3e3e3;
}

:global(.dark) .border-gray-300,
:global(.dark) .border-gray-400 {
  border-color: #404040;
}

:global(.dark) .bg-gray-200 {
  background-color: #3a3a3a;
  color: #f5f5f5;
}

:global(.dark) .bg-red-50 {
  background-color: #2d1f1f;
}

:global(.dark) .bg-red-100 {
  background-color: #3a2424;
}

:global(.dark) .bg-gray-100 {
  background-color: #423131;
}

:global(.dark) input,
:global(.dark) textarea {
  background-color: #3a3a3a;
  color: #f5f5f5;
  border-color: #404040;
}

:global(.dark) button.border-black {
  color: #f5f5f5;
  border-color: #404040;
}

:global(.dark) button.border-black:hover {
  background-color: #3a3a3a;
}
</style>
