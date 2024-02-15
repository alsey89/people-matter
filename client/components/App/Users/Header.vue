<template>
    <div v-auto-animate class="flex justify-start gap-0.5 items-center">
        <client-only>
            <div v-if="activeRoute.startsWith('/user')" @click="handleTabSelect(tab.to)" v-for="tab in tabs" :key="tab.name"
                :class="getTabClasses(tab.to)"
                class="min-w-24 text-center text-base font-bold rounded-t-md cursor-pointer shadow-[2px_0px_0px_rgba(0,0,0,1)]">
                {{ tab.name }}
            </div>
            <div class="mt-auto flex-grow border-b-2 border-black"></div>
        </client-only>
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
    { name: 'Users', to: `/user` },
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
</script>