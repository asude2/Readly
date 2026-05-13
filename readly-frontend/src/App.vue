<template>
  <div :class="isDarkMode ? 'dark' : ''">
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

<style>
/* Dark Mode Stilleri */
.dark {
  --bg-primary: #111317;
  --bg-secondary: #272b33;
  --bg-card: #d8d9dd;
  --text-primary: #111827;
  --text-secondary: #4b5563;
  --border-color: #cbd5e1;
  --accent-color: #000000;
}

.dark {
  background-color: #111317;
  color: #e3e3e3;
}

.dark .bg-red-50,
.dark .bg-red-100 {
  background-color: #aeafaf !important;
  color: var(--text-primary) !important;
}

.dark .bg-white,
.dark .bg-gray-200 {
  background-color: var(--bg-card) !important;
  color: var(--text-primary) !important;
}

.dark .text-red-600,
.dark .text-red-700 {
  color: #ff0000 !important;
}

.dark .text-gray-600,
.dark .text-gray-700,
.dark .text-gray-500,
.dark .text-black {
  color: var(--text-primary) !important;
}

.dark input,
.dark textarea {
  background-color: #ffffff;
  color: #111827;
  border-color: #f6f6f600;
}

.dark input::placeholder,
.dark textarea::placeholder {
  color: #000000e9;
}

.dark .border-gray-300,
.dark .border-gray-400,
.dark .border-gray-100 {
  border-color: #cbd5e1;
}

.dark .bg-gray-200 {
  background-color: var(--bg-card) !important;
}

.dark button.bg-black {
  background-color: #1f2937;
  color: #e3e3e3;
}

.dark button.bg-black:hover {
  background-color: #111827;
}

.dark button.bg-gray-200 {
  background-color: #e5e7eb;
  color: #111827;
}

.dark button.bg-red-500 {
  background-color: #b91c1c;
}

.dark button.bg-red-500:hover,
.dark .hover\:bg-red-700:hover {
  background-color: #991b1b;
}

/* Light Mode Stilleri */
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f8f8f8;
  --bg-tertiary: #f0f0f0;
  --text-primary: #000000;
  --text-secondary: #666666;
  --border-color: #e0e0e0;
  --accent-color: #000000;
}

body {
  transition: background-color 0.3s ease, color 0.3s ease;
}
</style>

