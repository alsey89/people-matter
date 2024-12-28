<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuthStore } from '@/stores/Auth';
import posthog from 'posthog-js';

const authStore = useAuthStore();

onMounted(() => {
  if (authStore.isAuthenticated) {
    posthog.identify(authStore.user.id);
    authStore.setIdentified(true); // Set a flag in your auth store to mark as identified
  }
});


</script>