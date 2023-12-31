<template>
    <div
        class="w-full h-full flex flex-row justify-center items-center gap-4 p-4 border-black border-2 rounded-md shadow-[2px_2px_0px_rgba(0,0,0,1)]">
        <slot name="header"></slot>
        <div class="flex-grow">
            <slot></slot> <!-- Default slot for main content -->
        </div>
        <slot name="footer"></slot>
    </div>
</template>
