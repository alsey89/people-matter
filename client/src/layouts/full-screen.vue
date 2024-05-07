<template>
    <div class="flex justify-center items-center min-h-screen bg-background">
        <div class="w-full max-w-md">
            <router-view></router-view>
        </div>
    </div>
</template>

<script>
export default {
    name: 'full-screen-layout',
};
</script>
