<template>
    <NBMenu>
        <NBMenuItemNavigation navigateTo="/" class="hover:bg-primary-light">
            Home
        </NBMenuItemNavigation>
        <NBMenuItemNavigation navigateTo="/signin" class="hover:bg-primary-light">
            Sign In
        </NBMenuItemNavigation>
        <NBMenuItem @click="userStore.signout" class="hover:bg-primary-light">
            Signout
        </NBMenuItem>
    </NBMenu>
</template>

<script setup>
const userStore = useUserStore()
</script>