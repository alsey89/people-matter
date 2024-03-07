<!-- todo: FIGURE OUT HEADER RERENDERING ON ROUTE CHANGE -->

<template>
    <div class="border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div class="px-2 gap-2 flex justify-between items-center p-1 md:p-2">
            <button @click="navigateTo('/company')" class="hover:text-primary">
                <div v-if="companyName" class="flex gap-2 font-bold">
                    <AppLogo :src="companyLogoUrl || '/defaultLogo.png'" shape="square" class="w-6 h-6" />
                    <div>{{ companyName }}</div>
                </div>
                <div v-else> No Company Data </div>
            </button>
            <button v-if="isMobileView">
                <Icon v-if="!showSidebar" @click="handleToggleMenuButtonClick" name="solar:list-down-bold"
                    class="w-6 h-6 text-primary font-bold" />
                <Icon v-else @click="handleToggleMenuButtonClick" name="solar:list-up-bold"
                    class="w-6 h-6 text-primary font-bold" />
            </button>
            <AppHeaderAvatarMenu v-else />
        </div>
    </div>
</template>

<script setup>
const companyName = persistedState.sessionStorage.getItem('companyName');
const companyLogoUrl = persistedState.sessionStorage.getItem('companyLogoUrl');

const { isMobileView, showSidebar } = defineProps({
    isMobileView: {
        type: Boolean,
        required: true,
    },
    showSidebar: {
        type: Boolean,
        required: true,

    },
});

const emit = defineEmits(['toggle']);

const handleToggleMenuButtonClick = () => {
    emit('toggle');
};
</script>