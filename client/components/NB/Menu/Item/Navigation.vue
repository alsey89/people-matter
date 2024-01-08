<template>
    <NuxtLink :to="navigateTo" :class="`w-full hover:${hoverStyle} px-1 rounded-sm`">
        <slot></slot>
    </NuxtLink>
</template>
  
<script setup>
const props = defineProps({
    navigateTo: {
        type: String,
        required: true
    },
});
const hoverStyle = inject('hoverStyle')
</script>
  