<template>
    <div class="flex flex-col gap-2">
        <div v-auto-animate class="flex justify-start gap-0.5 items-center">
            <client-only>
                <div v-if="activeRoute.startsWith('/user')" @click="handleTabSelect(tab.to)" v-for=" tab  in  tabs "
                    :key="tab.name" :class="getTabClasses(tab.to)"
                    class="min-w-12 md:min-w-24 text-center text-sm md:text-base font-bold rounded-t-md cursor-pointer shadow-[2px_0px_0px_rgba(0,0,0,1)]">
                    {{ tab.name }}
                </div>
                <div class="mt-auto flex-grow border-b-2 border-black"></div>
            </client-only>
        </div>
        <!-- !User Details -->
        <div class="min-h-max flex gap-2 items-center border-b-2 border-black py-2 pr-4">
            <AppImage :src="userStore.getSingleUserAvatarUrl || '/defaultAvatar.jpg'" shape="circle" class="w-8 h-8" />
            <h1 v-if="userStore.getSingleUserFullName" class="text-lg font-bold"> {{ userStore.getSingleUserFullName }}
            </h1>
            <Icon v-else name="svg-spinners:bars-scale-fade" class="h-6 w-6" />
        </div>
    </div>
</template>

<script setup>
import { useRoute } from 'vue-router';
const route = useRoute()

const userStore = useUserStore()

const activeRoute = computed(() => {
    return route.path
})
const userId = route.params.id

//tabs ---------------------
const tabs = computed(() => [
    { name: 'Info', to: `/user/${userId}` },
    { name: 'Job', to: `/user/${userId}/job` },
    { name: 'Salary', to: `/user/${userId}/salary` },
    { name: 'Attendance', to: `/user/${userId}/attendance` },
    { name: 'Leave', to: `/user/${userId}/leave` },
]);
const handleTabSelect = (tabRoute) => {
    if (activeRoute.value == tabRoute) {
        return;
    };
    if (tabRoute.includes('null')) {
        return;
    };
    navigateTo(tabRoute)
};
const getTabClasses = (tabRoute) => {
    let isDisabled = tabRoute.includes('null')
    return {
        // active tab
        'mt-1 px-2 md:px-4 py-1.5 text-primary border-black border-l border-r border-t border-b-transparent': activeRoute.value == tabRoute && !isDisabled,
        // inactive tab
        'mt-2 px-1 md:px-2 py-1 text-gray-500 border-gray-500 border-b-black border-l border-t border-r border-b-2': activeRoute.value != tabRoute && !isDisabled,
        // disabled tab
        'mt-2 px-1 md:px-2 py-1 text-gray-300 border-gray-500 border-b-black border-l border-t border-r border-b-2 hover:cursor-not-allowed': isDisabled,
    };
};

</script>  