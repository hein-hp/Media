/**
 * 快捷键配置
 */
export interface ShortcutConfig {
  /** 快捷键 (1-9, a-z) */
  key: string
  /** 目标文件夹（绝对路径或相对路径） */
  targetDir: string
  /** 显示名称 */
  label: string
}

/**
 * 移动记录（用于撤销）
 */
export interface MoveRecord {
  /** 原始路径 */
  sourcePath: string
  /** 目标路径 */
  targetPath: string
  /** 操作时间 */
  timestamp: string
}

/**
 * 分类状态
 */
export type ClassifyStatus = 'idle' | 'processing' | 'success' | 'error'

/**
 * 分类统计
 */
export interface ClassifyStats {
  /** 已处理数量 */
  processed: number
  /** 总数量 */
  total: number
  /** 各分类数量 */
  categories: Record<string, number>
}

