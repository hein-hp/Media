<template>
  <div class="min-h-screen bg-white">
    <Header>{{ mediaCount > 0 ? `共 ${mediaCount} 个文件，扫描出 ${mediaList.length} 个媒体文件` : '媒体预览' }}</Header>

    <!-- 主内容区域 -->
    <main class="pb-12">
      <MediaGrid
        v-if="mediaList.length > 0"
        :images="mediaList"
        @select="viewer.open"
      />
      <EmptyState v-else/>
    </main>

    <!-- 底部栏 -->
    <Footer
      :status="mediaList.length > 0 ? '已就绪' : '先选择文件夹'"
      :status-type="mediaList.length > 0 ? 'success' : 'info'">{{ selectedDir }}
    </Footer>

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
import {EmptyState, MediaGrid, MediaViewer} from '@/components'
import {useMediaList, useMediaViewer, useSelectedDir} from '@/composables'
import {Footer, Header} from '@/layout'

// 选中文件夹
const {selectedDir} = useSelectedDir()

// 媒体列表状态
const {mediaList, mediaCount} = useMediaList()

// 媒体查看器
const viewer = useMediaViewer(mediaList)
</script>
