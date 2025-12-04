/**
 * 路书状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  Roadbook,
  RoadbookDetailResponse,
  Waypoint,
  Tag,
  Comment,
  RoadbookListParams,
  CreateRoadbookRequest,
  UpdateRoadbookRequest,
  AddWaypointRequest,
  CreateCommentRequest,
  PaginatedResponse
} from '../types'
import {
  roadbookApi,
  waypointApi,
  commentApi,
  favoriteApi,
  likeApi,
  tagApi
} from '../api'

export const useRoadbookStore = defineStore('roadbook', () => {
  // ==================== 状态 ====================

  /** 当前路书详情 */
  const currentRoadbook = ref<RoadbookDetailResponse | null>(null)

  /** 路书列表 */
  const roadbookList = ref<Roadbook[]>([])

  /** 列表总数 */
  const total = ref(0)

  /** 当前页码 */
  const currentPage = ref(1)

  /** 每页数量 */
  const pageSize = ref(10)

  /** 加载状态 */
  const loading = ref(false)

  /** 标签列表 */
  const tagList = ref<Tag[]>([])

  /** 评论列表 */
  const commentList = ref<Comment[]>([])

  /** 评论总数 */
  const commentTotal = ref(0)

  /** 编辑中的路书草稿 */
  const draftRoadbook = ref<Partial<CreateRoadbookRequest>>({})

  /** 编辑中的途经点列表 */
  const draftWaypoints = ref<Partial<Waypoint>[]>([])

  // ==================== 计算属性 ====================

  /** 是否有更多数据 */
  const hasMore = computed(() => roadbookList.value.length < total.value)

  /** 当前路书的途经点按天分组 */
  const waypointsByDay = computed(() => {
    if (!currentRoadbook.value?.waypoints) return {}
    const grouped: Record<number, Waypoint[]> = {}
    currentRoadbook.value.waypoints.forEach((wp) => {
      const day = wp.dayIndex || 1
      if (!grouped[day]) {
        grouped[day] = []
      }
      grouped[day].push(wp)
    })
    // 按 sortOrder 排序
    Object.keys(grouped).forEach((day) => {
      grouped[Number(day)].sort((a, b) => a.sortOrder - b.sortOrder)
    })
    return grouped
  })

  /** 当前路书是否已收藏 */
  const isFavorited = computed(() => currentRoadbook.value?.isFavorited ?? false)

  /** 当前路书是否已点赞 */
  const isLiked = computed(() => currentRoadbook.value?.isLiked ?? false)

  // ==================== 路书操作 ====================

  /** 获取路书列表 */
  async function fetchRoadbookList(params?: RoadbookListParams) {
    loading.value = true
    try {
      const res = await roadbookApi.getList(params)
      const data = res as unknown as PaginatedResponse<Roadbook>
      roadbookList.value = data.list
      total.value = data.total
      currentPage.value = data.page
      pageSize.value = data.pageSize
    } finally {
      loading.value = false
    }
  }

  /** 加载更多路书 */
  async function loadMoreRoadbooks(params?: RoadbookListParams) {
    if (!hasMore.value || loading.value) return
    loading.value = true
    try {
      const res = await roadbookApi.getList({
        ...params,
        page: currentPage.value + 1,
        pageSize: pageSize.value
      })
      const data = res as unknown as PaginatedResponse<Roadbook>
      roadbookList.value.push(...data.list)
      total.value = data.total
      currentPage.value = data.page
    } finally {
      loading.value = false
    }
  }

  /** 获取路书详情 */
  async function fetchRoadbookDetail(id: number) {
    loading.value = true
    try {
      const res = await roadbookApi.getDetail(id)
      currentRoadbook.value = res as unknown as RoadbookDetailResponse
    } finally {
      loading.value = false
    }
  }

  /** 创建路书 */
  async function createRoadbook(data: CreateRoadbookRequest) {
    const res = await roadbookApi.create(data)
    return res as unknown as Roadbook
  }

  /** 更新路书 */
  async function updateRoadbook(id: number, data: UpdateRoadbookRequest) {
    const res = await roadbookApi.update(id, data)
    if (currentRoadbook.value?.id === id) {
      Object.assign(currentRoadbook.value, res)
    }
    return res as unknown as Roadbook
  }

  /** 删除路书 */
  async function deleteRoadbook(id: number) {
    await roadbookApi.delete(id)
    roadbookList.value = roadbookList.value.filter((r) => r.id !== id)
    if (currentRoadbook.value?.id === id) {
      currentRoadbook.value = null
    }
  }

  // ==================== 途经点操作 ====================

  /** 添加途经点 */
  async function addWaypoint(roadbookId: number, data: AddWaypointRequest) {
    const res = await waypointApi.add(roadbookId, data)
    const waypoint = res as unknown as Waypoint
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.waypoints = currentRoadbook.value.waypoints || []
      currentRoadbook.value.waypoints.push(waypoint)
    }
    return waypoint
  }

  /** 更新途经点 */
  async function updateWaypoint(
    roadbookId: number,
    waypointId: number,
    data: Partial<AddWaypointRequest>
  ) {
    const res = await waypointApi.update(roadbookId, waypointId, data)
    const waypoint = res as unknown as Waypoint
    if (currentRoadbook.value?.id === roadbookId && currentRoadbook.value.waypoints) {
      const index = currentRoadbook.value.waypoints.findIndex((w) => w.id === waypointId)
      if (index !== -1) {
        currentRoadbook.value.waypoints[index] = waypoint
      }
    }
    return waypoint
  }

  /** 删除途经点 */
  async function deleteWaypoint(roadbookId: number, waypointId: number) {
    await waypointApi.delete(roadbookId, waypointId)
    if (currentRoadbook.value?.id === roadbookId && currentRoadbook.value.waypoints) {
      currentRoadbook.value.waypoints = currentRoadbook.value.waypoints.filter(
        (w) => w.id !== waypointId
      )
    }
  }

  // ==================== 评论操作 ====================

  /** 获取评论列表 */
  async function fetchComments(roadbookId: number, params?: { page?: number; pageSize?: number }) {
    const res = await commentApi.getList(roadbookId, params)
    const data = res as unknown as PaginatedResponse<Comment>
    commentList.value = data.list
    commentTotal.value = data.total
  }

  /** 发表评论 */
  async function createComment(roadbookId: number, data: CreateCommentRequest) {
    const res = await commentApi.create(roadbookId, data)
    const comment = res as unknown as Comment
    commentList.value.unshift(comment)
    commentTotal.value++
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.commentCount++
    }
    return comment
  }

  /** 删除评论 */
  async function deleteComment(roadbookId: number, commentId: number) {
    await commentApi.delete(roadbookId, commentId)
    commentList.value = commentList.value.filter((c) => c.id !== commentId)
    commentTotal.value--
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.commentCount--
    }
  }

  // ==================== 收藏/点赞操作 ====================

  /** 收藏路书 */
  async function favoriteRoadbook(roadbookId: number) {
    await favoriteApi.add(roadbookId)
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.isFavorited = true
      currentRoadbook.value.favoriteCount++
    }
    // 更新列表中的数据
    const item = roadbookList.value.find((r) => r.id === roadbookId)
    if (item) {
      item.favoriteCount++
    }
  }

  /** 取消收藏 */
  async function unfavoriteRoadbook(roadbookId: number) {
    await favoriteApi.remove(roadbookId)
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.isFavorited = false
      currentRoadbook.value.favoriteCount--
    }
    const item = roadbookList.value.find((r) => r.id === roadbookId)
    if (item) {
      item.favoriteCount--
    }
  }

  /** 点赞路书 */
  async function likeRoadbook(roadbookId: number) {
    await likeApi.add(roadbookId)
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.isLiked = true
      currentRoadbook.value.likeCount++
    }
    const item = roadbookList.value.find((r) => r.id === roadbookId)
    if (item) {
      item.likeCount++
    }
  }

  /** 取消点赞 */
  async function unlikeRoadbook(roadbookId: number) {
    await likeApi.remove(roadbookId)
    if (currentRoadbook.value?.id === roadbookId) {
      currentRoadbook.value.isLiked = false
      currentRoadbook.value.likeCount--
    }
    const item = roadbookList.value.find((r) => r.id === roadbookId)
    if (item) {
      item.likeCount--
    }
  }

  // ==================== 标签操作 ====================

  /** 获取标签列表 */
  async function fetchTags(params?: { keyword?: string; isSystem?: boolean }) {
    const res = await tagApi.getList(params)
    tagList.value = res as unknown as Tag[]
  }

  // ==================== 草稿操作 ====================

  /** 设置草稿路书 */
  function setDraftRoadbook(data: Partial<CreateRoadbookRequest>) {
    draftRoadbook.value = data
  }

  /** 添加草稿途经点 */
  function addDraftWaypoint(waypoint: Partial<Waypoint>) {
    draftWaypoints.value.push(waypoint)
  }

  /** 更新草稿途经点 */
  function updateDraftWaypoint(index: number, waypoint: Partial<Waypoint>) {
    if (index >= 0 && index < draftWaypoints.value.length) {
      draftWaypoints.value[index] = { ...draftWaypoints.value[index], ...waypoint }
    }
  }

  /** 删除草稿途经点 */
  function removeDraftWaypoint(index: number) {
    if (index >= 0 && index < draftWaypoints.value.length) {
      draftWaypoints.value.splice(index, 1)
    }
  }

  /** 清空草稿 */
  function clearDraft() {
    draftRoadbook.value = {}
    draftWaypoints.value = []
  }

  // ==================== 重置状态 ====================

  /** 重置所有状态 */
  function $reset() {
    currentRoadbook.value = null
    roadbookList.value = []
    total.value = 0
    currentPage.value = 1
    pageSize.value = 10
    loading.value = false
    tagList.value = []
    commentList.value = []
    commentTotal.value = 0
    draftRoadbook.value = {}
    draftWaypoints.value = []
  }

  return {
    // 状态
    currentRoadbook,
    roadbookList,
    total,
    currentPage,
    pageSize,
    loading,
    tagList,
    commentList,
    commentTotal,
    draftRoadbook,
    draftWaypoints,
    // 计算属性
    hasMore,
    waypointsByDay,
    isFavorited,
    isLiked,
    // 路书操作
    fetchRoadbookList,
    loadMoreRoadbooks,
    fetchRoadbookDetail,
    createRoadbook,
    updateRoadbook,
    deleteRoadbook,
    // 途经点操作
    addWaypoint,
    updateWaypoint,
    deleteWaypoint,
    // 评论操作
    fetchComments,
    createComment,
    deleteComment,
    // 收藏/点赞操作
    favoriteRoadbook,
    unfavoriteRoadbook,
    likeRoadbook,
    unlikeRoadbook,
    // 标签操作
    fetchTags,
    // 草稿操作
    setDraftRoadbook,
    addDraftWaypoint,
    updateDraftWaypoint,
    removeDraftWaypoint,
    clearDraft,
    // 重置
    $reset
  }
})
