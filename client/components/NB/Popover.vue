<template>
    <div class="relative">
        <div @click="showPopover = !showPopover">
            <slot name="trigger"></slot>
        </div>
        <div @click="closePopover" v-if="showPopover" class="absolute right-8 top-6">
            <slot name="content"></slot>
        </div>
    </div>
</template>

<script setup>
const showPopover = ref(false)

const closePopover = () => {
    showPopover.value = false
}

</script>


