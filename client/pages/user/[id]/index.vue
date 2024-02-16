<template>
    <div class="grid grid-cols-1 md:grid-cols-12">
        <div class="col-span-1 md:col-span-3 order-1 md:order-2 bg-red-500">
            <div class="flex flex-col gap-2">
                <AppImage shape="circle" :src="userStore.getSingleUserAvatarUrl || 'defaultAvatar.jpg'" />
            </div>
        </div>
        <div class="col-span-1 md:col-span-9 order-2 md:order-1 bg-blue-500">
            user details
        </div>
    </div>
</template>

<script setup>
definePageMeta({
    title: "User Details",
    description: "User details",
    layout: "user-details",
});

const userStore = useUserStore();

const route = useRouter().currentRoute.value
</script>