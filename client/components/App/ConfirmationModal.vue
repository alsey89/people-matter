<template>
    <NBCard class="w-full h-full">
        <NBCardHeader>
            Confirmation Required
        </NBCardHeader>
        <div>
            {{ confirmationModalMessage }}
        </div>
        <NBCardFooter>
            <div class="flex justify-end gap-2">
                <NBButtonSquare @click=handleConfirmButtonClick size="md" class="px-4 hover:bg-secondary-bg">
                    Confirm
                </NBButtonSquare>
                <NBButtonSquare @click="handleCancelButtonClick" size="md" class="px-2 hover:bg-secondary-bg">
                    Cancel
                </NBButtonSquare>
            </div>
        </NBCardFooter>
    </NBCard>
</template>

<script setup>
const { confirmationModalMessage } = defineProps({
    confirmationModalMessage: String,
    default: "Are you sure?"
})
const emit = defineEmits(["confirm", "cancel"])
const handleConfirmButtonClick = () => {
    emit("confirm")
}
const handleCancelButtonClick = () => {
    emit("cancel")
}
</script>