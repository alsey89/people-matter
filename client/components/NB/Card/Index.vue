<template>
    <div
        class="w-full h-full flex flex-col justify-center gap-4 p-2 border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)]">
        <slot name="header"></slot>
        <div>
            <slot></slot> <!-- Default slot for main content -->
        </div>
        <slot name="footer"></slot>
    </div>
</template>