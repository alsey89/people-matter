<template>
    <button v-if="showPopover" @click="showPopover = false" class="fixed inset-0 z-5"></button>
    <div v-auto-animate class="relative">
        <div @click="showPopover = !showPopover">
            <slot name="trigger"></slot>
        </div>
        <div @click="showPopover = false" v-if="showPopover" :class="`absolute ${popoverPosition}`">
            <slot name="content"></slot>
        </div>
    </div>
</template>

<script setup>
const showPopover = ref(false)

const { position } = defineProps({
    position: {
        type: String,
        default: 'bottom'
    }
})

const popoverPosition = computed(() => {
    switch (position) {
        case 'bottom':
            return "top-10 right-0"
        case 'top':
            return "bottom-10 right-0"
        case 'top-right':
            return "bottom-10 left-10"
        case 'top-left':
            return "bottom-10 right-10"
        case 'bottom-right':
            return "top-10 left-10"
        case 'bottom-left':
            return "right-10 top-10"
        default:
            return "right-0 bottom-0"
    }
})

</script>

