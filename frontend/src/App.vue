<template>
  <div class="min-h-screen bg-white">
    <header class="sticky top-0 z-40 h-12 flex items-center justify-center
                   bg-white/80 backdrop-blur-md border-b border-gray-100"
            style="--wails-draggable: drag">
      <h1 class="text-sm font-medium text-gray-500 select-none">
        {{ mediaList.length > 0 ? `${mediaList.length} 个媒体文件` : 'Media' }}
      </h1>
    </header>

    <!-- 主内容区域 -->
    <main class="pb-6">
      <MediaGrid
        v-if="mediaList.length > 0"
        :images="mediaList"
        @select="viewer.open"
      />
      <EmptyState v-else />
    </main>

    <!-- 媒体预览器 -->
    <MediaViewer
      :is-open="viewer.isOpen.value"
      :media="viewer.currentMedia.value"
      :current-index="viewer.currentIndex.value"
      :total="mediaList.length"
      :scale="viewer.scale.value"
      :offset-x="viewer.offsetX.value"
      :offset-y="viewer.offsetY.value"
      :is-dragging="viewer.isDragging.value"
      :has-next="viewer.hasNext.value"
      :has-prev="viewer.hasPrev.value"
      @close="viewer.close"
      @prev="viewer.prev"
      @next="viewer.next"
      @zoom-in="viewer.zoomIn"
      @zoom-out="viewer.zoomOut"
      @reset-zoom="viewer.resetZoom"
      @start-drag="viewer.startDrag"
      @end-drag="viewer.endDrag"
      @drag="viewer.drag"
    />
  </div>
</template>

<script lang="ts" setup>
import { MediaGrid, MediaViewer, EmptyState } from '@/components'
import { useMediaList, useMediaViewer } from '@/composables'

// 媒体列表状态
const { mediaList } = useMediaList()

// 媒体查看器
const viewer = useMediaViewer(mediaList)
</script>