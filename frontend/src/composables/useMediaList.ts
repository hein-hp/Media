import { ref, onMounted, onUnmounted } from "vue";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";
import type { MediaInfo } from "@/types";

/**
 * 媒体列表管理 composable 监听 media-list/close 事件，管理媒体列表状态
 */
export function useMediaList() {
  const mediaList = ref<MediaInfo[]>([]);
  const isClose = ref(false);

  onMounted(() => {
    EventsOn("media-list", (data: MediaInfo[]) => {
      mediaList.value = data;
      isClose.value = false;
    });
    EventsOn("media-close", () => {
      mediaList.value = [];
      isClose.value = true;
    });
  });

  onUnmounted(() => {
    EventsOff("media-list");
    EventsOff("media-close");
  });

  return {
    mediaList,
    isClose,
  };
}
