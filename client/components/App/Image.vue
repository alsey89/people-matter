<template>
    <img :src="src" alt="Logo" :class="{
        'object-cover': true,
        'rounded-full': shape === 'circle',
        'rounded-md': shape === 'square'
    }">
</template>

<script setup>
const { src, shape } = defineProps({
    src: {
        type: String,
        required: true
    },
    shape: {
        type: String,
        default: 'square'
    }
});
</script>

  
  
  