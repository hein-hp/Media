import {onMounted, onUnmounted, ref} from "vue";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";
import type {MediaInfo} from "@/types";

/**
 * 媒体列表管理 composable 监听 media-list 事件，管理媒体列表状态
 */
export function useMediaList() {
  const mediaList = ref<MediaInfo[]>([]);
  const mediaCount = ref(0)

  onMounted(() => {
    EventsOn("media-list", (data: MediaInfo[]) => {
      mediaList.value = data;
    });
    EventsOn("media-count", (data: number) => {
      mediaCount.value = data;
    });
  });

  onUnmounted(() => {
    EventsOff("media-list");
    EventsOff("media-count");
  });

  return {
    mediaList,
    mediaCount
  };
}
