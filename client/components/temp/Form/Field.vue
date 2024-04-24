<template>
    <div class="h-full w-full border-2 border-black rounded-sm grid grid-cols-12 shadow-[2px_2px_0px_rgba(0,0,0,1)]">
        <div class="col-span-2 px-2 py-2 flex justify-center items-center">
            <slot name="icon"></slot>
        </div>
        <div class="col-span-10 w-full pr-2 py-2 rounded-nd">
            <slot></slot>
        </div>
    </div>
</template>
