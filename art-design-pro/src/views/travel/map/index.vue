<template>
  <div class="travel-map-page">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索地点"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button @click="handleSearch">搜索</el-button>
        </template>
      </el-input>
    </div>

    <!-- 地图容器 -->
    <AMapContainer
      ref="mapContainerRef"
      :center="mapStore.center"
      :zoom="mapStore.zoom"
      :show-zoom-control="true"
      :show-scale="true"
      @ready="handleMapReady"
      @click="handleMapClick"
    >
      <template #controls>
        <!-- 出行方式选择 -->
        <div class="travel-mode-selector">
          <el-radio-group v-model="currentTravelMode" size="small" @change="handleTravelModeChange">
            <el-radio-button value="driving">驾车</el-radio-button>
            <el-radio-button value="walking">步行</el-radio-button>
            <el-radio-button value="cycling">骑行</el-radio-button>
            <el-radio-button value="transit">公交</el-radio-button>
          </el-radio-group>
        </div>
      </template>
    </AMapContainer>

    <!-- 搜索结果面板 -->
    <div v-if="searchResults.length > 0" class="search-results-panel">
      <div class="panel-header">
        <span>搜索结果</span>
        <el-button text @click="clearSearch">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
      <div class="results-list">
        <div
          v-for="poi in searchResults"
          :key="poi.id"
          class="result-item"
          @click="handleSelectPOI(poi)"
        >
          <div class="poi-name">{{ poi.name }}</div>
          <div class="poi-address">{{ poi.address }}</div>
        </div>
      </div>
    </div>

    <!-- 路线信息面板 -->
    <div v-if="routeInfo.hasRoute" class="route-info-panel">
      <div class="panel-header">
        <span>路线信息</span>
        <el-button text @click="clearRoute">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
      <div class="route-summary">
        <div class="info-item">
          <span class="label">总距离</span>
          <span class="value">{{ routeInfo.formattedDistance }}</span>
        </div>
        <div class="info-item">
          <span class="label">预计时间</span>
          <span class="value">{{ routeInfo.formattedDuration }}</span>
        </div>
      </div>
    </div>

    <!-- 点击位置信息 -->
    <el-dialog
      v-model="showLocationDialog"
      title="位置信息"
      width="400px"
      :close-on-click-modal="true"
    >
      <div v-if="clickedLocation" class="location-info">
        <p><strong>经度:</strong> {{ clickedLocation.lnglat[0].toFixed(6) }}</p>
        <p><strong>纬度:</strong> {{ clickedLocation.lnglat[1].toFixed(6) }}</p>
      </div>
      <template #footer>
        <el-button @click="showLocationDialog = false">关闭</el-button>
        <el-button type="primary" @click="addAsWaypoint">添加为途经点</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, Close } from '@element-plus/icons-vue'
import AMapContainer from '../components/AMapContainer.vue'
import { useAMap, useRoute } from '../composables'
import { useMapStore } from '../store/map'
import type { POISearchResult, TravelMode } from '../types'

defineOptions({
  name: 'TravelMap'
})

// Store
const mapStore = useMapStore()

// Refs
const mapContainerRef = ref<InstanceType<typeof AMapContainer> | null>(null)
const searchKeyword = ref('')
const currentTravelMode = ref<TravelMode>('driving' as TravelMode)
const showLocationDialog = ref(false)
const clickedLocation = ref<{ lnglat: [number, number]; pixel: [number, number] } | null>(null)

// Hooks
const {
  searchPOI,
  searchResults,
  searchLoading,
  setMarkersFromWaypoints,
  clearMarkers
} = useAMap()

const routeInfo = useRoute()

// 计算属性
const hasSearchResults = computed(() => searchResults.value.length > 0)

// 地图就绪
function handleMapReady(map: any) {
  mapStore.setMapInstance(map)
}

// 地图点击
function handleMapClick(event: { lnglat: [number, number]; pixel: [number, number] }) {
  clickedLocation.value = event
  showLocationDialog.value = true
}

// 搜索
async function handleSearch() {
  if (!searchKeyword.value.trim()) return

  try {
    await searchPOI(searchKeyword.value)
  } catch (err) {
    console.error('搜索失败:', err)
  }
}

// 清除搜索
function clearSearch() {
  searchKeyword.value = ''
  searchResults.value = []
}

// 选择POI
function handleSelectPOI(poi: POISearchResult) {
  const map = mapContainerRef.value?.getMap()
  if (map) {
    map.setCenter([poi.location.lng, poi.location.lat])
    map.setZoom(15)
  }
  mapStore.setCenter(poi.location.lng, poi.location.lat)
}

// 切换出行方式
function handleTravelModeChange(mode: string | number | boolean | undefined) {
  if (mode) {
    routeInfo.setTravelMode(mode as TravelMode)
    mapStore.setTravelMode(mode as TravelMode)
  }
}

// 清除路线
function clearRoute() {
  const map = mapContainerRef.value?.getMap()
  if (map) {
    routeInfo.clearRouteFromMap(map)
  }
  routeInfo.clearRoute()
}

// 添加为途经点
function addAsWaypoint() {
  if (clickedLocation.value) {
    // 这里可以触发添加途经点的逻辑
    console.log('添加途经点:', clickedLocation.value.lnglat)
  }
  showLocationDialog.value = false
}

onMounted(() => {
  // 初始化
})
</script>

<style scoped lang="scss">
.travel-map-page {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 600px;

  .search-bar {
    position: absolute;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 100;
    width: 400px;

    :deep(.el-input-group__append) {
      background-color: var(--el-color-primary);
      color: #fff;
      border-color: var(--el-color-primary);
    }
  }

  .travel-mode-selector {
    background: #fff;
    padding: 8px;
    border-radius: 4px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  }

  .search-results-panel,
  .route-info-panel {
    position: absolute;
    top: 80px;
    right: 20px;
    width: 320px;
    max-height: 400px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    z-index: 100;
    overflow: hidden;

    .panel-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      font-weight: 500;
    }

    .results-list {
      max-height: 320px;
      overflow-y: auto;

      .result-item {
        padding: 12px 16px;
        cursor: pointer;
        transition: background-color 0.2s;

        &:hover {
          background-color: var(--el-fill-color-light);
        }

        .poi-name {
          font-size: 14px;
          color: var(--el-text-color-primary);
          margin-bottom: 4px;
        }

        .poi-address {
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
      }
    }

    .route-summary {
      padding: 16px;

      .info-item {
        display: flex;
        justify-content: space-between;
        margin-bottom: 12px;

        &:last-child {
          margin-bottom: 0;
        }

        .label {
          color: var(--el-text-color-secondary);
        }

        .value {
          font-weight: 500;
          color: var(--el-color-primary);
        }
      }
    }
  }

  .location-info {
    p {
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}
</style>
