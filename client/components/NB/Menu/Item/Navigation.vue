<template>
    <NuxtLink :to="navigateTo" class="px-1">
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
  