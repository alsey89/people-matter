<script setup lang="ts">
import { type HTMLAttributes, computed } from 'vue'
import { TabsList, type TabsListProps } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps<TabsListProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <TabsList v-bind="delegatedProps" :class="cn(
    'inline-flex h-10 items-center justify-center rounded-sm bg-card p-1 text-muted-foreground border-black border-2 shadow-[2px_2px_0px_rgba(0,0,0,1)]',
    props.class,
  )">
    <slot />
  </TabsList>
</template>
