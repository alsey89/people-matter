<template>
    <button
        :class="`py-${height} border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] active:translate-y-1 active:shadow-none rounded-full`">
        <slot></slot>
    </button>
</template>

<script setup>

const { size } = defineProps({
    size: {
        type: String,
        default: "md",
    },
});

const height = computed(() => {
    switch (size) {
        case "sm":
            return 1;
        case "md":
            return 2;
        case "lg":
            return 4;
        default:
            return 2;
    }
});

</script>