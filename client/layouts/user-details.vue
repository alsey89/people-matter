<template>
    <DefaultLayout>
        <div v-auto-animate class="w-full h-full flex flex-col gap-4">
            <!-- !Company Tabs -->
            <app-user-tabs />
            <!-- !Content -->
            <slot />
        </div>
    </DefaultLayout>
</template>

<script setup>
import { useRoute } from 'vue-router';
import DefaultLayout from '~/layouts/default.vue'

const userStore = useUserStore();

const route = useRoute();
const userId = route.params.id;

onBeforeMount(() => {
    userStore.fetchSingleUserData(userId);
});
</script>
