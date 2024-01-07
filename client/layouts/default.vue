<template>
    <div class="w-full h-full min-h-screen inset-0 flex justify-center bg-primary">
        <div class="w-full max-w-[1550px] flex">
            <div class="w-48 bg-red-500">
                SIDEBAR
            </div>
            <div class="w-[80%]">
                <router-view></router-view>
            </div>
        </div>
    </div>
</template>

<script setup>


</script>