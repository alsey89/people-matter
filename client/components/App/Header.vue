<template>
    <div class="w-full h-full border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div class="w-full h-full px-2 gap-2 flex justify-end items-center">
            <!-- <AppAvatar /> -->
            <NBPopover position="bottom-left">
                <template v-slot:trigger class="w-10 h-10">
                    <AppAvatar />
                </template>
                <template v-slot:content>
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
            </NBPopover>
        </div>
    </div>
</template>

<script setup>
const userStore = useUserStore()
</script>