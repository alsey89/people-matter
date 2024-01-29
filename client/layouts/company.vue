<template>
    <div class="w-full h-full min-h-screen inset-0 flex justify-center bg-primary-bg">
        <div class="w-full max-w-[1280px] flex flex-col py-2 gap-4">
            <!-- Header -->
            <div class="w-full h-12 z-10 px-2">
                <AppHeader />
            </div>
            <!-- Body -->
            <div class="w-full h-full flex gap-4" style="max-height: calc(100vh - 5rem)">
                <!-- Sidebar -->
                <div class="hidden md:block min-w-48 min-h-full md:pl-2">
                    <AppSidebar />
                </div>
                <!-- Content -->
                <div class="flex flex-grow overflow-y-auto px-2">
                    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
                        <div class="flex justify-start gap-0.5 items-center">
                            <client-only>
                                <Button @click="handleTabSelect(tab.to)" v-for="tab in tabs" :key="tab.name"
                                    :class="getTabClasses(tab.to)"
                                    class="text-lg font-bold rounded-t-md cursor-pointer shadow-[2px_0px_0px_rgba(0,0,0,1)]">
                                    {{ tab.name }}
                                </Button>
                                <div class="mt-auto flex-grow border-b-2 border-black"></div>
                            </client-only>
                        </div>
                        <router-view></router-view>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { useRoute } from 'vue-router';
const companyStore = useCompanyStore()

//tabs ---------------------
const route = useRoute()
const activeRoute = computed(() => {
    return route.path
})
const tabs = computed(() => [
    { name: 'Company', to: `/company` },
    { name: 'Departments', to: `/company/${companyStore.getCompanyId}/department` },
    { name: 'Locations', to: `/company/${companyStore.getCompanyId}/location` },
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
        'mt-1 px-4 py-1.5 text-primary border-black border-l border-t border-r': activeRoute.value == tabRoute && !isDisabled,
        // inactive tab
        'mt-2 px-2 py-1 text-gray-500 border-gray-500 border-b-black border-l border-t border-r border-b-2': activeRoute.value != tabRoute && !isDisabled,
        // disabled tab
        'mt-2 px-2 py-1 text-gray-300 border-gray-500 border-b-black border-l border-t border-r border-b-2 hover:cursor-not-allowed': isDisabled,
    };
};

// life cycle related ---------------------
const activeCompanyId = persistedState.sessionStorage.getItem('activeCompanyId')
onBeforeMount(() => {
    if (activeCompanyId) {
        companyStore.fetchCompanyListAndExpandById({ companyId: activeCompanyId })
    } else {
        companyStore.fetchCompanyListAndExpandDefault()
    }
});
</script>
