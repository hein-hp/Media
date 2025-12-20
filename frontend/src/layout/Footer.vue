<template>
  <footer class="fixed bottom-0 left-0 right-0 z-40 h-8 flex items-center justify-between px-4
                 bg-white/80 backdrop-blur-md border-t border-gray-100">
    <!-- 左侧信息 -->
    <div class="flex items-center gap-2 text-xs text-gray-400 select-none">
      <span v-if="status" class="flex items-center gap-1">
        <span class="w-1.5 h-1.5 rounded-full" :class="statusColor"></span>
        {{ status }}
      </span>
    </div>

    <!-- 插槽 -->
    <div class="flex-1 flex items-center justify-center text-xs text-gray-400 select-none">
      <slot/>
    </div>

    <!-- 右侧版本信息 -->
    <div class="text-xs text-gray-400 select-none">
      {{ version }}
    </div>
  </footer>
</template>

<script lang="ts" setup>
import {computed} from 'vue'

const props = withDefaults(defineProps<{
  version?: string
  status?: string
  statusType?: 'success' | 'warning' | 'error' | 'info'
}>(), {
  version: '1.0.0',
  status: '',
  statusType: 'info'
})

const statusColor = computed(() => {
  const colors = {
    success: 'bg-green-500',
    warning: 'bg-yellow-500',
    error: 'bg-red-500',
    info: 'bg-blue-500'
  }
  return colors[props.statusType]
})
</script>
