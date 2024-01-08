<template>
    <nav
        class="w-full flex flex-col gap-1 bg-white p-2 border-black border-2 rounded-md shadow-[2px_2px_0px_rgba(0,0,0,1)]">
        <slot></slot>
    </nav>
</template>
  
<script setup>

const { hoverStyle } = defineProps({
    hoverStyle: {
        type: String,
        default: ""
    }
});

provide('hoverStyle', hoverStyle)

</script>
  
 
  