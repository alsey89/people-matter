<template>
    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
        <!-- !Tab Navigation -->
        <div v-auto-animate class="flex justify-start gap-0.5 items-center">
            <div v-for="tab in ['Company', 'Departments', 'Locations']" :key="tab" :class="getTabClasses(tab)"
                class="text-lg font-bold rounded-t-md cursor-pointer shadow-[2px_0px_0px_rgba(0,0,0,1)]"
                @click="setActiveTab(tab)">
                {{ tab }}
            </div>
            <div class="mt-auto flex-grow border-b-2 border-black"></div>
        </div>
        <!-- !Company -->
        <AppCompany v-if="activeTab == 'Company'" />
        <!-- !Department -->
        <AppDepartment v-if="activeTab == 'Departments'" />
        <!-- !Location -->
        <AppLocation v-if="activeTab == 'Locations'" />
    </div>
</template>

<script setup>
const companyStore = useCompanyStore()
// Tabs
const activeTab = ref('Company');

const setActiveTab = (tabName) => {
    if (activeTab.value === tabName) {
        return
    };
    if (companyStore.getCompanyName === null && (tabName !== 'Company')) {
        return
    };
    activeTab.value = tabName;
};

const getTabClasses = (tabName) => {
    let isDisabled = companyStore.getCompanyName === null && (tabName === 'Departments' || tabName === 'Locations');
    return {
        'mt-1 px-4 py-1.5 text-primary border-l border-t border-r border-black': activeTab.value === tabName && !isDisabled,
        'mt-2 px-2 py-1 text-gray-500 border-gray-500 border-l border-t border-r border-b-2': activeTab.value !== tabName && !isDisabled,
        'mt-2 px-2 py-1 text-gray-300 border-gray-500 border-l border-t border-r border-b-2 cursor-not-allowed': isDisabled,
    };
};

//! meta & loading -----------------------------
definePageMeta({
    title: 'Company',
    layout: 'default',
})
const selectedCompanyId = persistedState.sessionStorage.getItem('activeCompanyId')
onBeforeMount(() => {
    if (selectedCompanyId) {
        companyStore.fetchCompanyListAndExpandById(selectedCompanyId)
        return
    }
    companyStore.fetchCompanyListAndExpandDefault()
})

</script>