<template>
  <div>
      <!--------------------------------------------HEADER KISMI------------------------------------------------------------------------>
        <div class="sticky top-0 z-50 backdrop-blur-md bg-white/80 dark:bg-gray-900/80 border-b border-gray-200 dark:border-gray-800 transition-colors w-full">
          <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center h-20">
            <!-- ------------logo kısmı-------------- -->
            <div class="flex items-center gap-4">
              <i @click="goLoginPg" class="fa-regular fa-circle-left text-2xl text-gray-500 hover:text-red-600 cursor-pointer transition-colors" title="Çıkış Yap"></i>
              <h1 class="text-red-600 text-3xl font-bold italic font-serif flex items-center cursor-default tracking-wide">
                  Readly
                  <i class="fa-solid fa-book-open text-xl px-2"></i>
              </h1>
            </div>

            <!-- Ortadaki Arama ve Navigasyon -->
            <div class="flex-1 max-w-2xl mx-8 hidden md:block">
              <nav class="w-full mb-3">
                <ul class="text-gray-600 dark:text-gray-400 flex items-center tracking-widest justify-center space-x-10 font-semibold text-xs uppercase">
                  <li>
                    <button @click="kitapAra" class="relative pb-1 hover:text-red-600 dark:hover:text-red-400 transition-colors" :class="{'text-red-600 dark:text-red-400 after:w-full': kitapGör, 'after:w-0': !kitapGör, 'after:bg-red-600 after:absolute after:bottom-0 after:left-0 after:h-[2px] after:transition-all after:duration-300': true}">
                      Keşfet
                    </button>
                  </li>
                  <li>
                    <button @click="takipAra" class="relative pb-1 hover:text-red-600 dark:hover:text-red-400 transition-colors" :class="{'text-red-600 dark:text-red-400 after:w-full': takipGör, 'after:w-0': !takipGör, 'after:bg-red-600 after:absolute after:bottom-0 after:left-0 after:h-[2px] after:transition-all after:duration-300': true}">
                      Takip Ettiklerim
                    </button>
                  </li>
                  <li>
                    <button @click="okurAra" class="relative pb-1 hover:text-red-600 dark:hover:text-red-400 transition-colors" :class="{'text-red-600 dark:text-red-400 after:w-full': okurGör, 'after:w-0': !okurGör, 'after:bg-red-600 after:absolute after:bottom-0 after:left-0 after:h-[2px] after:transition-all after:duration-300': true}">
                      Okurlar
                    </button>
                  </li>
                  <li>
                    <button @click="kitaplarAra" class="relative pb-1 hover:text-red-600 dark:hover:text-red-400 transition-colors" :class="{'text-red-600 dark:text-red-400 after:w-full': kitaplarGör, 'after:w-0': !kitaplarGör, 'after:bg-red-600 after:absolute after:bottom-0 after:left-0 after:h-[2px] after:transition-all after:duration-300': true}">
                      Kitaplar
                    </button>
                  </li>
                </ul>


              </nav>

              <!-- Okur Arama Çubuğu -->
              <div v-if="okurGör" class="transition-all duration-300 absolute w-full max-w-2xl left-1/2 -translate-x-1/2 mt-1">
                <div class="flex items-center bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full px-5 py-2 w-full focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-500/20 shadow-sm transition-all">
                  <i class="fa-solid fa-magnifying-glass text-gray-400 mr-3"></i>
                  <input
                    @keyup.enter="searchUsers"
                    v-model="searchQuery"
                    type="text"
                    placeholder="Kullanıcı ara..."
                    class="w-full bg-transparent text-black dark:text-white placeholder-gray-400 outline-none text-sm"
                  />
                  <button @click="searchUsers" class="bg-red-600 text-white rounded-full hover:bg-red-700 px-5 py-1.5 text-xs font-bold transition-colors shadow-sm ml-2">
                    ARA
                  </button>
                </div>

                <!-- Kullanıcı Listesi -->
                <div v-if="users.length > 0" class="absolute top-full left-0 w-full mt-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-xl max-h-60 overflow-y-auto z-50">
                  <div v-for="user in users.filter(u=>u.username!==username)" :key="user.username"
                      class="flex justify-between items-center p-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 border-b border-gray-100 dark:border-gray-700 last:border-0 transition-colors">
                    <div @click="goUsersProfile(user.username)" class="flex items-center font-semibold text-gray-900 dark:text-gray-200 cursor-pointer">
                      <img v-if="user.photo" :src="user.photo" alt="profile photo" class="w-10 h-10 rounded-full mr-3 object-cover border border-gray-200 dark:border-gray-600">
                      <i v-else class="fa-solid fa-circle-user text-3xl mr-3 text-gray-300 dark:text-gray-500"></i>
                      <span>{{ user.username }}</span>
                    </div>
                    <button
                      @click="follow(user.username, user.isFollowing, user.isPending)"
                      class="border rounded-full px-4 py-1.5 text-xs font-semibold transition-all"
                      :class="user.isFollowing ? 'bg-gray-100 dark:bg-gray-700 border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 hover:border-red-200' : user.isPending ? 'bg-gray-200 dark:bg-gray-600 text-gray-500 border-transparent cursor-wait' : 'bg-red-600 text-white border-red-600 hover:bg-red-700'"
                    >
                      {{ user.isFollowing ? 'Takipten Çık' : user.isPending ? 'İstek Gönderildi' : 'Takip Et' }}
                    </button>
                  </div>
                </div>
              </div>

              <!-- Kitap Yorumu Arama Çubuğu -->
              <div v-if="kitapGör" class="flex items-center bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full px-5 py-2 mt-1 w-full focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-500/20 transition-all shadow-sm">
                <i class="fa-solid fa-magnifying-glass text-gray-400 mr-3"></i>
                <input
                  v-model="searchBookQuery"
                  type="text"
                  placeholder="Kitap incelemesi ara..."
                  class="w-full bg-transparent text-black dark:text-white placeholder-gray-400 outline-none text-sm"
                />
              </div>
              <!-- Kitap Kataloğu Arama Çubuğu -->
              <div v-if="kitaplarGör" class="flex items-center bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full px-5 py-2 mt-1 w-full focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-500/20 transition-all shadow-sm">
                <i class="fa-solid fa-magnifying-glass text-gray-400 mr-3"></i>
                <input
                  v-model="searchCatalogQuery"
                  type="text"
                  placeholder="Katalogda kitap ara..."
                  class="w-full bg-transparent text-black dark:text-white placeholder-gray-400 outline-none text-sm"
                />
              </div>
            </div>


            <!-- Sağ Taraf: İkonlar ve Profil -->
            <div class="flex items-center gap-4">
              <!-- Messages -->
              <button @click="$router.push('/mesajlar')" class="relative p-2 text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-colors">
                <i class="fa-regular fa-paper-plane text-xl"></i>
                <span v-if="unreadMessageCount > 0" class="absolute top-0 right-0 flex items-center justify-center h-4 w-4 bg-red-600 text-white text-[10px] rounded-full border-2 border-white dark:border-gray-900 font-bold">
                  {{ unreadMessageCount > 9 ? '9+' : unreadMessageCount }}
                </span>
              </button>

              <!-- Notifications -->
              <div class="relative notifications-wrapper">
                <button @click="toggleNotifications" class="relative p-2 text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-colors">
                  <i class="fa-regular fa-bell text-xl"></i>
                <span v-if="pendingRequests.length || acceptedFollows.length || likeNotifications.length" class="absolute top-1 right-1 flex h-3 w-3">
                  </span>
                </button>

                <div v-if="isNotificationsOpen"
                    class="absolute right-0 mt-3 w-80 z-50 rounded-2xl bg-white dark:bg-gray-800 p-2 shadow-xl border border-gray-100 dark:border-gray-700">
                  <div class="p-3 border-b border-gray-100 dark:border-gray-700 mb-2">
                    <p class="text-sm font-bold text-black dark:text-white">Bildirimler</p>
                  </div>

                  <div class="max-h-60 overflow-y-auto">
                    <div v-if="pendingRequests.length || acceptedFollows.length || likeNotifications.length">
                      <div v-for="item in pendingRequests" :key="`pending-${item.follower}`"
                          class="flex items-center justify-between gap-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 p-3 mb-1 transition-colors">
                        <div class="text-xs text-gray-800 dark:text-gray-300">
                          <span class="font-bold text-black dark:text-white">{{ item.follower }}</span> seni takip etmek istiyor.
                        </div>
                        <div class="flex gap-2">
                          <button @click="approveRequest(item.follower)"
                                  class="px-3 py-1.5 rounded-lg bg-red-600 text-white hover:bg-red-700 text-[10px] font-semibold shadow-sm">
                            Onayla
                          </button>
                          <button @click="rejectRequest(item.follower)"
                                  class="px-3 py-1.5 rounded-lg bg-gray-200 dark:bg-gray-600 text-gray-800 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-500 text-[10px] font-semibold transition-colors">
                            Sil
                          </button>
                        </div>
                      </div>

                      <div v-for="item in acceptedFollows" :key="`accepted-${item.follower}`"
                          class="rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 p-3 mb-1 text-xs text-gray-600 dark:text-gray-400 transition-colors">
                        <span class="font-bold text-black dark:text-white">{{ item.follower }}</span> seni takip etti.
                      </div>

                      <div v-for="item in likeNotifications" :key="`like-${item.sender}-${item.book_id}`"
                          class="rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 p-3 mb-1 text-xs text-gray-600 dark:text-gray-400 transition-colors">
                        {{ item.message }}
                      </div>
                    </div>

                    <div v-else class="py-8 text-center">
                      <p class="text-sm text-gray-400 dark:text-gray-500">Henüz yeni bir bildirim yok.</p>
                    </div>
                  </div>
                </div>
              </div>

              <div class="w-px h-6 bg-gray-200 dark:bg-gray-700 mx-2 hidden md:block"></div>

              <!-- Profile -->
              <div @click.prevent="goToProfile" class="cursor-pointer flex items-center gap-3 pl-2 group">
                <img v-if="profile" :src="profile" alt="Profil" class="w-10 h-10 rounded-full object-cover ring-2 ring-gray-200 dark:ring-gray-700 group-hover:ring-red-500 transition-all duration-300">
                <p class="text-sm font-semibold text-gray-800 dark:text-gray-200 group-hover:text-red-600 hidden lg:block transition-colors">{{ username }}</p>
              </div>
            </div>
          </div>
        </div>

        <!--------------------------------------------CONTENT KISMI-------------------------------------------------------------------->
        <!-- Filtreleme Sekmeleri -->
        <div v-if="kitapGör || takipGör" class="w-full max-w-5xl mx-auto px-6 mt-8 flex justify-center">
          <div class="flex flex-wrap justify-center gap-2 bg-gray-100 dark:bg-slate-800/80 p-1.5 rounded-xl border border-gray-200 dark:border-slate-700">
            <button @click="activeTab = 'Tümü'" :class="activeTab === 'Tümü' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Tümü</button>
            <button @click="activeTab = 'Kitap İncelemesi'" :class="activeTab === 'Kitap İncelemesi' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Kitaplar</button>
            <button @click="activeTab = 'Şiir'" :class="activeTab === 'Şiir' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Şiirler</button>
            <button @click="activeTab = 'Duygu & Düşünce'" :class="activeTab === 'Duygu & Düşünce' ? 'bg-white dark:bg-slate-700 shadow-sm text-red-600 font-semibold' : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium'" class="px-4 md:px-6 py-2 rounded-lg transition-all text-sm">Duygular</button>
          </div>
        </div>

        <!-- Kitap Kataloğu Görünümü -->
        <div v-if="kitaplarGör" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-12 pb-20 w-full">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            <div v-for="book in filteredCatalogBooks" :key="book.id" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 dark:border-slate-700 overflow-hidden group">
              <div class="relative aspect-[3/4] overflow-hidden">
                <img :src="book.image" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" :alt="book.title" />
                <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-end p-4">
                  <p class="text-white text-xs line-clamp-3">{{ book.description }}</p>
                </div>
              </div>
              <div class="p-5">
                <h3 class="font-bold text-gray-900 dark:text-white mb-1 line-clamp-1">{{ book.title }}</h3>
                <p class="text-gray-500 dark:text-gray-400 text-sm mb-4">{{ book.author }}</p>

                <div class="flex flex-col gap-3 pt-3 border-t border-gray-100 dark:border-slate-700">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-1 text-yellow-500">
                      <i v-for="i in 5" :key="i" class="fa-star" :class="i <= Math.round(book.avg_rating) ? 'fa-solid' : 'fa-regular'"></i>
                    </div>
                    <span class="text-xs font-bold text-gray-600 dark:text-gray-400">{{ book.avg_rating ? book.avg_rating.toFixed(1) : '0.0' }} ({{ book.rating_count }})</span>
                  </div>

                  <div class="flex items-center justify-between">
                    <span class="text-[10px] uppercase tracking-wider font-bold text-gray-400">Senin Puanın:</span>
                    <div class="flex items-center gap-1">
                      <i v-for="i in 5" :key="i"
                         @click="submitRating(book, i)"
                         class="fa-star cursor-pointer transition-colors"
                         :class="[i <= (book.temp_rating || book.user_rating) ? 'fa-solid text-red-500' : 'fa-regular text-gray-300 dark:text-gray-600',
                                 'hover:text-red-400']"
                         @mouseenter="book.temp_rating = i"
                         @mouseleave="book.temp_rating = null"></i>
                    </div>
                  </div>
                </div>
              </div>
              <div class="px-5 pb-5 border-t border-gray-100 dark:border-slate-700">
                <div @click="toggleComments(book)" class="flex items-center cursor-pointer text-gray-600 dark:text-gray-400 hover:text-red-500 dark:hover:text-red-400 transition-colors w-max font-medium">
                  <i class="fa-regular fa-comment mr-2 text-lg"></i>
                  <span>Yorumlar ({{ book.comment_count || 0 }})</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="kitapGör || takipGör" class="flex flex-col justify-center items-center max-w-5xl mx-auto mt-4 pb-20">


          <div v-for="book in filteredBooks" :key="book.id" class="rounded-2xl shadow-sm hover:shadow-md border border-gray-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-6 transition-all duration-300 w-full mb-8 relative justify-center">
            <p @click="goUsersProfile(book.username)" class="flex items-center cursor-pointer font-semibold text-lg mb-6 mt-2 px-4 py-1.5 inline-flex bg-gray-100 dark:bg-slate-700 text-gray-900 dark:text-slate-200 hover:bg-red-600 dark:hover:bg-red-600 hover:text-white rounded-lg transition-colors">
              <img v-if="book.user_photo" :src="book.user_photo" class="w-8 h-8 rounded-full mr-2 object-cover border border-white">
              <i v-else class="fa-solid fa-circle-user text-sm mr-2"></i>{{ book.username }}
            </p>
            <img v-if="book.image" class="w-full sm:w-1/4 h-48 sm:h-56 m-5 object-cover rounded-xl shadow-sm border border-gray-100 dark:border-slate-700" alt="book image" :src="book.image" />
            <div class="px-5 py-4 overflow-hidden w-full">
              <div class="flex justify-between items-start mb-2">
                <h2 class="text-2xl font-bold text-black dark:text-white">{{book.title}}</h2>
                <span v-if="book.category" class="text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider shrink-0"
                      :class="{'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300': book.category === 'Şiir',
                               'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300': book.category === 'Duygu & Düşünce',
                               'bg-gray-100 text-gray-800 dark:bg-slate-700 dark:text-gray-300': book.category === 'Kitap İncelemesi'}">
                  {{ book.category }}
                </span>
              </div>

              <p class="text-gray-600 dark:text-slate-300 text-sm mb-4 leading-relaxed">{{book.description}}</p>

              <!-- Repost Kutusu -->
              <div v-if="book.repost_id && book.repost" @click="openOriginalPost(book.repost)" class="mt-4 border border-gray-200 dark:border-slate-700 rounded-xl p-4 bg-gray-50 dark:bg-slate-800/50 cursor-pointer hover:bg-gray-100 dark:hover:bg-slate-700/80 transition-colors">
                 <p @click.stop="goUsersProfile(book.repost.username)" class="flex items-center cursor-pointer font-semibold text-sm mb-3 w-max">
                    <img v-if="book.repost.user_photo" :src="book.repost.user_photo" class="w-6 h-6 rounded-full mr-2 object-cover border border-white">
                    <i v-else class="fa-solid fa-circle-user text-sm mr-2 text-gray-400"></i>
                    <span class="text-gray-800 dark:text-slate-200 hover:text-red-600 transition-colors">{{ book.repost.username }}</span>
                 </p>
                 <div class="flex flex-col sm:flex-row gap-4">
                    <img v-if="book.repost.image" :src="book.repost.image" class="w-full sm:w-24 h-24 object-cover rounded-lg border border-gray-200 dark:border-slate-600 shadow-sm">
                    <div>
                      <h3 class="font-bold text-black dark:text-white mb-1">{{ book.repost.title }}</h3>
                      <p class="text-gray-600 dark:text-slate-300 text-sm line-clamp-3">{{ book.repost.description }}</p>
                    </div>
                 </div>
              </div>

              <div class="flex items-center gap-4 mt-5">
                <div class="flex items-center gap-2">
                  <button @click="toggleLike(book)" class="p-2 rounded-full hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-200 group">
                    <i class="text-xl transition-transform duration-200 group-hover:scale-110" :class="book.is_liked ? 'fa-solid fa-heart text-red-500' : 'fa-regular fa-heart text-gray-400 dark:text-gray-500 group-hover:text-red-500'"></i>
                  </button>
                  <span class="text-sm font-semibold text-gray-800 dark:text-slate-300">{{ book.like_count }}</span>
                </div>

                <button @click="openRepostModal(book)" class="p-2 rounded-full hover:bg-green-50 dark:hover:bg-green-900/20 transition-colors duration-200 group text-gray-400 dark:text-gray-500 hover:text-green-500">
                  <i class="fa-solid fa-retweet text-xl transition-transform duration-200 group-hover:scale-110"></i>
                </button>
              </div>

              <div class="mt-6 border-t border-gray-100 dark:border-slate-700 pt-4">
                <div @click="toggleComments(book)" class="flex items-center cursor-pointer text-gray-600 dark:text-gray-400 hover:text-red-500 dark:hover:text-red-400 transition-colors w-max font-medium">
                  <i class="fa-regular fa-comment mr-2 text-lg"></i>
                  <span>Yorumlar ({{ book.comment_count || 0 }})</span>
                </div>
              </div>

            </div>
          </div>
        </div>

        <!-- Yorumlar Modal -->
        <transition name="zoom">
          <div v-if="selectedCommentBook" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50 p-4" @click="closeCommentModal">
            <div class="bg-white dark:bg-slate-800 rounded-3xl w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-2xl border border-gray-200 dark:border-slate-700" @click.stop>
              <div class="flex items-start justify-between gap-4 p-5 border-b border-gray-200 dark:border-slate-700">
                <div>
                  <h2 class="text-2xl font-bold text-black dark:text-white">Yorumlar</h2>
                  <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ selectedCommentBook.title }}</p>
                </div>
                <button @click="closeCommentModal" class="text-gray-500 hover:text-red-500 transition-colors"><i class="fa-solid fa-times text-2xl"></i></button>
              </div>
              <div class="p-5">
                <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">Toplam: <span class="font-semibold text-gray-900 dark:text-white">{{ selectedCommentBook.comment_count || 0 }}</span> yorum</p>
                <div v-if="selectedCommentBook.comments.length === 0" class="rounded-2xl bg-gray-50 dark:bg-slate-700 p-4 text-sm text-gray-500 dark:text-gray-400 mb-4">
                  Henüz yorum yok. İlk yorumu sen yap!
                </div>
                <div v-for="c in selectedCommentBook.comments" :key="c.id" class="mb-4 text-sm border-b border-gray-200 dark:border-slate-600 pb-4 last:pb-0 last:border-0">
                  <div class="flex justify-between items-start gap-3">
                    <div class="flex-1">
                      <p class="font-semibold text-red-600 dark:text-red-400 mb-1">{{ c.username }}</p>
                      <p class="text-gray-800 dark:text-slate-200 leading-relaxed">{{ c.content }}</p>
                    </div>
                    <button v-if="c.username === username" @click="deleteComment(selectedCommentBook, c.id)" class="text-gray-400 hover:text-red-600 dark:hover:text-red-400 p-2 rounded-full transition-colors">
                      <i class="fa-solid fa-trash"></i>
                    </button>
                  </div>
                </div>
                <div class="mt-4 flex flex-col gap-3">
                  <textarea v-model="selectedCommentBook.newComment" rows="3" placeholder="Yorumunu yaz..." class="w-full border border-slate-300 dark:border-slate-600 rounded-3xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-red-500 dark:bg-slate-700 dark:text-white dark:placeholder-slate-400 transition-shadow"></textarea>
                  <div class="flex justify-end">
                    <button @click="submitComment(selectedCommentBook)" class="bg-red-600 text-white px-6 py-2 rounded-full text-sm font-semibold hover:bg-red-700 shadow-sm transition-all">Gönder</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </transition>

        <!-- Alıntı Modal -->
        <transition name="zoom">
          <div v-if="showRepostModal" class="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm z-50" @click="showRepostModal = false">
            <div class="bg-white dark:bg-slate-800 p-8 rounded-xl w-full max-w-lg shadow-2xl border border-gray-200 dark:border-slate-700" @click.stop>
              <h2 class="text-xl font-bold mb-6 text-black dark:text-white"><i class="fa-solid fa-retweet mr-2 text-green-500"></i>Alıntıla</h2>

              <textarea v-model="repostText" rows="3" placeholder="Bu gönderi hakkında ne düşünüyorsun?" class="w-full border border-gray-300 dark:border-slate-600 rounded-lg px-3 py-2 mb-4 bg-gray-50 dark:bg-slate-700/50 text-black dark:text-white focus:outline-none focus:ring-2 focus:ring-green-500"></textarea>

              <!-- Alıntılanan post önizlemesi -->
              <div v-if="repostingBook" class="border border-gray-200 dark:border-slate-700 rounded-xl p-4 bg-gray-50 dark:bg-slate-800/50 mb-6">
                 <div class="flex items-center gap-2 mb-2">
                    <img v-if="repostingBook.user_photo" :src="repostingBook.user_photo" class="w-5 h-5 rounded-full object-cover">
                    <i v-else class="fa-solid fa-circle-user text-xs text-gray-400"></i>
                    <span class="text-xs font-semibold text-gray-800 dark:text-slate-200">{{ repostingBook.username }}</span>
                 </div>
                 <h3 class="font-bold text-sm text-black dark:text-white mb-1">{{ repostingBook.title }}</h3>
                 <p class="text-xs text-gray-600 dark:text-slate-400 line-clamp-2">{{ repostingBook.description }}</p>
              </div>

              <div class="flex justify-end gap-3">
                <button class="px-5 py-2 rounded-lg bg-slate-200 dark:bg-slate-700 text-gray-900 dark:text-slate-200 font-semibold hover:bg-slate-300 dark:hover:bg-slate-600 transition-colors" @click="showRepostModal = false">İptal</button>
                <button class="px-5 py-2 rounded-lg bg-green-500 text-white font-semibold hover:bg-green-600 transition-colors shadow-sm flex items-center" @click="submitRepost">
                   Paylaş
                </button>
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

               <p @click="goUsersProfile(viewingOriginalPost.username)" class="flex items-center cursor-pointer font-semibold text-lg mb-5 px-4 py-1.5 inline-flex bg-gray-100 dark:bg-slate-700 text-gray-900 dark:text-slate-200 hover:bg-red-600 dark:hover:bg-red-600 hover:text-white rounded-lg transition-colors w-max">
                  <img v-if="viewingOriginalPost.user_photo" :src="viewingOriginalPost.user_photo" class="w-8 h-8 rounded-full mr-2 object-cover border border-white">
                  <i v-else class="fa-solid fa-circle-user text-sm mr-2"></i>{{ viewingOriginalPost.username }}
               </p>
               <img v-if="viewingOriginalPost.image" class="w-full h-auto max-h-96 object-cover rounded-xl shadow-sm border border-gray-100 dark:border-slate-700 mb-6" alt="book image" :src="viewingOriginalPost.image" />
               <h2 class="text-2xl font-bold mb-3 text-black dark:text-white">{{ viewingOriginalPost.title }}</h2>
               <p class="text-gray-600 dark:text-slate-300 text-base leading-relaxed whitespace-pre-wrap">{{ viewingOriginalPost.description }}</p>
            </div>
          </div>
        </transition>

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
const takipGör=ref(false)
const kitaplarGör=ref(false)
const books=ref([])
const catalogBooks=ref([])
const searchBookQuery = ref("")
const searchCatalogQuery = ref("")
const activeTab = ref('Tümü')

const pendingRequests = ref([])
const acceptedFollows = ref([])
const likeNotifications = ref([])

const showRepostModal = ref(false)
const repostingBook = ref(null)
const repostText = ref('')

const showOriginalPostModal = ref(false)
const viewingOriginalPost = ref(null)
const selectedCommentBook = ref(null)

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
  let result = books.value

  if (activeTab.value !== 'Tümü') {
    result = result.filter(book => book.category === activeTab.value)
  }

  if (searchBookQuery.value) {
    result = result.filter(book =>
      book.title.toLowerCase().includes(searchBookQuery.value.toLowerCase())
    )
  }

  return result
})

const filteredCatalogBooks = computed(() => {
  if (!searchCatalogQuery.value) return catalogBooks.value
  return catalogBooks.value.filter(book =>
    book.title.toLowerCase().includes(searchCatalogQuery.value.toLowerCase()) ||
    book.author.toLowerCase().includes(searchCatalogQuery.value.toLowerCase())
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
      likeNotifications.value = data.likes || []
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
  takipGör.value=false
  kitaplarGör.value=false
  okurGör.value=true
}
const kitapAra=()=>{
  okurGör.value=false
  takipGör.value=false
  kitaplarGör.value=false
  kitapGör.value=true
  fetchAllBooks()
}
const takipAra=()=>{
  okurGör.value=false
  kitapGör.value=false
  kitaplarGör.value=false
  takipGör.value=true
  fetchFollowingBooks()
}
const kitaplarAra=()=>{
  okurGör.value=false
  kitapGör.value=false
  takipGör.value=false
  kitaplarGör.value=true
  fetchCatalogBooks()
}

const fetchCatalogBooks = async () => {
  try {
    const res = await fetch(`http://localhost:8000/api/catalog/books?currentUser=${username.value}`)
    const data = await res.json()
    if (res.ok) {
      catalogBooks.value = data.map(b => ({
        ...b,
        temp_rating: null,
        showComments: false,
        comments: [],
        newComment: '',
        comment_count: b.comment_count || 0,
        catalog_book_id: b.id
      }))
    }
  } catch (err) {
    console.error('Katalog kitapları fetch hatası:', err)
  }
}

const submitRating = async (book, rating) => {
  try {
    const res = await fetch('http://localhost:8000/api/catalog/rate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value,
        book_id: book.id,
        rating: rating
      })
    })
    if (res.ok) {
      // Refresh only this book's data or all catalog
      fetchCatalogBooks()
      await fetchNotifications()
    }
  } catch (err) {
    console.error('Puanlama hatası:', err)
  }
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

const fetchFollowingBooks= async () => {
  try{
    const res= await fetch(`http://localhost:8000/api/getFollowingBooks?currentUser=${username.value}`);
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
  const isSameBook = selectedCommentBook.value && selectedCommentBook.value.id === book.id && (!!selectedCommentBook.value.catalog_book_id === !!book.catalog_book_id)
  if (isSameBook) {
    selectedCommentBook.value = null
    return
  }

  selectedCommentBook.value = book
  if (!book.comments || book.comments.length === 0) {
    try {
      const param = book.catalog_book_id ? `catalog_book_id=${book.id}` : `book_id=${book.id}`;
      const res = await fetch(`http://localhost:8000/api/getComments?${param}`);
      const data = await res.json();
      book.comments = data || [];
    } catch (err) { console.error('Yorum yükleme hatası:', err); }
  }
};

const closeCommentModal = () => {
  selectedCommentBook.value = null
};

const submitComment = async (book) => {
  if (!book.newComment.trim()) return;

  try {
    const payload = {
      username: username.value,
      content: book.newComment
    };
    if (book.catalog_book_id) payload.catalog_book_id = book.id;
    else payload.book_id = book.id;

    const res = await fetch('http://localhost:8000/api/addComment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });
    if (res.ok) {
      const saveComment = await res.json();
      book.comments.push(saveComment); //listeye yeni yorumu ekle
      book.comment_count = (book.comment_count || 0) + 1
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

const openRepostModal = (book) => {
  repostingBook.value = book.repost_id && book.repost ? book.repost : book; // Eğer zaten bir alıntıysa asıl postu alıntıla
  repostText.value = '';
  showRepostModal.value = true;
};

const submitRepost = async () => {
  if (!repostingBook.value) return;

  const newPost = {
    title: "Alıntı",
    description: repostText.value || "Bu gönderiyi paylaştı.",
    image: "",
    username: username.value,
    repost_id: repostingBook.value.id || repostingBook.value.repost_id
  };

  try {
    const res = await fetch('http://localhost:8000/api/addBook', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newPost)
    });

    if (res.ok) {
      showRepostModal.value = false;
      if (takipGör.value) {
         await fetchFollowingBooks();
      } else {
         await fetchAllBooks();
      }
    } else {
      const error = await res.json();
      console.error('Repost hatası:', error.error);
      alert('Paylaşılamadı: ' + error.error);
    }
  } catch (err) {
    console.error('Repost gönderme hatası:', err);
    alert('Bir hata oluştu.');
  }
};

const openOriginalPost = (repost) => {
  viewingOriginalPost.value = repost;
  showOriginalPostModal.value = true;
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

