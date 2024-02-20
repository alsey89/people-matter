<template>
    <UICard id="modal" class="w-full h-full">
        <UICardHeader>
            Confirmation Required
        </UICardHeader>
        <div>
            {{ confirmationModalMessage }}
        </div>
        <UICardFooter>
            <div class="flex justify-end gap-2">
                <UIButtonSquare @click=handleConfirmButtonClick size="md" class="px-4 hover:bg-secondary-bg">
                    Confirm
                </UIButtonSquare>
                <UIButtonSquare @click="handleCancelButtonClick" size="md" class="px-2 hover:bg-secondary-bg">
                    Cancel
                </UIButtonSquare>
            </div>
        </UICardFooter>
    </UICard>
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
// Scroll to bottom of modal when it is mounted
onMounted(() => {
    const modal = document.getElementById("modal")
    modal.scrollIntoView({ behavior: "smooth", block: "end" })
})
</script>