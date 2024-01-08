<template>
    <div :class="`w-full hover:${hoverStyle} px-1 rounded-sm`">
        <slot></slot>
    </div>
</template>
  
<script setup>

const hoverStyle = inject('hoverStyle')
</script>
  