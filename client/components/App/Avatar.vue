<template>
    <div class="avatar">
        <img :src="avatarUrl" class="rounded-full w-8 h-8 object-cover border-2 border-black" />
    </div>
</template>

<script setup>
const userStore = useUserStore()

// Computed property to determine the avatar URL
const avatarUrl = computed(() => {
    return userStore.avatarUrl || '/defaultAvatar.jpg';
});
</script>
