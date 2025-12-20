import {computed, onMounted, onUnmounted, ref, type Ref} from "vue";
import type {MediaInfo} from "@/types";
import {RemoveMedia} from "../../wailsjs/go/app/App";
import {handler} from "../../wailsjs/go/models"

/**
 * 媒体查看器配置项
 *
 * @property step 缩放步长（默认0.25）
 * @property maxScale 最大缩放比例（默认3）
 * @property minScale 最小缩放比例（默认0.5）
 * @property enableKeyboard 是否启用键盘导航（默认true）
 */
export interface MediaViewerOptions {
  step?: number;
  maxScale?: number;
  minScale?: number;
  enableKeyboard?: boolean;
}

/**
 * 媒体查看器 composable 管理媒体预览状态、导航和缩放，支持自定义配置和键盘操作
 *
 * @param mediaList 媒体列表
 * @param options 可选配置项
 * @returns 媒体查看器的状态和操作方法
 */
export function useMediaViewer(
  mediaList: Ref<MediaInfo[]>,
  options: MediaViewerOptions = {}
) {
  // 合并默认配置
  const {
    step = 0.1,
    maxScale = 3,
    minScale = 0.5,
    enableKeyboard = true,
  } = options;

  // 核心响应式状态
  const isOpen = ref(false); // 查看器开关
  const currentIndex = ref(0); // 当前媒体索引
  const scale = ref(1); // 媒体缩放比例
  const isLoading = ref(false); // 当前媒体加载状态

  // 拖拽相关状态
  const offsetX = ref(0); // X轴偏移量
  const offsetY = ref(0); // Y轴偏移量
  const isDragging = ref(false); // 是否正在拖拽

  // 当前预览的媒体信息
  const currentMedia = computed(() => {
    const list = mediaList.value;
    if (
      list.length === 0 ||
      currentIndex.value < 0 ||
      currentIndex.value >= list.length
    ) {
      return null;
    }
    return list[currentIndex.value];
  });

  const hasNext = computed(() => currentIndex.value < mediaList.value.length - 1);
  const hasPrev = computed(() => currentIndex.value > 0);

  /** 重置位置和缩放 */
  const resetTransform = () => {
    scale.value = 1;
    offsetX.value = 0;
    offsetY.value = 0;
  };

  const open = (index = 0) => {
    const listLength = mediaList.value.length;
    if (listLength === 0) return; // 空列表不打开

    // 索引合法性校验
    currentIndex.value = Math.max(0, Math.min(index, listLength - 1));
    resetTransform();
    isOpen.value = true; // 打开查看器
  };

  const close = () => {
    isOpen.value = false;
    resetTransform();
  };

  const remove = async () => {
    const list = mediaList.value;
    const curIdx = currentIndex.value;

    // 列表为空或索引越界则直接返回
    if (list.length === 0 || curIdx < 0 || curIdx >= list.length) {
      return;
    }
    const mediaInfo = mediaList.value[curIdx];
    const media = handler.MediaInfo.createFrom({
        path: mediaInfo.path,
        name: mediaInfo.name,
        size: mediaInfo.size,
        url: mediaInfo.url,
        type: mediaInfo.type
      }
    )
    await RemoveMedia(media)

    // 移除当前索引对应的媒体项
    mediaList.value.splice(curIdx, 1);

    const newListLength = mediaList.value.length;
    if (newListLength === 0) {
      currentIndex.value = 0;
      close();
      return;
    }

    // 若移除的是最后一个元素，索引回退一位
    if (curIdx >= newListLength) {
      currentIndex.value = newListLength - 1;
    }

    resetTransform();
  }

  const next = () => {
    if (hasNext.value) {
      currentIndex.value++;
      resetTransform();
      isLoading.value = true;
    }
  };

  const prev = () => {
    if (hasPrev.value) {
      currentIndex.value--;
      resetTransform();
      isLoading.value = true;
    }
  };

  const zoomIn = () => {
    scale.value = Math.min(scale.value + step, maxScale);
  };

  const zoomOut = () => {
    scale.value = Math.max(scale.value - step, minScale);
  };

  const resetZoom = () => {
    resetTransform();
  };

  const loaded = () => {
    isLoading.value = false;
  };

  /** 开始拖拽 */
  const startDrag = () => {
    if (scale.value > 1) {
      isDragging.value = true;
    }
  };

  /** 结束拖拽 */
  const endDrag = () => {
    isDragging.value = false;
  };

  /** 拖拽移动 */
  const drag = (deltaX: number, deltaY: number) => {
    if (isDragging.value && scale.value > 1) {
      offsetX.value += deltaX;
      offsetY.value += deltaY;
    }
  };

  const handleKeydown = (e: KeyboardEvent) => {
    // 查看器未打开或禁用键盘时不处理
    if (!isOpen.value || !enableKeyboard) return;

    switch (e.key) {
      case "Escape": // ESC关闭
        close();
        break;
      case "ArrowRight": // 右箭头下一张
        next();
        break;
      case "ArrowLeft": // 左箭头上一张
        prev();
        break;
      case "+": // 加号放大
      case "=": // 等号（部分键盘+号需要按shift）
        zoomIn();
        break;
      case "-": // 减号缩小
        zoomOut();
        break;
      case "0": // 数字0重置缩放
        resetZoom();
        break;
      case "d": // d 删除
      case "D":
        remove()
        break
      default:
        break;
    }
  };

  onMounted(() => {
    if (enableKeyboard) {
      document.addEventListener("keydown", handleKeydown);
    }
  });

  onUnmounted(() => {
    if (enableKeyboard) {
      document.removeEventListener("keydown", handleKeydown);
    }
  });

  return {
    isOpen,
    currentIndex,
    currentMedia,
    scale,
    offsetX,
    offsetY,
    isDragging,
    isLoading,
    hasNext,
    hasPrev,
    open,
    close,
    next,
    prev,
    zoomIn,
    zoomOut,
    resetZoom,
    loaded,
    startDrag,
    endDrag,
    drag,
  };
}
