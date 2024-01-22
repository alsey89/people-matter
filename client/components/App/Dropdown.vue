<template>
    <div class="relative">
        <select v-model="selectedOption"
            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500" :required="required">
            <option disabled value="">Select an option...</option>
            <option v-for="option in options" :key="option.id" :value="option">
                {{ option.name }} [ID: {{ option.ID }}]
            </option>
        </select>
    </div>
</template>

<script setup>
const selectedOption = ref(null);
const { required, options } = defineProps({
    options: {
        type: Array,
        default: () => [],
    },
    required: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['selected']);
watch(selectedOption, (newSelection) => {
    emit('selected', newSelection);
});

</script>