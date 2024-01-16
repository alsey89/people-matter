<template>
    <button
        :class="`py-${height} text-${textSize} font-bold border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] active:translate-y-1 active:shadow-none rounded-sm`">
        <slot></slot>
    </button>
</template>

<script setup>

const { size, textSize } = defineProps({
    size: {
        type: String,
        default: "md",
    },
    textSize: {
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
            return 3;
        default:
            return 2;
    }
});

</script>