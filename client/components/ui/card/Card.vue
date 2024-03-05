<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div :class="cn(
    'rounded-sm bg-card text-card-foreground border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)]',
    props.class,
  )
    ">
    <slot />
  </div>
</template>
