<template>
  <div class="h-screen bg-gray-50 dark:bg-slate-900 flex flex-col transition-colors duration-300">
    <!-- HEADER -->
    <div class="bg-white dark:bg-slate-800 border-b border-gray-200 dark:border-slate-700 shadow-sm px-6 py-4 flex justify-between items-center z-10 relative transition-colors">
      <h1 @click="goHome" class="text-red-600 text-3xl font-medium italic font-serif flex items-center cursor-pointer hover:opacity-80 transition-opacity">
        Readly
        <i class="fa-solid fa-book-open text-2xl px-2"></i>
      </h1>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-slate-100">Mesajlar</h2>
      <i @click="goHome" class="fa-regular fa-circle-left text-3xl text-gray-600 dark:text-gray-400 cursor-pointer hover:text-red-600 transition-colors"></i>
    </div>

    <!-- MAIN CHAT AREA -->
    <div class="flex-1 flex overflow-hidden">
      
      <!-- SIDEBAR: Conversations List -->
      <div class="w-1/3 bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-slate-700 flex flex-col transition-colors">
        <div class="p-4 border-b border-gray-200 dark:border-slate-700 flex justify-between items-center bg-gray-50 dark:bg-slate-800/50">
          <h3 class="font-semibold text-gray-800 dark:text-slate-200">Sohbetler</h3>
          <button @click="showNewChatModal = true" class="bg-red-600 text-white p-2 rounded-full hover:bg-red-700 transition-colors shadow-sm" title="Yeni Mesaj">
            <i class="fa-solid fa-pen-to-square"></i>
          </button>
        </div>

        <div class="flex-1 overflow-y-auto">
          <div v-if="conversations.length === 0" class="p-4 text-center text-gray-500 dark:text-gray-400 text-sm mt-10">
            Henüz bir sohbetiniz yok. Sağ üstten yeni bir sohbet başlatın.
          </div>
          <div v-for="conv in conversations" :key="conv.id" 
               @click="selectConversation(conv)"
               class="p-4 border-b border-gray-100 dark:border-slate-700/50 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors flex items-center gap-3"
               :class="{'bg-red-50 dark:bg-red-900/20': activeConversation?.id === conv.id}">
            
            <div class="relative">
              <img v-if="conv.type === 'direct' && conv.otherPhoto" :src="conv.otherPhoto" class="w-12 h-12 rounded-full object-cover border" />
              <div v-else-if="conv.type === 'direct'" class="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center border text-gray-500">
                <i class="fa-solid fa-user"></i>
              </div>
              <div v-else class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center border border-red-200 text-red-600">
                <i class="fa-solid fa-users"></i>
              </div>
              <div v-if="conv.unreadCount > 0" class="absolute -top-1 -right-1 bg-red-600 text-white text-[10px] w-5 h-5 flex items-center justify-center rounded-full font-bold border-2 border-white">
                {{ conv.unreadCount > 9 ? '9+' : conv.unreadCount }}
              </div>
            </div>

            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center">
                <h4 class="font-semibold text-sm truncate">{{ conv.type === 'group' ? conv.name : conv.otherUsername }}</h4>
                <span class="text-[10px] text-gray-400">{{ formatTime(conv.lastMessageAt) }}</span>
              </div>
              <p class="text-xs text-gray-500 truncate mt-0.5" :class="{'font-bold text-gray-800': conv.unreadCount > 0}">
                {{ conv.lastMessage || 'Mesajlaşmaya başlayın...' }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- CHAT VIEW -->
      <div class="flex-1 flex flex-col bg-gray-50 dark:bg-slate-900/50 relative transition-colors">
        <div v-if="!activeConversation" class="absolute inset-0 flex items-center justify-center flex-col text-gray-400 dark:text-gray-500">
          <i class="fa-regular fa-comments text-6xl mb-4"></i>
          <p class="text-lg">Bir sohbet seçin veya yeni bir mesaj başlatın.</p>
        </div>

        <template v-else>
          <!-- Chat Header -->
          <div class="px-6 py-4 bg-white dark:bg-slate-800 border-b border-gray-200 dark:border-slate-700 flex items-center gap-4 z-10 shadow-sm cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors" @click="handleHeaderClick">
            <img v-if="activeConversation.type === 'direct' && activeConversation.otherPhoto" :src="activeConversation.otherPhoto" class="w-10 h-10 rounded-full object-cover border" />
            <div v-else-if="activeConversation.type === 'direct'" class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center border">
               <i class="fa-solid fa-user text-gray-500"></i>
            </div>
            <div v-else class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center border text-red-600">
               <i class="fa-solid fa-users"></i>
            </div>
            <div>
              <h3 class="font-semibold hover:text-red-600 transition-colors">{{ activeConversation.type === 'group' ? activeConversation.name : activeConversation.otherUsername }}</h3>
              <p class="text-xs text-gray-500">{{ activeConversation.type === 'group' ? 'Grup Sohbeti - Üyeleri gör' : 'Özel Mesaj - Profile git' }}</p>
            </div>
          </div>

          <!-- Messages -->
          <div class="flex-1 overflow-y-auto p-6 flex flex-col gap-3" ref="messagesContainer">
             <div v-if="messages.length === 0" class="text-center text-gray-400 dark:text-gray-500 my-auto text-sm">
                İlk mesajı siz gönderin!
             </div>
             <div v-for="msg in messages" :key="msg.id" 
                  class="max-w-[70%] rounded-2xl px-4 py-2 flex flex-col"
                  :class="msg.senderUsername === currentUser ? 'bg-red-600 text-white self-end rounded-tr-sm shadow-sm' : 'bg-white dark:bg-slate-800 border border-gray-200 dark:border-slate-700 text-gray-900 dark:text-slate-100 self-start rounded-tl-sm shadow-sm'">
                <span v-if="activeConversation.type === 'group' && msg.senderUsername !== currentUser" 
                      @click="goToProfile(msg.senderUsername)"
                      class="text-[10px] font-bold text-red-700 mb-0.5 cursor-pointer hover:underline">{{ msg.senderUsername }}</span>
                <p class="text-sm break-words">{{ msg.content }}</p>
                <span class="text-[10px] self-end mt-1 opacity-70">{{ formatTime(msg.createdAt) }}</span>
             </div>
          </div>

          <!-- Message Input -->
          <div class="p-4 bg-white dark:bg-slate-800 border-t border-gray-200 dark:border-slate-700 flex items-center gap-3 transition-colors">
             <input type="text" v-model="newMessage" @keyup.enter="sendMessage"
                    placeholder="Bir mesaj yazın..." 
                    class="flex-1 bg-gray-100 dark:bg-slate-700 border-transparent rounded-full px-5 py-2.5 text-sm text-black dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:bg-white dark:focus:bg-slate-800 focus:border-red-400 focus:ring-1 focus:ring-red-400 border transition-colors" />
             <button @click="sendMessage" :disabled="!newMessage.trim()"
                     class="bg-red-600 text-white w-10 h-10 rounded-full flex items-center justify-center hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
               <i class="fa-solid fa-paper-plane text-sm"></i>
             </button>
          </div>
        </template>
      </div>

    </div>

    <!-- Group Details Modal -->
    <div v-if="showGroupModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
       <div class="bg-white rounded-2xl w-full max-w-md overflow-hidden shadow-xl flex flex-col max-h-[80vh]">
          <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
             <h3 class="font-bold text-lg">{{ activeConversation?.name }} - Üyeler</h3>
             <i @click="showGroupModal = false" class="fa-solid fa-xmark text-xl cursor-pointer hover:text-red-600 transition-colors"></i>
          </div>
          
          <div class="p-6 overflow-y-auto flex-1">
             <!-- Add new member -->
             <div class="flex gap-2 mb-6">
               <input type="text" v-model="newGroupMember" @keyup.enter="addGroupMember" placeholder="Yeni kullanıcı adı" class="flex-1 border px-3 py-2 rounded-lg text-sm focus:outline-none focus:border-red-500" />
               <button @click="addGroupMember" class="bg-red-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-red-700 transition-colors">Ekle</button>
             </div>

             <!-- Members List -->
             <div class="space-y-3">
               <div v-for="member in groupMembersList" :key="member.username" class="flex items-center justify-between p-2 hover:bg-gray-50 rounded-lg border border-transparent hover:border-gray-100 transition-all">
                 <div class="flex items-center gap-3 cursor-pointer" @click="goToProfile(member.username)">
                   <img v-if="member.photo" :src="member.photo" class="w-8 h-8 rounded-full object-cover border border-gray-300" />
                   <i v-else class="fa-solid fa-circle-user text-2xl text-gray-400"></i>
                   <span class="font-medium text-sm text-gray-800 hover:text-red-600">{{ member.username }}</span>
                 </div>
                 <button v-if="member.username !== currentUser" @click="removeGroupMember(member.username)" class="text-gray-400 hover:text-red-600 p-2 transition-colors" title="Gruptan Çıkar">
                   <i class="fa-solid fa-user-minus"></i>
                 </button>
                 <span v-else class="text-xs text-gray-400 font-medium italic">Sen</span>
               </div>
             </div>
          </div>
       </div>
    </div>

    <!-- New Chat Modal -->
    <div v-if="showNewChatModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
       <div class="bg-white rounded-2xl w-full max-w-md overflow-hidden shadow-xl">
          <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
             <h3 class="font-bold text-lg">Yeni Sohbet</h3>
             <i @click="showNewChatModal = false" class="fa-solid fa-xmark text-xl cursor-pointer hover:text-red-600 transition-colors"></i>
          </div>
          
          <div class="p-6">
             <div class="flex gap-2 mb-6 bg-gray-100 p-1 rounded-lg">
                <button @click="chatMode = 'direct'" class="flex-1 py-1.5 text-sm font-medium rounded-md transition-all" :class="chatMode === 'direct' ? 'bg-white shadow text-red-600' : 'text-gray-500 hover:text-gray-700'">Kişi</button>
                <button @click="chatMode = 'group'" class="flex-1 py-1.5 text-sm font-medium rounded-md transition-all" :class="chatMode === 'group' ? 'bg-white shadow text-red-600' : 'text-gray-500 hover:text-gray-700'">Grup</button>
             </div>

             <!-- Direct Chat -->
             <div v-if="chatMode === 'direct'" class="space-y-4">
                <div>
                   <label class="block text-sm font-semibold mb-1 text-gray-700">Kullanıcı Adı</label>
                   <input type="text" v-model="targetUser" placeholder="Kiminle konuşmak istersiniz?" class="w-full border px-4 py-2 rounded-lg focus:outline-none focus:border-red-500 text-sm" />
                </div>
                <button @click="startDirectChat" class="w-full bg-red-600 text-white py-2.5 rounded-lg font-semibold hover:bg-red-700 transition-colors">Sohbet Başlat</button>
             </div>

             <!-- Group Chat -->
             <div v-else class="space-y-4">
                <div>
                   <label class="block text-sm font-semibold mb-1 text-gray-700">Grup Adı</label>
                   <input type="text" v-model="groupName" placeholder="Gruba bir isim verin" class="w-full border px-4 py-2 rounded-lg focus:outline-none focus:border-red-500 text-sm" />
                </div>
                <div>
                   <label class="block text-sm font-semibold mb-1 text-gray-700">Üyeler (Virgülle ayırın)</label>
                   <input type="text" v-model="groupMembersInput" placeholder="örn: jonsnow, arya, sansa" class="w-full border px-4 py-2 rounded-lg focus:outline-none focus:border-red-500 text-sm" />
                </div>
                <button @click="startGroupChat" class="w-full bg-red-600 text-white py-2.5 rounded-lg font-semibold hover:bg-red-700 transition-colors">Grup Oluştur</button>
             </div>
          </div>
       </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const currentUser = localStorage.getItem('username') || ''

const conversations = ref([])
const activeConversation = ref(null)
const messages = ref([])
const newMessage = ref('')
const messagesContainer = ref(null)

const showNewChatModal = ref(false)
const chatMode = ref('direct')
const targetUser = ref('')
const groupName = ref('')
const groupMembersInput = ref('')

const showGroupModal = ref(false)
const groupMembersList = ref([])
const newGroupMember = ref('')

let pollingInterval = null

const goToProfile = (username) => {
  router.push(`/profil/${username}`)
}

const handleHeaderClick = () => {
  if (activeConversation.value.type === 'direct') {
    goToProfile(activeConversation.value.otherUsername)
  } else if (activeConversation.value.type === 'group') {
    fetchGroupMembers()
    showGroupModal.value = true
  }
}

const fetchGroupMembers = async () => {
  if (!activeConversation.value) return
  try {
    const res = await fetch(`http://localhost:8000/api/messages/groupMembers?conversationId=${activeConversation.value.id}`)
    if (res.ok) {
      groupMembersList.value = await res.json() || []
    }
  } catch (err) {
    console.error('Members fetch error:', err)
  }
}

const addGroupMember = async () => {
  if (!newGroupMember.value.trim() || !activeConversation.value) return
  try {
    const res = await fetch('http://localhost:8000/api/messages/addGroupMember', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        conversationId: activeConversation.value.id,
        username: newGroupMember.value.trim()
      })
    })
    if (res.ok) {
      newGroupMember.value = ''
      await fetchGroupMembers()
    } else {
      const text = await res.text()
      alert('Kişi eklenemedi: ' + text)
    }
  } catch (err) {
    console.error('Add member error:', err)
  }
}

const removeGroupMember = async (usernameToRemove) => {
  if (!confirm(`${usernameToRemove} adlı kullanıcıyı gruptan çıkarmak istediğinize emin misiniz?`)) return
  try {
    const res = await fetch('http://localhost:8000/api/messages/removeGroupMember', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        conversationId: activeConversation.value.id,
        username: usernameToRemove
      })
    })
    if (res.ok) {
      await fetchGroupMembers()
    } else {
      alert('Kişi çıkarılamadı.')
    }
  } catch (err) {
    console.error('Remove member error:', err)
  }
}

const goHome = () => {
  router.push('/anasayfa')
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const fetchConversations = async () => {
  if (!currentUser) return
  try {
    const res = await fetch(`http://localhost:8000/api/messages/conversations?currentUser=${currentUser}`)
    if (res.ok) {
      conversations.value = await res.json() || []
    }
  } catch (err) {
    console.error('Conversations error:', err)
  }
}

const fetchMessages = async () => {
  if (!activeConversation.value) return
  try {
    const res = await fetch(`http://localhost:8000/api/messages/history?conversationId=${activeConversation.value.id}`)
    if (res.ok) {
      const msgs = await res.json() || []
      const isNewMessage = msgs.length > messages.value.length
      messages.value = msgs
      
      if (isNewMessage) {
        scrollToBottom()
      }
    }
  } catch (err) {
    console.error('Messages error:', err)
  }
}

const markAsRead = async () => {
  if (!activeConversation.value || !currentUser) return
  try {
    await fetch('http://localhost:8000/api/messages/markRead', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ conversationId: activeConversation.value.id, currentUser })
    })
    // Update local count
    const conv = conversations.value.find(c => c.id === activeConversation.value.id)
    if (conv) conv.unreadCount = 0
  } catch (err) {
    console.error('Mark read error:', err)
  }
}

const selectConversation = async (conv) => {
  activeConversation.value = conv
  messages.value = []
  await fetchMessages()
  await markAsRead()
  scrollToBottom()
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || !activeConversation.value) return
  const text = newMessage.value.trim()
  newMessage.value = ''
  
  try {
    const res = await fetch('http://localhost:8000/api/messages/send', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        conversationId: activeConversation.value.id,
        senderUsername: currentUser,
        content: text
      })
    })
    if (res.ok) {
      await fetchMessages()
      await fetchConversations()
    }
  } catch (err) {
    console.error('Send error:', err)
  }
}

const startDirectChat = async () => {
  if (!targetUser.value.trim()) return
  try {
    const res = await fetch('http://localhost:8000/api/messages/startDirect', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ currentUser, targetUser: targetUser.value.trim() })
    })
    if (res.ok) {
      showNewChatModal.value = false
      targetUser.value = ''
      await fetchConversations()
      const data = await res.json()
      const conv = conversations.value.find(c => c.id === data.conversationId)
      if (conv) selectConversation(conv)
    } else {
      alert('Sohbet başlatılamadı.')
    }
  } catch (err) {
    console.error(err)
  }
}

const startGroupChat = async () => {
  if (!groupName.value.trim() || !groupMembersInput.value.trim()) return
  const members = groupMembersInput.value.split(',').map(m => m.trim()).filter(m => m)
  if (members.length === 0) return

  try {
    const res = await fetch('http://localhost:8000/api/messages/createGroup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: groupName.value.trim(), currentUser, members })
    })
    if (res.ok) {
      showNewChatModal.value = false
      groupName.value = ''
      groupMembersInput.value = ''
      await fetchConversations()
      const data = await res.json()
      const conv = conversations.value.find(c => c.id === data.conversationId)
      if (conv) selectConversation(conv)
    } else {
      alert('Grup oluşturulamadı.')
    }
  } catch (err) {
    console.error(err)
  }
}

const polling = async () => {
  await fetchConversations()
  if (activeConversation.value) {
    await fetchMessages()
  }
}

onMounted(() => {
  if (!currentUser) {
    router.push('/')
    return
  }
  fetchConversations()
  // Poll every 3 seconds for real-time feel
  pollingInterval = setInterval(polling, 3000)
})

onUnmounted(() => {
  if (pollingInterval) clearInterval(pollingInterval)
})

</script>

<style scoped>
/* Scrollbar saklama / özelleştirme */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-track {
  background: transparent; 
}
::-webkit-scrollbar-thumb {
  background: #cbd5e1; 
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #94a3b8; 
}
</style>
