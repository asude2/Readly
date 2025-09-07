<template>
    <div class="flex justify-between p-4">
        <div>
            <h1 class="text-red-500 text-3xl font-medium italic font-serif flex  cursor-default">
                Readly
                <i class="fa-solid fa-book-open text-2xl p-2 hover:scale-110 transition-transform transition-colors duration-700 border-2 border-transparent rounded-full"></i>
            </h1>
            <i class="fa-regular fa-circle-left text-3xl p-2 cursor-pointer"></i>
        </div>
        <div>
            <nav class="w-full">
                <ul class="flex items-center tracking-widest justify-center space-x-12">
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      kitaplar
                    </a>
                  </li>
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      okurlar
                    </a>
                  </li>
                  <li>
                    <a href="#" class="relative after:block after:h-px after:w-0 after:bg-black after:transition-all after:duration-500 after:absolute after:bottom-0 after:left-0 hover:after:w-full">
                      yorumlamalar
                    </a>
                  </li>
                </ul>
            </nav>
            <div class="flex items-center border border-gray-300 rounded-full px-3 my-4 w-full">
                <i class="fa-solid fa-magnifying-glass text-gray-500"></i>
                <input type="text" placeholder="arama yap.." class="w-full px-2 py-1 outline-none rounded-full"/>
            </div>
        </div>

        <div @click.prevent="goToProfile">
          <i class="fa-solid fa-circle-user p-3 text-3xl flex cursor-pointer"><p class="px-1 text-2xl">profil</p></i>
        </div>
    </div>
</template>
<script>
export default{
  data(){
    return{}
  },
  methods:{
    goToProfile(){
      this.$router.push('/profil')
    }
  }
}
</script>

