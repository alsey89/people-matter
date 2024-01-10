<template>
    <NuxtLink :to="navigateTo" class="px-1 rounded-sm">
        <slot></slot>
    </NuxtLink>
</template>
  
<script setup>
const { navigateTo } = defineProps({
    navigateTo: {
        type: String,
        required: true
    },
});
</script>
  