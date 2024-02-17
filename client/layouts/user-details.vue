<template>
    <div class="w-full h-full min-h-screen inset-0 flex justify-center bg-primary-bg">
        <div class="w-full max-w-[1280px] flex flex-col py-2 gap-4">
            <!-- !HEADER -->
            <div class="w-full h-12 z-10 px-2">
                <AppHeader />
            </div>
            <!-- !BODY -->
            <div class="w-full h-full flex gap-4" style="max-height: calc(100vh - 5rem)">
                <!-- !Sidebar -->
                <div class="hidden md:block min-w-48 min-h-full md:pl-2">
                    <AppSidebar />
                </div>
                <div class="flex flex-grow overflow-y-auto px-2">
                    <div v-auto-animate class="w-full h-full flex flex-col gap-4">
                        <!-- !User Details Tabs -->
                        <app-user-header />
                        <!-- !User Details -->
                        <div class="min-h-max flex gap-2 items-center border-b-2 border-black py-2 pr-4">
                            <AppImage :src="userStore.getSingleUserAvatarUrl || 'defaultAvatar.jpg'" shape="circle"
                                class="w-8 h-8" />
                            <h1 class="text-lg font-bold"> {{ userStore.getSingleUserFullName }} </h1>
                        </div>
                        <!-- <NBCard v-if="userStore.getSingleUserData">
                            <div class="flex justify-between items-center p-2">
                                <div class="flex gap-4 items-center">
                                    <AppImage :src="userStore.getSingleUserAvatarUrl || 'defaultAvatar.jpg'" shape="square"
                                        class="w-8 h-8" />
                                    <div class="flex flex-col overflow-x-hidden">
                                        <h1 class="text-sm md:text-base text-nowrap font-bold"> {{
                                            userStore.getSingleUserFullName }}
                                        </h1>
                                        <p class="text-sm md:text-base"> Unassigned </p>
                                    </div>
                                </div>
                            </div>
                        </NBCard>
                        <NBCard v-else>
                            <div v-if="userStore.isLoading">
                                Loading...
                            </div>
                            <div v-else class="m-auto">
                                No Data
                            </div>
                        </NBCard> -->
                        <!-- !Content -->
                        <router-view></router-view>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
const userStore = useUserStore();

</script>

