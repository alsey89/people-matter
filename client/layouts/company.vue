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
                        <AppCompanyHeader />
                        <router-view></router-view>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


