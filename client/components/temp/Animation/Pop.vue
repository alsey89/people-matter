<template>
    <transition name="pop" appear>
        <slot> </slot>
    </transition>
</template>

<style>
.pop-enter-active,
.pop-leave-active {
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.pop-enter,
.pop-leave-to {
    opacity: 0;
    transform: scale(0.5);
}
</style>