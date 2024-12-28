<template>
    <div class="w-full h-screen flex bg-background">
        <!-- button to toggle sidebar -->
        <Icon v-if="isMobile && !showSidebar" @click="toggleSidebar" icon="material-symbols:menu-rounded"
            class="fixed top-3 left-3 w-6 h-6 text-primary z-10" />
        <!-- overlay to hide sidebar -->
        <div v-if="showSidebar && isMobile" @click="toggleSidebar" class="fixed inset-0 bg-black bg-opacity-40 z-40">
        </div>
        <!-- sidebar -->
        <Sidebar :showSidebar="showSidebar" :isMobile="isMobile" @closeSidebar="toggleSidebar" />
        <!-- main content -->
        <div :class="{ 'ml-[250px]': showSidebar && !isMobile, 'ml-0': !showSidebar || isMobile }"
            class="w-full h-full overflow-y-auto">
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { Sidebar } from "@/components/app";
import { Icon } from '@iconify/vue';

const showSidebar = ref(window.innerWidth >= 768);
const isMobile = ref(window.innerWidth < 768);

const toggleSidebar = () => {
    showSidebar.value = !showSidebar.value;
};

const handleResize = () => {
    const width = window.innerWidth;
    isMobile.value = width < 768;
    showSidebar.value = width >= 768;
};

onMounted(() => {
    window.addEventListener('resize', handleResize);
    handleResize();
});

onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize);
});
</script>
