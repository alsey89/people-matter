<template>
    <DefaultLayout>
        <div v-auto-animate class="w-full h-full flex flex-col gap-4">
            <!-- !Users List Tabs -->
            <app-user-list-tabs />
            <!-- !Content -->
            <router-view></router-view>
        </div>
    </DefaultLayout>
</template>

<script setup>
import DefaultLayout from '~/layouts/default.vue'

const userStore = useUserStore()

onBeforeMount(() => {
    userStore.fetchAllUsersData()
});
</script>
