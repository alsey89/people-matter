<template>
    <div class="flex justify-center items-center min-h-screen bg-primary">
        <div class="w-full max-w-xs">
            <slot></slot>
        </div>
    </div>
</template>

<script>
export default {
    name: 'full-screen-layout',
};
</script>
