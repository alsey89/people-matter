<template>
    <transition name="slide-pop">
        <slot></slot>
    </transition>
</template>


<style>
.slide-pop-enter-active,
.slide-pop-leave-active {
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.slide-pop-enter,
.slide-pop-leave-to {
    opacity: 0;
    transform: translateX(30px);
}
</style>