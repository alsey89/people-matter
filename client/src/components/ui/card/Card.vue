<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div :class="cn(
    'rounded-none border-2 border-secondary bg-card text-card-foreground',
    // 'rounded-lg border bg-card text-card-foreground shadow-sm',
    props.class,
  )
    ">
    <slot />
  </div>
</template>
