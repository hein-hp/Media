<template>
  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component"/>
    </transition>
  </router-view>
</template>

<script lang="ts" setup>
import {onMounted, onUnmounted} from 'vue'
import {useRouter} from 'vue-router'
import {EventsOff, EventsOn} from '../wailsjs/runtime'

const router = useRouter()

onMounted(() => {
  // 监听后端发送的路由跳转事件
  EventsOn('router', (route: string) => {
    router.push(route)
  })
})

onUnmounted(() => {
  EventsOff('router')
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
