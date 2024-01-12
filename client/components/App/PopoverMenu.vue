<template>
    <NBMenu>
        <NBMenuItemNavigation navigateTo="/" class="hover:bg-secondary-bg">
            Home
        </NBMenuItemNavigation>
        <NBMenuItemNavigation navigateTo="/signin" class="hover:bg-secondary-bg">
            Sign In
        </NBMenuItemNavigation>
        <NBMenuItem @click="userStore.signout" class="hover:bg-secondary-bg">
            Signout
        </NBMenuItem>
    </NBMenu>
</template>

<script setup>
const userStore = useUserStore()
</script>