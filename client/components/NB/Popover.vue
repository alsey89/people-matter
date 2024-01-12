<template>
    <button v-if="showPopover" @click="showPopover = false" class="fixed inset-0 z-5"></button>
    <div class="relative">
        <div @click="showPopover = !showPopover">
            <slot name="trigger"></slot>
        </div>
        <NBAnimationPop>
            <div @click="showPopover = false" v-if="showPopover" class="absolute right-8 top-6">
                <slot name="content"></slot>
            </div>
        </NBAnimationPop>
    </div>
</template>

<script setup>
const showPopover = ref(false)

</script>

