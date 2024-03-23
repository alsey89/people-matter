<template>
    <div class="w-full h-full min-h-screen inset-0 flex justify-center bg-primary-bg">
        <div class="w-full max-w-[1280px] flex flex-col py-2 gap-1 md:gap-4">
            <!-- !HEADER -->
            <div class="w-full h-12 z-10 px-2">
                <AppHeader :isMobileView="isMobileView" :showSidebar="showSidebar" @toggle="handleSidebarToggle" />
            </div>
            <!-- !BODY -->
            <div v-auto-animate class="w-full h-full flex gap-4" style="max-height: calc(100vh - 5rem)">
                <!-- !Sidebar -->
                <div v-if="showSidebar" class="min-w-full min-h-max md:min-w-48 px-2 md:px-0 md:pl-2">
                    <AppSidebar />
                </div>
                <!-- !Content -->
                <div v-if="showContent" class="flex flex-grow overflow-y-auto px-2">
                    <slot />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { debounce } from 'lodash-es'

const showSidebar = ref(true)
const showContent = ref(true)
const isMobileView = ref(false)

const handleSidebarToggle = () => {
    showSidebar.value = !showSidebar.value
    showContent.value = !showContent.value
}

const debouncedHandleResize = debounce(() => {
    if (window.innerWidth < 768) {
        showSidebar.value = false
        showContent.value = true
        isMobileView.value = true
    } else {
        showSidebar.value = true
        showContent.value = true
        isMobileView.value = false
    }
}, 200) // 200ms debounce time

onMounted(() => {
    debouncedHandleResize()
    window.addEventListener('resize', debouncedHandleResize)
})
onUnmounted(() => {
    window.removeEventListener('resize', debouncedHandleResize)
})
</script>
