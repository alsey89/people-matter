<template>
    <div class="w-full h-full min-h-screen inset-0 flex justify-center bg-primary-bg">
        <div class="w-full max-w-[1280px] flex flex-col gap-4 p-2">
            <div class="w-full h-12">
                <AppHeader />
            </div>
            <div class="w-full h-full flex gap-4" style="max-height: calc(100vh - 5rem)">
                <div class="w-32 min-h-full">
                    <AppSidebar />
                </div>
                <div class="flex flex-grow overflow-y-auto">
                    <router-view></router-view>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
// Your script here
</script>
