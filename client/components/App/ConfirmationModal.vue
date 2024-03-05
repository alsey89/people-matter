<template>
    <Card id="modal" class="w-full h-full">
        <CardHeader>
            Confirmation Required
        </CardHeader>
        <div>
            {{ confirmationModalMessage }}
        </div>
        <CardFooter>
            <div class="flex justify-end gap-2">
                <Button @click=handleConfirmButtonClick size="md" class="px-4 hover:bg-secondary-bg">
                    Confirm
                </Button>
                <Button @click="handleCancelButtonClick" size="md" class="px-2 hover:bg-secondary-bg">
                    Cancel
                </Button>
            </div>
        </CardFooter>
    </Card>
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