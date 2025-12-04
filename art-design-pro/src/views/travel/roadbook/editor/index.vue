<template>
  <div class="page-content roadbook-editor">
    <!-- 顶部操作栏 -->
    <div class="editor-header">
      <div class="header-left">
        <ElButton :icon="ArrowLeft" @click="handleBack">返回</ElButton>
        <h2 class="page-title">{{ isEdit ? '编辑路书' : '创建路书' }}</h2>
      </div>
      <div class="header-right">
        <ElButton @click="handleSaveDraft" :loading="saving">保存草稿</ElButton>
        <ElButton type="primary" @click="handlePublish" :loading="saving">
          {{ isEdit ? '更新发布' : '发布路书' }}
        </ElButton>
      </div>
    </div>

    <!-- 编辑器主体 -->
    <div class="editor-body">
      <ElRow :gutter="20">
        <!-- 左侧：基本信息和途经点 -->
        <ElCol :lg="14" :md="24">
          <div class="editor-left">
            <!-- 基本信息表单 -->
            <ElCard class="info-card" shadow="never">
              <template #header>
                <div class="card-header">
                  <span>基本信息</span>
                </div>
              </template>
              <ElForm
                ref="formRef"
                :model="formData"
                :rules="formRules"
                label-width="100px"
                label-position="top"
              >
                <!-- 封面图片上传 -->
                <ElFormItem label="封面图片" prop="coverUrl">
                  <div class="cover-upload">
                    <ElUpload
                      class="cover-uploader"
                      :show-file-list="false"
                      :before-upload="beforeCoverUpload"
                      :http-request="handleCoverUpload"
                      accept="image/*"
                    >
                      <div v-if="formData.coverUrl" class="cover-preview">
                        <ElImage :src="formData.coverUrl" fit="cover" />
                        <div class="cover-mask">
                          <ElIcon><Plus /></ElIcon>
                          <span>更换封面</span>
                        </div>
                      </div>
                      <div v-else class="cover-placeholder">
                        <ElIcon :size="40"><Plus /></ElIcon>
                        <span>上传封面图片</span>
                        <span class="tip">建议尺寸 800x500，支持 JPG、PNG</span>
                      </div>
                    </ElUpload>
                  </div>
                </ElFormItem>

                <!-- 标题 -->
                <ElFormItem label="路书标题" prop="title">
                  <ElInput
                    v-model="formData.title"
                    placeholder="请输入路书标题"
                    maxlength="100"
                    show-word-limit
                  />
                </ElFormItem>

                <!-- 描述 -->
                <ElFormItem label="路书描述" prop="description">
                  <ElInput
                    v-model="formData.description"
                    type="textarea"
                    placeholder="请输入路书描述"
                    :rows="4"
                    maxlength="500"
                    show-word-limit
                  />
                </ElFormItem>

                <!-- 日期范围 -->
                <ElFormItem label="行程日期">
                  <ElDatePicker
                    v-model="dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="YYYY-MM-DD"
                    :disabled-date="disabledDate"
                    style="width: 100%"
                  />
                </ElFormItem>

                <ElRow :gutter="16">
                  <!-- 出行方式 -->
                  <ElCol :span="12">
                    <ElFormItem label="出行方式">
                      <ElSelect v-model="formData.travelMode" style="width: 100%">
                        <ElOption label="自驾" value="driving" />
                        <ElOption label="步行" value="walking" />
                        <ElOption label="骑行" value="cycling" />
                        <ElOption label="公交" value="transit" />
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>

                  <!-- 预算 -->
                  <ElCol :span="12">
                    <ElFormItem label="预计预算">
                      <ElInputNumber
                        v-model="formData.totalBudget"
                        :min="0"
                        :precision="2"
                        placeholder="预算金额"
                        style="width: 100%"
                      >
                        <template #suffix>元</template>
                      </ElInputNumber>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <!-- 可见性 -->
                <ElFormItem label="可见性">
                  <ElRadioGroup v-model="formData.visibility">
                    <ElRadio :value="1">公开</ElRadio>
                    <ElRadio :value="2">私有</ElRadio>
                  </ElRadioGroup>
                </ElFormItem>

                <!-- 标签选择 -->
                <ElFormItem label="标签">
                  <TagSelector v-model="formData.tagIds" />
                </ElFormItem>
              </ElForm>
            </ElCard>

            <!-- 途经点列表 -->
            <ElCard class="waypoints-card" shadow="never">
              <template #header>
                <div class="card-header">
                  <span>途经点 ({{ waypoints.length }})</span>
                  <ElButton type="primary" link :icon="Plus" @click="handleAddWaypoint">
                    添加途经点
                  </ElButton>
                </div>
              </template>
              <WaypointList
                v-model="waypoints"
                :days="dayCount"
                @edit="handleEditWaypoint"
                @delete="handleDeleteWaypoint"
                @reorder="handleReorderWaypoints"
              />
            </ElCard>

            <!-- 日程时间线（仅在有多天行程时显示） -->
            <ElCard v-if="dayCount > 1 && waypoints.length > 0" class="timeline-card" shadow="never">
              <template #header>
                <div class="card-header">
                  <span>日程预览</span>
                  <span class="header-tip">按天查看行程安排</span>
                </div>
              </template>
              <ScheduleTimeline
                :waypoints="waypoints"
                :days="dayCount"
                :start-date="formData.startDate"
              />
            </ElCard>
          </div>
        </ElCol>

        <!-- 右侧：地图预览 -->
        <ElCol :lg="10" :md="24">
          <div class="editor-right">
            <ElCard class="map-card" shadow="never">
              <template #header>
                <div class="card-header">
                  <span>地图预览</span>
                  <div class="map-actions">
                    <ElButton
                      type="primary"
                      link
                      :icon="Refresh"
                      @click="handleRefreshRoute"
                      :loading="routeLoading"
                    >
                      刷新路线
                    </ElButton>
                  </div>
                </div>
              </template>
              <div class="map-container">
                <AMapContainer
                  ref="mapContainerRef"
                  :center="mapCenter"
                  :zoom="12"
                  @ready="handleMapReady"
                />
              </div>
              <!-- 路线信息 -->
              <div v-if="routeInfo" class="route-info">
                <div class="info-item">
                  <span class="label">总距离</span>
                  <span class="value">{{ routeInfo.distance }}</span>
                </div>
                <div class="info-item">
                  <span class="label">预计时间</span>
                  <span class="value">{{ routeInfo.duration }}</span>
                </div>
              </div>
            </ElCard>
          </div>
        </ElCol>
      </ElRow>
    </div>

    <!-- 途经点编辑弹窗 -->
    <WaypointDialog
      v-model:visible="waypointDialogVisible"
      :waypoint="editingWaypoint"
      :days="dayCount"
      @save="handleSaveWaypoint"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, shallowRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadRequestOptions } from 'element-plus'
import { ArrowLeft, Plus, Refresh } from '@element-plus/icons-vue'
import AMapContainer from '../../components/AMapContainer.vue'
import TagSelector from '../../components/TagSelector.vue'
import WaypointList from '../../components/WaypointList.vue'
import WaypointDialog from '../../components/WaypointDialog.vue'
import ScheduleTimeline from '../../components/ScheduleTimeline.vue'
import { useRoadbookStore } from '../../store/roadbook'
import { useRoute as useRouteHook } from '../../composables/useRoute'
import { getAMap } from '../../composables/useAMap'
import type { Waypoint, CreateRoadbookRequest, TravelMode } from '../../types'
import { RoadbookVisibility, RoadbookStatus, WaypointType } from '../../types'

defineOptions({
  name: 'RoadbookEditor'
})

const router = useRouter()
const route = useRoute()
const roadbookStore = useRoadbookStore()

// 是否编辑模式
const isEdit = computed(() => !!route.params.id)
const roadbookId = computed(() => Number(route.params.id) || 0)

// 表单相关
const formRef = ref<FormInstance>()
const saving = ref(false)

// 表单数据
const formData = ref<CreateRoadbookRequest & { status?: RoadbookStatus }>({
  title: '',
  description: '',
  coverUrl: '',
  startDate: '',
  endDate: '',
  visibility: RoadbookVisibility.PUBLIC,
  travelMode: 'driving' as TravelMode,
  totalBudget: 0,
  tagIds: []
})

// 日期范围
const dateRange = ref<[string, string] | null>(null)

// 表单验证规则
const formRules: FormRules = {
  title: [
    { required: true, message: '请输入路书标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
  ]
}

// 途经点列表
const waypoints = ref<Partial<Waypoint>[]>([])

// 计算天数
const dayCount = computed(() => {
  if (!formData.value.startDate || !formData.value.endDate) return 0
  const start = new Date(formData.value.startDate)
  const end = new Date(formData.value.endDate)
  const diff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
  return diff + 1
})

// 监听日期范围变化
watch(dateRange, (val) => {
  if (val && val.length === 2) {
    formData.value.startDate = val[0]
    formData.value.endDate = val[1]
  } else {
    formData.value.startDate = ''
    formData.value.endDate = ''
  }
})

// 禁用过去的日期
const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

// 地图相关
const mapContainerRef = ref<InstanceType<typeof AMapContainer> | null>(null)
const mapCenter = ref<[number, number]>([116.397428, 39.90923])
const routePolyline = shallowRef<any>(null)
const { planRouteFromWaypoints, formattedDistance, formattedDuration, loading: routeLoading } = useRouteHook()

// 路线信息
const routeInfo = computed(() => {
  if (!formattedDistance.value || !formattedDuration.value) return null
  return {
    distance: formattedDistance.value,
    duration: formattedDuration.value
  }
})

// 途经点编辑弹窗
const waypointDialogVisible = ref(false)
const editingWaypoint = ref<Partial<Waypoint> | null>(null)
const editingWaypointIndex = ref(-1)

// 地图就绪
const handleMapReady = () => {
  if (waypoints.value.length > 0) {
    updateMapMarkers()
    handleRefreshRoute()
  }
}

// 地图标记列表
const mapMarkers = shallowRef<any[]>([])

// 更新地图标记
const updateMapMarkers = () => {
  const map = mapContainerRef.value?.getMap()
  if (!map) return

  const AMap = getAMap()
  if (!AMap) return

  // 清除旧标记
  mapMarkers.value.forEach((marker) => map.remove(marker))
  mapMarkers.value = []

  // 添加新标记
  const validWaypoints = waypoints.value.filter(
    (wp): wp is Waypoint => !!(wp.longitude && wp.latitude)
  )

  validWaypoints.forEach((wp, index) => {
    const marker = new AMap.Marker({
      position: [wp.longitude, wp.latitude],
      title: wp.name
    })
    marker.setLabel({
      content: `<div class="amap-marker-label">${index + 1}</div>`,
      offset: new AMap.Pixel(-10, -10)
    })
    map.add(marker)
    mapMarkers.value.push(marker)
  })
}

// 清除路线
const clearRoutePolyline = () => {
  const map = mapContainerRef.value?.getMap()
  if (map && routePolyline.value) {
    map.remove(routePolyline.value)
    routePolyline.value = null
  }
}

// 刷新路线
const handleRefreshRoute = async () => {
  const validWaypoints = waypoints.value.filter(
    (wp): wp is Waypoint => !!(wp.longitude && wp.latitude)
  )

  if (validWaypoints.length < 2) {
    clearRoutePolyline()
    return
  }

  try {
    const result = await planRouteFromWaypoints(validWaypoints, formData.value.travelMode)
    const map = mapContainerRef.value?.getMap()
    const AMap = getAMap()

    if (result && map && AMap) {
      // 清除旧路线
      clearRoutePolyline()

      // 绘制新路线
      routePolyline.value = new AMap.Polyline({
        path: result.polyline,
        strokeColor: '#3366FF',
        strokeWeight: 6,
        strokeOpacity: 0.8,
        lineJoin: 'round',
        lineCap: 'round'
      })
      map.add(routePolyline.value)
      map.setFitView()
    }
  } catch (err) {
    console.warn('Route planning failed:', err)
    // 即使路线规划失败，也尝试自适应显示所有标记
    const map = mapContainerRef.value?.getMap()
    if (map) {
      map.setFitView()
    }
  }
}

// 添加途经点
const handleAddWaypoint = () => {
  editingWaypoint.value = {
    name: '',
    address: '',
    longitude: 0,
    latitude: 0,
    dayIndex: 1,
    sortOrder: waypoints.value.length,
    stayDuration: 60,
    budget: 0,
    notes: '',
    images: [],
    waypointType: WaypointType.WAYPOINT
  }
  editingWaypointIndex.value = -1
  waypointDialogVisible.value = true
}

// 编辑途经点
const handleEditWaypoint = (waypoint: Partial<Waypoint>, index: number) => {
  editingWaypoint.value = { ...waypoint }
  editingWaypointIndex.value = index
  waypointDialogVisible.value = true
}

// 删除途经点
const handleDeleteWaypoint = (index: number) => {
  waypoints.value.splice(index, 1)
  // 重新排序
  waypoints.value.forEach((wp, i) => {
    wp.sortOrder = i
  })
  updateMapMarkers()
  handleRefreshRoute()
}

// 保存途经点
const handleSaveWaypoint = (waypoint: Partial<Waypoint>) => {
  if (editingWaypointIndex.value >= 0) {
    waypoints.value[editingWaypointIndex.value] = waypoint
  } else {
    waypoints.value.push(waypoint)
  }
  waypointDialogVisible.value = false
  updateMapMarkers()
  handleRefreshRoute()
}

// 途经点重新排序
const handleReorderWaypoints = (newWaypoints: Partial<Waypoint>[]) => {
  waypoints.value = newWaypoints.map((wp, index) => ({
    ...wp,
    sortOrder: index
  }))
  handleRefreshRoute()
}

// 封面上传前验证
const beforeCoverUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 封面上传处理
const handleCoverUpload = async (options: UploadRequestOptions) => {
  // TODO: 实现实际的图片上传逻辑
  // 这里暂时使用本地预览
  const file = options.file
  const reader = new FileReader()
  reader.onload = (e) => {
    formData.value.coverUrl = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

// 返回
const handleBack = () => {
  router.back()
}

// 保存草稿
const handleSaveDraft = async () => {
  saving.value = true
  try {
    formData.value.status = RoadbookStatus.DRAFT
    await saveRoadbook()
    ElMessage.success('草稿保存成功')
  } catch (err: any) {
    ElMessage.error(err.message || '保存失败')
  } finally {
    saving.value = false
  }
}

// 发布
const handlePublish = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    formData.value.status = RoadbookStatus.PUBLISHED
    await saveRoadbook()
    ElMessage.success(isEdit.value ? '更新成功' : '发布成功')
    router.push({ name: 'RoadbookList' })
  } catch (err: any) {
    ElMessage.error(err.message || '发布失败')
  } finally {
    saving.value = false
  }
}

// 格式化日期为 ISO 8601 格式
const formatDateToISO = (dateStr: string): string => {
  if (!dateStr) return ''
  // 如果已经是 ISO 格式，直接返回
  if (dateStr.includes('T')) return dateStr
  // 将 YYYY-MM-DD 转换为 ISO 格式
  return `${dateStr}T00:00:00Z`
}

// 保存路书
const saveRoadbook = async () => {
  // 准备请求数据，转换日期格式
  const requestData = {
    ...formData.value,
    startDate: formatDateToISO(formData.value.startDate || ''),
    endDate: formatDateToISO(formData.value.endDate || '')
  }

  if (isEdit.value) {
    await roadbookStore.updateRoadbook(roadbookId.value, requestData)
    // TODO: 实现途经点批量更新
  } else {
    const roadbook = await roadbookStore.createRoadbook(requestData)
    // 添加途经点
    for (const wp of waypoints.value) {
      if (wp.name && wp.longitude && wp.latitude) {
        await roadbookStore.addWaypoint(roadbook.id, {
          name: wp.name,
          address: wp.address || '',
          longitude: wp.longitude,
          latitude: wp.latitude,
          poiId: wp.poiId || '',
          poiType: wp.poiType || '',
          dayIndex: wp.dayIndex || 1,
          sortOrder: wp.sortOrder || 0,
          stayDuration: wp.stayDuration || 60,
          budget: wp.budget || 0,
          notes: wp.notes || '',
          images: wp.images || [],
          waypointType: wp.waypointType || WaypointType.WAYPOINT
        })
      }
    }
  }
}

// 加载路书数据（编辑模式）
const loadRoadbook = async () => {
  if (!isEdit.value) return

  try {
    await roadbookStore.fetchRoadbookDetail(roadbookId.value)
    const roadbook = roadbookStore.currentRoadbook
    if (roadbook) {
      formData.value = {
        title: roadbook.title,
        description: roadbook.description,
        coverUrl: roadbook.coverUrl,
        startDate: roadbook.startDate,
        endDate: roadbook.endDate,
        visibility: roadbook.visibility,
        travelMode: roadbook.travelMode,
        totalBudget: roadbook.totalBudget,
        tagIds: roadbook.tags?.map(t => t.id) || []
      }
      
      if (roadbook.startDate && roadbook.endDate) {
        dateRange.value = [roadbook.startDate, roadbook.endDate]
      }
      
      waypoints.value = roadbook.waypoints || []
      
      // 更新地图
      if (waypoints.value.length > 0) {
        const firstWp = waypoints.value[0]
        if (firstWp.longitude && firstWp.latitude) {
          mapCenter.value = [firstWp.longitude, firstWp.latitude]
        }
      }
    }
  } catch (err: any) {
    ElMessage.error('加载路书失败')
    router.back()
  }
}

onMounted(() => {
  loadRoadbook()
})
</script>


<style scoped lang="scss">
.roadbook-editor {
  .editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--art-border-color);

    .header-left {
      display: flex;
      align-items: center;
      gap: 16px;

      .page-title {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
        color: var(--art-text-gray-800);
      }
    }

    .header-right {
      display: flex;
      gap: 12px;
    }
  }

  .editor-body {
    .editor-left,
    .editor-right {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }

    .el-card {
      border-radius: calc(var(--custom-radius) / 2 + 2px);

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-weight: 600;

        .header-tip {
          font-size: 12px;
          font-weight: 400;
          color: var(--art-text-gray-400);
        }
      }
    }

    // 封面上传
    .cover-upload {
      width: 100%;

      .cover-uploader {
        width: 100%;

        :deep(.el-upload) {
          width: 100%;
        }
      }

      .cover-preview {
        position: relative;
        width: 100%;
        aspect-ratio: 16/10;
        border-radius: 8px;
        overflow: hidden;
        cursor: pointer;

        .el-image {
          width: 100%;
          height: 100%;
        }

        .cover-mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          gap: 8px;
          background: rgba(0, 0, 0, 0.5);
          color: #fff;
          opacity: 0;
          transition: opacity 0.3s;

          .el-icon {
            font-size: 32px;
          }
        }

        &:hover .cover-mask {
          opacity: 1;
        }
      }

      .cover-placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 8px;
        width: 100%;
        aspect-ratio: 16/10;
        border: 2px dashed var(--art-border-color);
        border-radius: 8px;
        background: var(--art-gray-100);
        color: var(--art-text-gray-500);
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: var(--el-color-primary);
          color: var(--el-color-primary);
        }

        .tip {
          font-size: 12px;
          color: var(--art-text-gray-400);
        }
      }
    }

    // 地图容器
    .map-card {
      .map-container {
        height: 400px;
        border-radius: 8px;
        overflow: hidden;
      }

      .route-info {
        display: flex;
        gap: 24px;
        margin-top: 16px;
        padding: 12px 16px;
        background: var(--art-gray-100);
        border-radius: 8px;

        .info-item {
          display: flex;
          flex-direction: column;
          gap: 4px;

          .label {
            font-size: 12px;
            color: var(--art-text-gray-500);
          }

          .value {
            font-size: 16px;
            font-weight: 600;
            color: var(--el-color-primary);
          }
        }
      }
    }

  }
}

// 响应式布局
@media only screen and (max-width: 1200px) {
  .roadbook-editor {
    .editor-body {
      .el-col {
        margin-bottom: 20px;
      }
    }
  }
}

@media only screen and (max-width: 768px) {
  .roadbook-editor {
    .editor-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;

      .header-right {
        width: 100%;
        justify-content: flex-end;
      }
    }
  }
}
</style>
