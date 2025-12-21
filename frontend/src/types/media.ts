/**
 * 媒体类型
 */
export type MediaType = 'image' | 'video'

/**
 * 媒体文件信息接口
 * 对应后端 handler.MediaInfo 结构
 */
export interface MediaInfo {
  /** 文件绝对路径 */
  path: string
  /** 文件名 */
  name: string
  /** 文件大小（字节） */
  size: number
  /** HTTP 访问 URL */
  url: string
  /** 媒体类型 */
  type: MediaType
}

/**
 * 判断是否为视频类型
 */
export function isVideo(media: MediaInfo): boolean {
  return media.type === 'video'
}

/**
 * 判断是否为图片类型
 */
export function isImage(media: MediaInfo): boolean {
  return media.type === 'image'
}

/**
 * 相似图片信息
 */
export interface SimilarImage {
  /** 文件绝对路径 */
  path: string
  /** 文件名 */
  name: string
  /** HTTP 访问 URL */
  url: string
  /** 文件大小（字节） */
  size: number
  /** 修改时间 */
  modTime: string
}

/**
 * 相似图片组
 */
export interface SimilarityResult {
  /** 组 ID */
  groupId: number
  /** 该组中的相似图片 */
  images: SimilarImage[]
}
