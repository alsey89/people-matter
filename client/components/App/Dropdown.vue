<template>
    <div class="relative">
        <select v-model="internalSelectedOption"
            class="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500" :required="required">
            <option :value="null">Select an option... </option> <!-- Added unselect option -->
            <option v-for="option in options" :key="option.id" :value="option.ID">
                {{ option.name || option.title }} [ID: {{ option.ID }}]
            </option>
        </select>
    </div>
</template>

<script setup>


const props = defineProps({
    modelValue: {
        type: Number,
        default: null
    },
    options: Array,
    required: Boolean
});

const emit = defineEmits(['update:modelValue']);

// Set the internal selected option based on the modelValue
const internalSelectedOption = ref(props.modelValue);

watch(internalSelectedOption, (newValue) => {
    emit('update:modelValue', newValue);
});

// Watch for external changes in modelValue
watch(() => props.modelValue, (newValue) => {
    internalSelectedOption.value = newValue;
});
</script>
