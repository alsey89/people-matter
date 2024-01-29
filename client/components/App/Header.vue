<template>
    <div class="border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div class="px-2 gap-2 flex justify-between items-center p-1 md:p-2">
            <button @click="navigateTo('/company')" class="hover:text-primary">
                <ClientOnly>
                    <div v-if="companyStore.getCompanyName && companyStore.getCompanyLogoUrl" class="flex gap-2 font-bold">
                        <AppLogo :src="companyStore.getCompanyLogoUrl" shape="square" class="w-6 h-6" />
                        <div>{{ companyStore.getCompanyName }}</div>
                    </div>
                    <Icon v-else name="svg-spinners:bars-fade" class="w-12 h-6 text-primary font-bold" />
                </ClientOnly>
            </button>
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
const companyStore = useCompanyStore()

const selectedCompanyId = persistedState.sessionStorage.getItem('activeCompanyId')

onBeforeMount(() => {
    if (!selectedCompanyId) {
        companyStore.fetchCompanyListAndExpandDefault()
    }
    if (!companyStore.getCompanyName || !companyStore.getCompanyLogoUrl) {
        companyStore.fetchCompanyDetails(selectedCompanyId)
    }
})

</script>