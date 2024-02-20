<template>
    <div class="border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div v-auto-animate class="px-2 gap-2 flex justify-between items-center p-1 md:p-2">
            <button @click="navigateTo('/company')" class="hover:text-primary">
                <client-only>
                    <div v-if="companyName" class="flex gap-2 font-bold">
                        <AppLogo :src="companyLogoUrl || '/defaultLogo.png'" shape="square" class="w-6 h-6" />
                        <div>{{ companyName }}</div>
                    </div>
                    <div v-else> No Company Data </div>
                </client-only>
            </button>
            <UIPopover position="bottom-left">
                <template v-slot:trigger class="w-10 h-10">
                    <AppAvatar />
                </template>
                <template v-slot:content>
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
            </UIPopover>
        </div>
    </div>
</template>

<script setup>
const companyName = persistedState.sessionStorage.getItem('companyName');
const companyLogoUrl = persistedState.sessionStorage.getItem('companyLogoUrl');

</script>