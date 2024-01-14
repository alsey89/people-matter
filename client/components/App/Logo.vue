<template>
    <div class="w-full h-auto">
        <img :src="src" alt="Logo">
    </div>
</template>
  
<script setup>
const { src } = defineProps({
    src: {
        type: String,
        required: true
    }
})
</script>
  
  