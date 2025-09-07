<template>
  <div class="flex justify-center mt-20 gap-20 relative">

    <!-- Sol taraf: Profil resmi ve kitap ekleme -->
    <div>
      <div>
        <img :src="profilePhoto" alt="photo" class="w-40 h-40 rounded-full cursor-pointer" @click="showModal = true">
      </div>
      <div class="flex flex-col cursor-pointer mt-5">
        <i @click="addBook" class="fa-solid fa-plus bg-gray-300 text-gray-800 py-5 pl-6 rounded-full hover:bg-gray-400 transition">Add Book</i>
      </div>
    </div>

    <!-- Sağ taraf: Profil bilgileri -->
    <div class="w-72">
      <div class="flex gap-10 items-center">
        <h1 class="username text-2xl font-semibold">{{ username }}</h1>
        <div class="flex items-center gap-2">
          <button class="bg-black text-white font-medium rounded-[15px] px-4 py-1 hover:bg-gray-800" @click="showEditModal = true">Edit Profile</button>
          <i class="fa-solid fa-gear text-xl cursor-pointer"></i>
        </div>
      </div>
      <div class="flex gap-12 pt-2">
        <p>Followers:</p>
        <p>Following:</p>
      </div>
      <div class="pt-5">
        <p class="isim font-semibold text-xl">{{ firstname }} {{ lastname }}</p>
        <p class="biography pt-1 text-gray-500 font-light leading-relaxed">{{ bio }}</p>
      </div>
    </div>

    <!-- Profil resmi büyütme modal -->
    <transition name="zoom">
      <div v-if="showModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50">
        <img :src="profilePhoto" alt="photo" class="w-96 h-96 rounded-full shadow-lg cursor-pointer" @click="showModal = false">
      </div>
    </transition>

    <!-- Profil düzenleme modal -->
    <transition name="zoom">
      <div v-if="showEditModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showEditModal = false">
        <div class="bg-white p-6 rounded-lg w-80" @click.stop>
          <h2 class="text-xl font-semibold mb-4">Edit Profile</h2>
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

  </div>

  <!-- Kitap kartları -->
  <div class="mt-10 grid grid-cols-3 gap-20 w-2/3 m-auto">
    <div v-for="book in books" :key="book.id" class="relative rounded-2xl shadow-lg bg-white p-6 transition-all duration-300">
      <img v-if="editingBook.image" :src="editingBook.image" class="w-full h-40 object-cover mt-2"  alt="book image" />
      <button @click="addBookImage" class="flex hover:bg-gray-200 font-medium border border-gray-300 px-1 rounded-[4px] mt-1">
          <i class="fa-regular fa-image pt-1"></i>
          <p class="pl-1">kitap resmi ekle</p>
      </button>
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
        <h2 class="text-xl font-semibold mb-4">Edit Book</h2>
        <label class="block mb-2 font-medium">Book Image</label>
        <input type="file" accept="image/*" @change="e => addBookImage(e, editingBook)">
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
import { ref, onMounted} from 'vue'
import profilePhoto from '@/assets/WIN_20240125_01_35_51_Pro.jpg'

const showModal = ref(false)
const showEditModal = ref(false)
const showEditBook = ref(false)
const editingBook = ref({})

const firstname = ref('')
const lastname = ref('')
const bio = ref('')
const username = ref('')
const books = ref([])

//profil fetchleme
const fetchProfile = async () => {
  username.value = localStorage.getItem('username') || ''

    if (!username.value) {
    console.error('Username bulunamadı')
    return
  }
  try {
    const res = await fetch(`http://localhost:8000/api/profile?username=${username.value}`)
    const data = await res.json()
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
    showEditModal.value = false;
    await fetchProfile();
  } catch(err){
    console.log("Profile edit error:", err)
  }
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
  if (!username.value) return console.error("Username yok, kitap eklenemiyor")

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
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = e => {
    const file = e.target.files[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = e => {
      book.image = e.target.result // base64 string
    }
    reader.readAsDataURL(file)
  }
  input.click()
}


onMounted(() => {
  fetchProfile()
  fetchBooks()
})
</script>

<style scoped>
.zoom-enter-active, .zoom-leave-active { transition: all 0.3s ease; }
.zoom-enter-from, .zoom-leave-to { transform: scale(0.5); opacity: 0; }
.zoom-enter-to, .zoom-leave-from { transform: scale(1); opacity: 1; }
</style>
