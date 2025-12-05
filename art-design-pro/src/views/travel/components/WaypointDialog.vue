<template>
  <ElDialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑途经点' : '添加途经点'"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <ElForm
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
      label-position="top"
    >
      <!-- 地图选点 -->
      <ElFormItem label="选择位置" required>
        <div class="map-picker">
          <div class="map-tip">
            <ElIcon>
              <InfoFilled />
            </ElIcon>
            <span>点击地图选择位置</span>
          </div>
          <div ref="mapContainerRef" class="map-container"></div>
          <div class="coord-display">
            <span class="coord-label">经度:</span>
            <span class="coord-value">{{ formData.longitude?.toFixed(6) || '-' }}</span>
            <span class="coord-label">纬度:</span>
            <span class="coord-value">{{ formData.latitude?.toFixed(6) || '-' }}</span>
          </div>
        </div>
      </ElFormItem>

      <!-- 地点名称 -->
      <ElFormItem label="地点名称" prop="name">
        <ElInput v-model="formData.name" placeholder="请输入地点名称" maxlength="100" />
      </ElFormItem>

      <!-- 地址 -->
      <ElFormItem label="详细地址" prop="address">
        <ElInput v-model="formData.address" placeholder="请输入详细地址" maxlength="300" />
      </ElFormItem>

      <!-- 所属天数 -->
      <ElFormItem label="所属天数(天)" prop="dayIndex">
        <ElInput v-model="formData.dayIndex" placeholder="请输入所属天数" maxlength="300" />
      </ElFormItem>

      <!-- 停留时间 -->
      <ElFormItem label="停留时间" prop="stayDuration">
        <ElSelect v-model="formData.stayDuration" style="width: 100%">
          <ElOption label="30分钟" :value="30" />
          <ElOption label="1小时" :value="60" />
          <ElOption label="1.5小时" :value="90" />
          <ElOption label="2小时" :value="120" />
          <ElOption label="3小时" :value="180" />
          <ElOption label="半天" :value="240" />
          <ElOption label="一天" :value="480" />
        </ElSelect>
      </ElFormItem>

      <!-- 预算 -->
      <ElFormItem label="预算金额" prop="budget">
        <ElInputNumber
          v-model="formData.budget"
          :min="0"
          :precision="2"
          placeholder="预算金额"
          style="width: 100%"
        >
          <template #suffix>元</template>
        </ElInputNumber>
      </ElFormItem>

      <!-- 备注 -->
      <ElFormItem label="备注" prop="notes">
        <ElInput
          v-model="formData.notes"
          type="textarea"
          placeholder="添加备注信息"
          :rows="3"
          maxlength="500"
          show-word-limit
        />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton type="primary" @click="handleSave" :loading="saving">
        {{ isEdit ? '保存修改' : '添加' }}
      </ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ref, computed, watch, shallowRef, nextTick, onUnmounted } from 'vue'
  import type { FormInstance, FormRules } from 'element-plus'
  import { InfoFilled } from '@element-plus/icons-vue'
  import type { Waypoint } from '../types'
  import { WaypointType } from '../types'
  import { loadAMapScript } from '../composables/useAMap'

  defineOptions({
    name: 'WaypointDialog'
  })

  // Props
  interface Props {
    visible: boolean
    waypoint?: Partial<Waypoint> | null
    days?: number
  }

  const props = withDefaults(defineProps<Props>(), {
    visible: false,
    waypoint: null,
    days: 1
  })

  // Emits
  const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'save', waypoint: Partial<Waypoint>): void
  }>()

  // 状态
  const formRef = ref<FormInstance>()
  const saving = ref(false)

  // 地图相关
  const mapContainerRef = ref<HTMLElement | null>(null)
  const mapInstance = shallowRef<any>(null)
  const markerInstance = shallowRef<any>(null)
  const geocoder = shallowRef<any>(null)

  // 对话框可见性
  const dialogVisible = computed({
    get: () => props.visible,
    set: (val) => emit('update:visible', val)
  })

  // 是否编辑模式
  const isEdit = computed(() => !!props.waypoint?.id || !!props.waypoint?.name)

  // 实际天数（至少为1）
  // const actualDays = computed(() => Math.max(props.days, 1))

  // 表单数据
  const formData = ref<Partial<Waypoint>>({
    name: '',
    address: '',
    longitude: 0,
    latitude: 0,
    dayIndex: 1,
    sortOrder: 0,
    stayDuration: 60,
    budget: 0,
    notes: '',
    images: [],
    waypointType: WaypointType.WAYPOINT
  })

  // 表单验证规则
  const formRules: FormRules = {
    name: [{ required: true, message: '请输入地点名称', trigger: 'blur' }]
  }

  // 初始化地图
  const initMap = async () => {
    if (!mapContainerRef.value || mapInstance.value) return

    try {
      const AMap = await loadAMapScript()

      // 默认中心点（北京）或已有坐标
      const center: [number, number] =
        formData.value.longitude && formData.value.latitude
          ? [formData.value.longitude, formData.value.latitude]
          : [116.397428, 39.90923]

      mapInstance.value = new AMap.Map(mapContainerRef.value, {
        center,
        zoom: 14,
        resizeEnable: true
      })

      // 初始化地理编码服务
      geocoder.value = new AMap.Geocoder({
        city: '全国',
        radius: 1000
      })

      // 如果已有坐标，添加标记
      if (formData.value.longitude && formData.value.latitude) {
        addMarker([formData.value.longitude, formData.value.latitude])
      }

      // 地图点击事件
      mapInstance.value.on('click', handleMapClick)
    } catch (err) {
      console.error('Failed to init map:', err)
    }
  }

  // 添加/更新标记
  const addMarker = (position: [number, number]) => {
    if (!mapInstance.value) return

    const AMap = (window as any).AMap
    if (!AMap) return

    // 移除旧标记
    if (markerInstance.value) {
      mapInstance.value.remove(markerInstance.value)
    }

    // 创建新标记
    markerInstance.value = new AMap.Marker({
      position,
      draggable: true
    })

    mapInstance.value.add(markerInstance.value)

    // 标记拖拽结束事件
    markerInstance.value.on('dragend', (e: any) => {
      const pos = e.target.getPosition()
      updateLocationFromCoords([pos.getLng(), pos.getLat()])
    })
  }

  // 地图点击处理
  const handleMapClick = (e: any) => {
    const lng = e.lnglat.getLng()
    const lat = e.lnglat.getLat()
    console.log(e)

    updateLocationFromCoords([lng, lat])
  }

  // 根据坐标更新位置信息
  const updateLocationFromCoords = (position: [number, number]) => {
    const [lng, lat] = position

    // 更新坐标
    formData.value.longitude = lng
    formData.value.latitude = lat

    // 添加/更新标记
    addMarker(position)

    // 逆地理编码获取地址
    // if (geocoder.value) {
    //   geocoder.value.getAddress(position, (status: string, result: any) => {
    //     if (status === 'complete' && result.regeocode) {
    //       const address = result.regeocode.formattedAddress || ''
    //       const addressComponent = result.regeocode.addressComponent
    //       const poi = result.regeocode.pois?.[0]

    //       // 优先使用 POI 名称，其次是建筑物名称，最后是街道+门牌号
    //       const name =
    //         poi?.name ||
    //         addressComponent?.building ||
    //         (addressComponent?.street
    //           ? `${addressComponent.street}${addressComponent.streetNumber || ''}`
    //           : '未命名地点')

    //       // 自动填充地点名称和地址
    //       formData.value.name = name
    //       formData.value.address = address
    //     } else {
    //       // 逆地理编码失败，使用坐标
    //       formData.value.name = '自定义地点'
    //       formData.value.address = `经度: ${lng.toFixed(6)}, 纬度: ${lat.toFixed(6)}`
    //     }
    //   })
    // }
  }

  // 销毁地图
  const destroyMap = () => {
    if (mapInstance.value) {
      mapInstance.value.off('click', handleMapClick)
      mapInstance.value.destroy()
      mapInstance.value = null
      markerInstance.value = null
    }
  }

  // 重置表单
  const resetForm = () => {
    formData.value = {
      name: '',
      address: '',
      longitude: 0,
      latitude: 0,
      dayIndex: 1,
      sortOrder: 0,
      stayDuration: 60,
      budget: 0,
      notes: '',
      images: [],
      waypointType: WaypointType.WAYPOINT
    }
    formRef.value?.resetFields()
  }

  // 监听 waypoint 变化
  watch(
    () => props.waypoint,
    (val) => {
      if (val) {
        formData.value = { ...val }
      } else {
        resetForm()
      }
    },
    { immediate: true }
  )

  // 监听弹窗打开，初始化地图
  watch(
    () => props.visible,
    async (visible) => {
      if (visible) {
        await nextTick()
        initMap()
      } else {
        destroyMap()
      }
    }
  )

  // 组件卸载时销毁地图
  onUnmounted(() => {
    destroyMap()
  })

  // 关闭对话框
  const handleClose = () => {
    dialogVisible.value = false
    resetForm()
  }

  // 保存
  const handleSave = async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return

    saving.value = true
    try {
      emit('save', { ...formData.value })
      handleClose()
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped lang="scss">
  :deep(.el-dialog__body) {
    padding-top: 10px;
  }

  .map-picker {
    width: 100%;

    .map-tip {
      display: flex;
      align-items: center;
      gap: 6px;
      margin-bottom: 8px;
      padding: 6px 10px;
      background: var(--el-color-primary-light-9);
      border-radius: 4px;
      font-size: 12px;
      color: var(--el-color-primary);

      .el-icon {
        font-size: 14px;
      }
    }

    .map-container {
      width: 100%;
      height: 250px;
      border-radius: 8px;
      overflow: hidden;
      border: 1px solid var(--el-border-color);
    }

    .coord-display {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 8px;
      padding: 8px 12px;
      background: var(--el-fill-color-light);
      border-radius: 4px;
      font-size: 13px;

      .coord-label {
        color: var(--el-text-color-secondary);
      }

      .coord-value {
        font-family: monospace;
        color: var(--el-text-color-primary);
      }
    }
  }
</style>
