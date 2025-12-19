<template>
  <div class="version-detail">
    <!-- 基本信息 / Basic Info -->
    <ElDescriptions :column="2" border>
      <ElDescriptionsItem label="版本号">V{{ version.versionNumber }}</ElDescriptionsItem>
      <ElDescriptionsItem label="版本名称">{{ version.versionName }}</ElDescriptionsItem>
      <ElDescriptionsItem label="状态">
        <ElTag :type="getStatusType(version.status)">{{ getStatusText(version.status) }}</ElTag>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="创建时间">{{ formatDate(version.createdAt) }}</ElDescriptionsItem>
      <ElDescriptionsItem label="版本说明" :span="2">
        {{ version.description || '暂无说明' }}
      </ElDescriptionsItem>
    </ElDescriptions>

    <!-- 设计图 / Design Images -->
    <div class="section" v-if="version.designImages?.length">
      <h4>设计图</h4>
      <div class="image-grid">
        <div v-for="(img, index) in version.designImages" :key="index" class="image-item">
          <ElImage
            :src="img.url"
            :preview-src-list="version.designImages.map(i => i.url)"
            :initial-index="index"
            fit="cover"
          />
          <div class="image-info">
            <span class="name">{{ img.name }}</span>
            <span class="type">{{ getImageTypeText(img.type) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 布局信息 / Layout Info -->
    <div class="section" v-if="version.layoutInfo?.length">
      <h4>布局信息</h4>
      <ElTable :data="version.layoutInfo" stripe size="small">
        <ElTableColumn prop="zone" label="区域" />
        <ElTableColumn prop="area" label="面积">
          <template #default="{ row }">{{ row.area }} m²</template>
        </ElTableColumn>
        <ElTableColumn prop="position" label="位置" />
      </ElTable>
    </div>

    <!-- 面积信息 / Area Info -->
    <div class="section" v-if="version.areaInfo?.length">
      <h4>面积信息</h4>
      <ElTable :data="version.areaInfo" stripe size="small">
        <ElTableColumn prop="name" label="名称" />
        <ElTableColumn prop="value" label="面积">
          <template #default="{ row }">{{ row.value }} {{ row.unit }}</template>
        </ElTableColumn>
      </ElTable>
    </div>

    <!-- 风格信息 / Style Info -->
    <div class="section" v-if="version.styleInfo?.length">
      <h4>风格信息</h4>
      <ElTable :data="version.styleInfo" stripe size="small">
        <ElTableColumn prop="category" label="类别" />
        <ElTableColumn prop="value" label="风格" />
        <ElTableColumn prop="details" label="详情" />
      </ElTable>
    </div>

    <!-- 材料信息 / Material Info -->
    <div class="section" v-if="version.materialInfo?.length">
      <h4>材料信息</h4>
      <ElTable :data="version.materialInfo" stripe size="small">
        <ElTableColumn prop="name" label="材料名称" />
        <ElTableColumn prop="category" label="类别" />
        <ElTableColumn prop="quantity" label="数量">
          <template #default="{ row }">{{ row.quantity }} {{ row.unit }}</template>
        </ElTableColumn>
      </ElTable>
    </div>
  </div>
</template>

<script setup lang="ts">
/***
 * Version Detail Component
 * 版本详情组件
 * Requirements: 7.2
 ***/
import dayjs from 'dayjs'
import type { DesignVersionResponse } from '@/api/designer-version'

defineProps<{
  version: DesignVersionResponse
}>()

/*** Format date ***/
const formatDate = (date: string) => {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
}

/*** Get status type ***/
type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
const getStatusType = (status: string): TagType => {
  const map: Record<string, TagType> = {
    draft: 'info',
    submitted: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

/*** Get status text ***/
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    draft: '草稿',
    submitted: '已提交',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return map[status] || status
}

/*** Get image type text ***/
const getImageTypeText = (type: string) => {
  const map: Record<string, string> = {
    floor_plan: '平面图',
    elevation: '立面图',
    '3d_render': '3D效果图',
    section: '剖面图',
    detail: '详图'
  }
  return map[type] || type
}
</script>

<style scoped lang="scss">
.version-detail {
  .section {
    margin-top: 24px;

    h4 {
      margin-bottom: 12px;
      font-size: 14px;
      font-weight: 500;
      color: var(--el-text-color-primary);
    }
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;

    .image-item {
      border: 1px solid var(--el-border-color);
      border-radius: 8px;
      overflow: hidden;

      :deep(.el-image) {
        width: 100%;
        height: 150px;
      }

      .image-info {
        padding: 8px;
        background: var(--el-fill-color-light);

        .name {
          display: block;
          font-size: 13px;
          color: var(--el-text-color-primary);
        }

        .type {
          display: block;
          font-size: 12px;
          color: var(--el-text-color-secondary);
          margin-top: 4px;
        }
      }
    }
  }
}
</style>
