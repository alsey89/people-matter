<template>
    <div v-html="renderedMarkdown" class="markdown-content"></div>
</template>
  
<script setup>
import { ref, watchEffect } from 'vue';
import MarkdownIt from 'markdown-it';
import { defineProps } from 'vue';

const props = defineProps({
    content: String
});

const md = new MarkdownIt({
    breaks: true // Convert '\n' in paragraphs into <br>
});
const renderedMarkdown = ref('');

watchEffect(() => {
    renderedMarkdown.value = md.render(props.content || '');
});
</script>