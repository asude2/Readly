<template>
  <div class="flex justify-between items-center px-8 py-6 w-full max-w-6xl mx-auto">
    <h1 class="text-red-600 text-3xl font-medium italic font-serif flex items-center cursor-default">
        Readly
        <i class="fa-solid fa-book-open text-2xl px-2"></i>
    </h1>
    <i @click="goMainPg" class="fa-regular fa-circle-left text-4xl text-gray-600 dark:text-gray-400 cursor-pointer hover:text-red-600 dark:hover:text-red-600 transition-colors"></i>
  </div>

  <div class="flex flex-col md:flex-row justify-center items-center md:items-start gap-12 md:gap-20 mt-8 mb-12 w-full max-w-4xl mx-auto px-6">

    <!-- Sol taraf: Profil resmi ve kitap ekleme -->
    <div class="flex flex-col items-center">
      <div class="relative group">
        <img v-if="photo" :src="photo" @click="showModal=true" class="w-40 h-40 md:w-44 md:h-44 rounded-full object-cover cursor-pointer ring-4 ring-gray-100 dark:ring-gray-200 shadow-lg group-hover:ring-red-100 dark:group-hover:ring-red-900/30 transition-all duration-300" alt="profilePhoto" />
        <div v-if="isOwner" @click="showEditModal = true" class="absolute bottom-2 right-2 bg-gray-900 dark:bg-gray-700 text-white p-2.5 rounded-full cursor-pointer shadow-md hover:bg-gray-800 transition-colors border-2 border-white dark:border-gray-800">
          <i class="fa-solid fa-pen text-sm"></i>
        </div>
      </div>
      <button v-show="isOwner" @click="addBook" class="mt-6 flex items-center gap-2 bg-red-600 text-white font-semibold py-2.5 px-6 rounded-full hover:bg-red-700 shadow-md hover:shadow-lg transition-all duration-200">
        <i class="fa-solid fa-plus"></i> Yeni Gönderi
      </button>
    </div>

    <!-- Sağ taraf: Profil bilgileri -->
    <div class="w-full md:w-[28rem] flex flex-col items-center md:items-start text-center md:text-left">
      <div class="flex flex-col md:flex-row md:items-center gap-4 md:gap-6 mb-5">
        <h1 class="text-2xl md:text-3xl font-bold text-black dark:text-white">{{ username }}</h1>
        <div class="flex items-center justify-center gap-3">
          <button v-show="isOwner" class="border border-gray-300 dark:border-slate-600 text-gray-800 dark:text-slate-300 font-medium rounded-lg px-4 py-1.5 hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors duration-200 text-sm" @click="showEditModal = true">Profili Düzenle</button>
          <i v-show="isOwner" @click="showSettingsModal = true" class="fa-solid fa-gear text-gray-400 hover:text-gray-800 dark:hover:text-slate-300 text-xl cursor-pointer transition-colors"></i>
        </div>
      </div>

      <div class="mb-5">
          <button v-if="!isOwner" @click="toggleFollow" class="font-semibold text-sm border rounded-lg px-6 py-2 transition-all duration-200 shadow-sm"
             :class="isFollowing ? 'bg-gray-100 dark:bg-slate-800 border-gray-200 dark:border-slate-700 text-gray-800 dark:text-slate-300 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 hover:border-red-200' : isPending ? 'bg-slate-200 dark:bg-slate-700 text-gray-600 dark:text-gray-400 border-transparent cursor-wait' : 'bg-red-600 text-white border-red-600 hover:bg-red-700 hover:shadow-md'">
            {{ isFollowing ? 'Takipten Çık' : isPending ? 'İstek Gönderildi' : 'Takip Et' }}
          </button>
      </div>

      <div class="flex gap-8 mb-6 bg-gray-50 dark:bg-slate-800/50 py-3 px-6 rounded-xl border border-gray-100 dark:border-slate-800">
        <div @click="showFollowers=true" class="flex flex-col items-center cursor-pointer group">
            <span class="font-bold text-xl text-black dark:text-white group-hover:text-red-600 transition-colors">{{ followersCount }}</span>
            <span class="text-sm text-gray-500 dark:text-gray-400">Takipçi</span>
        </div>
        <div class="w-px bg-slate-200 dark:bg-slate-700"></div>
        <div @click="showFollowing=true" class="flex flex-col items-center cursor-pointer group">
            <span class="font-bold text-xl text-black dark:text-white group-hover:text-red-600 transition-colors">{{ followingCount }}</span>
            <span class="text-sm text-gray-500 dark:text-gray-400">Takip</span>
        </div>
      </div>

      <div class="w-full">
        <p class="font-semibold text-lg text-black dark:text-slate-100">{{ firstname }} {{ lastname }}</p>
        <p class="mt-2 text-gray-600 dark:text-slate-300 leading-relaxed max-w-sm">{{ bio }}</p>
        <p class="mt-4 text-xs font-medium text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-slate-800 inline-flex px-3 py-1 rounded-full"><i class="fa-solid fa-lock mr-1.5 text-[10px]"></i> {{ profilePrivacy }}</p>
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
        <div class="bg-white dark:bg-slate-800 p-6 rounded-xl w-80 shadow-2xl border border-gray-100 dark:border-slate-700" @click.stop>
          <h2 class="text-xl font-bold mb-6 text-black dark:text-white">Profili Düzenle</h2>
          <div class="flex items-center block font-medium mb-2 text-gray-800 dark:text-slate-200">
              <i class="fa-regular fa-image pt-1"></i>
              <p class="pl-2 text-sm">Profil Fotoğrafı</p>
          </div>
          <input type="file" accept="image/*" @change="addProfilePhoto" class="pb-4 block w-full text-xs text-gray-500 dark:text-gray-400 file:mr-4 file:py-1 file:px-3 file:rounded-full file:border-0 file:text-xs file:font-semibold file:bg-red-50 file:text-red-700 hover:file:bg-red-100 transition-colors">
          <button @click="removePP" class="text-xs text-gray-500 hover:text-red-600 mb-6 transition-colors underline decoration-dotted">Fotoğrafı Kaldır</button>
          
          <label class="block mb-2 font-semibold text-sm text-gray-800 dark:text-slate-200">Biyografi</label>
          <textarea v-model="bio" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-4 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500 text-sm" rows="3"></textarea>
          
          <label class="block mb-2 font-semibold text-sm text-gray-800 dark:text-slate-200">Ad</label>
          <input v-model="firstname" type="text" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-4 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500 text-sm" />
          
          <label class="block mb-2 font-semibold text-sm text-gray-800 dark:text-slate-200">Soyad</label>
          <input v-model="lastname" type="text" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-6 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500 text-sm" />
          
          <div class="flex justify-end gap-2">
            <button class="px-5 py-2 rounded-lg bg-red-600 text-white font-semibold hover:bg-red-700 transition-all shadow-md active:scale-95" @click="saveProfile">Kaydet</button>
          </div>
        </div>

      </div>
    </transition>

    <!-- Ayarlar Modal -->
    <transition name="zoom">
      <div v-if="showSettingsModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showSettingsModal = false">
        <div class="bg-white dark:bg-slate-800 p-6 rounded-xl w-80 shadow-2xl border border-gray-100 dark:border-slate-700" @click.stop>
          <h2 class="text-xl font-bold mb-6 flex items-center text-black dark:text-white">
            <i class="fa-solid fa-sliders text-red-700 mr-2"></i>
            Ayarlar
          </h2>


          <!-- Tema Seçimi -->
          <div class="mb-6">
            <label class="block text-sm font-semibold mb-3 text-gray-700 dark:text-slate-300">Tema Seçimi</label>

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
            <label class="block text-sm font-semibold mb-3 text-gray-700 dark:text-slate-300">Profil Gizliliği</label>
            <div class="space-y-2 text-gray-800 dark:text-slate-300 text-sm">
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
            <label class="block text-sm font-semibold mb-3 text-gray-700 dark:text-slate-300">Hesap Ayarları</label>
            <div class="space-y-3">
              <input type="text" v-model="newUsername" placeholder="Yeni Kullanıcı Adı (İsteğe Bağlı)" 
                     class="w-full px-3 py-2 border dark:border-slate-600 rounded-lg focus:outline-none focus:border-red-600 text-sm bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white" />
              <input type="password" v-model="oldPassword" placeholder="Eski Şifre" 
                     class="w-full px-3 py-2 border dark:border-slate-600 rounded-lg focus:outline-none focus:border-red-600 text-sm bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white" />
              <input type="password" v-model="newPassword" placeholder="Yeni Şifre" 
                     class="w-full px-3 py-2 border dark:border-slate-600 rounded-lg focus:outline-none focus:border-red-600 text-sm bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white" />
              <p class="text-[10px] text-gray-500 dark:text-gray-400">* Şifre değiştirmek için eski şifrenizi girmelisiniz.</p>
            </div>
          </div>


          <!-- Kapatma ve Kaydet Butonları -->
          <div class="flex justify-end gap-3">
            <button
              @click="showSettingsModal = false"
              class="px-4 py-2 bg-gray-200 dark:bg-slate-700 text-gray-800 dark:text-slate-200 rounded-lg hover:bg-gray-300 dark:hover:bg-slate-600 transition-colors text-sm font-medium"
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
      <div v-if="showFollowers" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 p-4" @click="showFollowers=false">
        <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl w-80 max-h-[80vh] overflow-y-auto shadow-2xl border border-gray-100 dark:border-slate-700" @click.stop>
          <h2 class="text-xl font-bold mb-5 text-center border-b dark:border-slate-700 pb-3 text-black dark:text-white">Takipçiler</h2>
          <div v-if="followers.length === 0" class="text-center text-gray-500 dark:text-gray-400 py-8 text-sm">Henüz takipçiniz yok.</div>
          <div v-for="f in followers" :key="f.username" class="flex items-center gap-3 p-2.5 hover:bg-gray-100 dark:hover:bg-slate-700/50 rounded-xl cursor-pointer transition-colors group" @click="goToProfile(f.username)">
            <img :src="f.photo || defaultPhoto" class="w-10 h-10 rounded-full object-cover border-2 border-gray-100 dark:border-slate-600 group-hover:border-red-500 transition-all" alt="profile" />
            <span class="font-semibold text-gray-800 dark:text-slate-200 group-hover:text-red-600">{{ f.username }}</span>
          </div>
        </div>
      </div>
    </transition>


    <!--takip etiklerimi görme-->
    <transition name="zoom">
      <div v-if="showFollowing" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 p-4" @click="showFollowing=false">
        <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl w-80 max-h-[80vh] overflow-y-auto shadow-2xl border border-gray-100 dark:border-slate-700" @click.stop>
          <h2 class="text-xl font-bold mb-5 text-center border-b dark:border-slate-700 pb-3 text-black dark:text-white">Takip Edilenler</h2>
          <div v-if="following.length === 0" class="text-center text-gray-500 dark:text-gray-400 py-8 text-sm">Henüz kimseyi takip etmiyorsunuz.</div>
          <div v-for="f in following" :key="f.username" class="flex items-center gap-3 p-2.5 hover:bg-gray-100 dark:hover:bg-slate-700/50 rounded-xl cursor-pointer transition-colors group" @click="goToProfile(f.username)">
            <img :src="f.photo || defaultPhoto" class="w-10 h-10 rounded-full object-cover border-2 border-gray-100 dark:border-slate-600 group-hover:border-red-500 transition-all" alt="profile" />
            <span class="font-semibold text-gray-800 dark:text-slate-200 group-hover:text-red-600">{{ f.username }}</span>
          </div>
        </div>
      </div>
    </transition>


  <!-- Filtreleme Sekmeleri -->
  <div class="w-full max-w-6xl mx-auto px-6 mt-12 flex justify-center">
    <div class="flex flex-wrap justify-center gap-2 bg-gray-100 dark:bg-slate-800/80 p-1.5 rounded-xl border border-gray-200 dark:border-slate-700">
      <button @click="activeTab = 'Tümü'" :class="activeTab === 'Tümü' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Tümü</button>
      <button @click="activeTab = 'Kitap İncelemesi'" :class="activeTab === 'Kitap İncelemesi' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Kitaplar</button>
      <button @click="activeTab = 'Şiir'" :class="activeTab === 'Şiir' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Şiirler</button>
      <button @click="activeTab = 'Duygu & Düşünce'" :class="activeTab === 'Duygu & Düşünce' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Duygular</button>
    </div>
  </div>

  <!-- Kitap kartları -->
  <div class="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full max-w-6xl mx-auto px-6 pb-20">
    <div v-for="book in filteredProfileBooks" :key="book.id" class="relative rounded-2xl shadow-sm hover:shadow-md border border-gray-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-6 transition-all duration-300 flex flex-col">
      <img v-if="book.image" :src="book.image" class="w-full h-48 object-cover mt-2 rounded-xl border border-gray-100 dark:border-slate-700 shadow-sm"  alt="book image" />
      <div class="px-2 py-4 overflow-hidden flex-1 flex flex-col">
        <div class="flex justify-between items-start mb-2 gap-2">
          <h2 class="text-xl font-bold text-black dark:text-white">{{ book.title }}</h2>
          <span class="text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider shrink-0"
                :class="{'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300': book.category === 'Şiir',
                         'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300': book.category === 'Duygu & Düşünce',
                         'bg-gray-100 text-gray-800 dark:bg-slate-700 dark:text-gray-300': book.category === 'Kitap İncelemesi'}">
            {{ book.category }}
          </span>
        </div>
        <p class="text-gray-600 dark:text-slate-300 text-sm mb-5 leading-relaxed" :style="book.expanded ? 'max-height:none' : 'max-height:4rem; overflow:hidden;'">{{ book.description }}</p>
        
        <!-- Repost Kutusu -->
        <div v-if="book.repost_id && book.repost" @click="openOriginalPost(book.repost)" class="mb-5 border border-gray-200 dark:border-slate-700 rounded-xl p-4 bg-gray-50 dark:bg-slate-800/50 cursor-pointer hover:bg-gray-100 dark:hover:bg-slate-700/80 transition-colors">
           <div class="flex items-center gap-2 mb-2 w-max" @click.stop="goToProfile(book.repost.username)">
              <img v-if="book.repost.user_photo" :src="book.repost.user_photo" class="w-5 h-5 rounded-full object-cover">
              <i v-else class="fa-solid fa-circle-user text-xs text-gray-400"></i>
              <span class="text-xs font-semibold text-gray-800 dark:text-slate-200 hover:text-red-600 transition-colors">{{ book.repost.username }}</span>
           </div>
           <div class="flex flex-col sm:flex-row gap-3">
              <img v-if="book.repost.image" :src="book.repost.image" class="w-full sm:w-16 h-16 object-cover rounded-lg border border-gray-200 dark:border-slate-600 shadow-sm">
              <div>
                <h3 class="font-bold text-sm text-black dark:text-white mb-1">{{ book.repost.title }}</h3>
                <p class="text-xs text-gray-600 dark:text-slate-400 line-clamp-2">{{ book.repost.description }}</p>
              </div>
           </div>
        </div>

        <div class="flex flex-wrap gap-2 mt-auto">
          <button class="bg-red-600 hover:bg-red-700 text-white text-sm font-semibold py-2 px-4 rounded-lg shadow-sm transition-colors" @click="book.expanded = !book.expanded">{{ book.expanded ? 'Daha Az Gör' : 'Devamını Oku' }}</button>
          <button v-show="isOwner" @click="openEditBook(book)" class="bg-slate-800 dark:bg-slate-600 text-white hover:bg-slate-900 dark:hover:bg-gray-500 text-sm font-semibold py-2 px-4 rounded-lg shadow-sm transition-colors">Düzenle</button>
          <button v-show="isOwner" @click="deleteBook(book.id)" class="border border-red-600 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-sm rounded-lg px-4 py-2 font-semibold transition-colors">Sil</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Kitap düzenleme modal -->
  <transition name="zoom">
    <div v-if="showEditBook" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showEditBook = false">
      <div class="bg-white dark:bg-slate-800 p-8 rounded-xl w-96 shadow-2xl border border-gray-200 dark:border-slate-700" @click.stop>
        <h2 class="text-xl font-bold mb-6 text-black dark:text-white">Kitabı Düzenle</h2>
        <div class="flex items-center mb-2 font-medium text-gray-800 dark:text-slate-300">
            <i class="fa-regular fa-image"></i>
            <p class="pl-2">Kapak Fotoğrafı</p>
        </div>
        <input type="file" accept="image/*" @change="e=>addBookImage(e,editingBook)" class="mb-5 block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-red-50 file:text-red-700 hover:file:bg-red-100 dark:file:bg-slate-700 dark:file:text-slate-300 transition-colors">
        
        <label class="block mb-2 font-medium text-sm text-gray-800 dark:text-slate-300">Kategori</label>
        <select v-model="editingBook.category" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-5 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500 appearance-none">
          <option value="Kitap İncelemesi">Kitap İncelemesi</option>
          <option value="Şiir">Şiir</option>
          <option value="Duygu & Düşünce">Duygu & Düşünce</option>
        </select>

        <label class="block mb-2 font-medium text-sm text-gray-800 dark:text-slate-300">Başlık</label>
        <input v-model="editingBook.title" type="text" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-5 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500" />
        <label class="block mb-2 font-medium text-sm text-gray-800 dark:text-slate-300">Açıklama</label>
        <textarea v-model="editingBook.description" rows="4" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-6 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-red-500"></textarea>
        <div class="flex justify-end gap-3">
          <button class="px-5 py-2 rounded-lg bg-slate-200 dark:bg-slate-700 text-gray-900 dark:text-slate-200 font-semibold hover:bg-slate-300 dark:hover:bg-slate-600 transition-colors" @click="showEditBook = false">İptal</button>
          <button class="px-5 py-2 rounded-lg bg-red-600 text-white font-semibold hover:bg-red-700 transition-colors shadow-sm" @click="saveEditedBook">Kaydet</button>
        </div>
      </div>
    </div>
  </transition>

  <!-- Orijinal Gönderi Modal -->
  <transition name="zoom">
    <div v-if="showOriginalPostModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 p-4" @click="showOriginalPostModal = false">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl w-full max-w-2xl shadow-2xl border border-gray-200 dark:border-slate-700 max-h-[90vh] overflow-y-auto" @click.stop>
         <div class="flex justify-between items-center mb-6 border-b border-gray-100 dark:border-slate-700 pb-3">
           <h2 class="text-xl font-bold text-black dark:text-white"><i class="fa-solid fa-book-open mr-2 text-red-500"></i>Orijinal Gönderi</h2>
           <button @click="showOriginalPostModal = false" class="text-gray-500 hover:text-red-500 transition-colors"><i class="fa-solid fa-times text-2xl"></i></button>
         </div>
         
         <p @click="goToProfile(viewingOriginalPost.username)" class="flex items-center cursor-pointer font-semibold text-lg mb-5 px-4 py-1.5 inline-flex bg-gray-100 dark:bg-slate-700 text-gray-900 dark:text-slate-200 hover:bg-red-600 dark:hover:bg-red-600 hover:text-white rounded-lg transition-colors w-max">
            <img v-if="viewingOriginalPost.user_photo" :src="viewingOriginalPost.user_photo" class="w-8 h-8 rounded-full mr-2 object-cover border border-white">
            <i v-else class="fa-solid fa-circle-user text-sm mr-2"></i>{{ viewingOriginalPost.username }}
         </p>
         <img v-if="viewingOriginalPost.image" class="w-full h-auto max-h-96 object-cover rounded-xl shadow-sm border border-gray-100 dark:border-slate-700 mb-6" alt="book image" :src="viewingOriginalPost.image" />
         <h2 class="text-2xl font-bold mb-3 text-black dark:text-white">{{ viewingOriginalPost.title }}</h2>
         <p class="text-gray-600 dark:text-slate-300 text-base leading-relaxed whitespace-pre-wrap">{{ viewingOriginalPost.description }}</p>
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
const showOriginalPostModal = ref(false)
const viewingOriginalPost = ref(null)
const showSettingsModal = ref(false)
const editingBook = ref({})
const currentTheme = ref('light')
const profilePrivacy = ref('Herkes')
const activeTab = ref('Tümü')


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


const openOriginalPost = (repost) => {
  viewingOriginalPost.value = repost;
  showOriginalPostModal.value = true;
};

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

// Gönderi ekleme
const addBook = async () => {
  const newBook = { title: 'Yeni Gönderi', description: 'İçerik buraya gelecek...', username: username.value, category: 'Kitap İncelemesi' }
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

const filteredProfileBooks = computed(() => {
  if (activeTab.value === 'Tümü') return books.value
  return books.value.filter(book => book.category === activeTab.value)
})


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
</style>
