<template>
    <div class="flex justify-center w-full min-h-screen bg-primary-bg">
        <!-- Container with max-width and centered -->
        <div class="flex w-full max-w-[1550px]">
            <!-- Sidebar -->
            <div class="w-48 bg-red-500">
                SIDEBAR
            </div>

            <!-- Main Content Area -->
            <div class="flex flex-col w-full p-24">
                <!-- Profile Section -->
                <div class="p-4 bg-gray-100 border-2 border-black rounded-md shadow-[2px_2px_0px_rgba(0,0,0,1)]">
                    <div class="flex gap-2 w-full">
                        <!-- User Info and Navigation -->
                        <div class="flex flex-col gap-4 w-[15%]">
                            <div class="flex flex-col items-center w-full">
                                <client-only>{{ userStore.getFullName }}</client-only>
                                <img src="https://via.placeholder.com/100" class="rounded-full" />
                            </div>
                            <AppProfileNavigation />
                        </div>

                        <div class="border border-black"></div>

                        <!-- Router View -->
                        <div class="flex-grow">
                            <router-view></router-view>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
const userStore = useUserStore();

</script>
