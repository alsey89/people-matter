<template>
    <div class="w-full h-full border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)] bg-secondary-bg">
        <div class="w-full h-full px-2 gap-2 flex justify-end items-center">
            <!-- <AppAvatar /> -->
            <NBPopover>
                <template v-slot:trigger class="w-8 h-8">
                    <AppAvatar />
                </template>
                <template v-slot:content>
                    <AppPopoverMenu hoverStyle="bg-red-500" />
                </template>
            </NBPopover>
        </div>
    </div>
</template>

<script setup>

</script>