<template>
    <UIMenu>
        <UIMenuItemNavigation navigateTo="/" class="hover:bg-secondary-bg">
            Home
        </UIMenuItemNavigation>
        <UIMenuItemNavigation navigateTo="/signin" class="hover:bg-secondary-bg">
            Sign In
        </UIMenuItemNavigation>
        <UIMenuItem @click="userStore.signout" class="hover:bg-secondary-bg">
            Signout
        </UIMenuItem>
    </UIMenu>
</template>

<script setup>
const userStore = useUserStore()
</script>