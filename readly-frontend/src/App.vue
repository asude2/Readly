<template>
  <div class="min-h-screen bg-white text-black dark:bg-slate-900 dark:text-slate-100 transition-colors duration-300">
    <router-view></router-view>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'

const isDarkMode = ref(false)

// Tema bilgisini localStorage'dan yükle
const loadTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'light'
  isDarkMode.value = savedTheme === 'dark'
  applyTheme()
}

// Temayı document'e uygula
const applyTheme = () => {
  const htmlElement = document.documentElement
  if (isDarkMode.value) {
    htmlElement.classList.add('dark')
  } else {
    htmlElement.classList.remove('dark')
  }
}

// Tema değişimini izle ve localStorage'a kaydet
watch(isDarkMode, (newVal) => {
  localStorage.setItem('theme', newVal ? 'dark' : 'light')
  applyTheme()
})

onMounted(() => {
  loadTheme()
})

// Global window'a tema değişim fonksiyonu ekle
window.toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
}
window.getTheme = () => isDarkMode.value ? 'dark' : 'light'
</script>

