<template>
  <div class="xhs-favorites-page">
    <ElCard shadow="never" class="toolbar-card">
      <div class="toolbar">
        <div class="folder-tabs">
          <ElButton :type="currentFolderId === 0 ? 'primary' : 'default'" @click="selectFolder(0)">
            全部收藏
          </ElButton>
          <ElButton
            v-for="folder in folders"
            :key="folder.id"
            :type="currentFolderId === folder.id ? 'primary' : 'default'"
            @click="selectFolder(folder.id)"
          >
            {{ folder.name }}
          </ElButton>
          <ElButton :icon="Plus" @click="showFolderDialog()">新建文件夹</ElButton>
        </div>
      </div>
    </ElCard>

    <div v-loading="loading" class="favorites-list">
      <ElRow :gutter="16">
        <ElCol
          v-for="item in favoritesList"
          :key="item.id"
          :xs="24"
          :sm="12"
          :lg="8"
          :xl="6"
          class="mb-4"
        >
          <ElCard shadow="hover" class="favorite-card">
            <div class="card-header">
              <span class="card-date">{{ formatDate(item.created_at) }}</span>
              <ElDropdown @command="handleCommand($event, item)">
                <ElButton link size="small">
                  <ElIcon><MoreFilled /></ElIcon>
                </ElButton>
                <template #dropdown>
                  <ElDropdownMenu>
                    <ElDropdownItem command="move">移动到文件夹</ElDropdownItem>
                    <ElDropdownItem command="unfavorite">取消收藏</ElDropdownItem>
                  </ElDropdownMenu>
                </template>
              </ElDropdown>
            </div>
            <h4 class="card-title">{{ item.summary_title || item.note_title || '未命名' }}</h4>
            <p class="card-content">{{ getPreview(item.summary_content) }}</p>
            <div class="card-tags" v-if="item.tags?.length">
              <ElTag
                v-for="tag in item.tags.slice(0, 3)"
                :key="tag"
                size="small"
                effect="plain"
                class="mr-1"
              >
                {{ tag }}
              </ElTag>
            </div>
          </ElCard>
        </ElCol>
      </ElRow>

      <ElEmpty v-if="!loading && favoritesList.length === 0" description="暂无收藏内容">
        <ElButton type="primary" @click="$router.push('/xhs/history')">去历史记录查看</ElButton>
      </ElEmpty>

      <div v-if="total > pageSize" class="pagination-wrapper">
        <ElPagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadFavorites"
        />
      </div>
    </div>

    <!-- 新建/编辑文件夹弹窗 -->
    <ElDialog
      v-model="folderDialogVisible"
      :title="editingFolder ? '编辑文件夹' : '新建文件夹'"
      width="400px"
    >
      <ElForm :model="folderForm" label-width="80px">
        <ElFormItem label="名称">
          <ElInput v-model="folderForm.name" placeholder="请输入文件夹名称" maxlength="20" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="folderDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveFolder" :loading="saving">保存</ElButton>
      </template>
    </ElDialog>

    <!-- 移动到文件夹弹窗 -->
    <ElDialog v-model="moveDialogVisible" title="移动到文件夹" width="400px">
      <ElRadioGroup v-model="targetFolderId" class="folder-radio-group">
        <ElRadio :label="0">无文件夹</ElRadio>
        <ElRadio v-for="folder in folders" :key="folder.id" :label="folder.id">
          {{ folder.name }}
        </ElRadio>
      </ElRadioGroup>
      <template #footer>
        <ElButton @click="moveDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="moveToFolder">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { Plus, MoreFilled } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import { xhsApi, type XHSSummary, type XHSFolder } from '@/api/xhs'
  import dayjs from 'dayjs'

  defineOptions({ name: 'XHSFavorites' })

  const currentFolderId = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(12)
  const total = ref(0)
  const loading = ref(false)
  const favoritesList = ref<XHSSummary[]>([])
  const folders = ref<XHSFolder[]>([])

  const folderDialogVisible = ref(false)
  const editingFolder = ref<XHSFolder | null>(null)
  const folderForm = reactive({ name: '' })
  const saving = ref(false)

  const moveDialogVisible = ref(false)
  const targetFolderId = ref(0)
  const movingItem = ref<XHSSummary | null>(null)

  const formatDate = (date: string) => dayjs(date).format('MM-DD HH:mm')
  const getPreview = (content: string) =>
    content?.length > 80 ? content.slice(0, 80) + '...' : content

  const selectFolder = (id: number) => {
    currentFolderId.value = id
    currentPage.value = 1
    loadFavorites()
  }

  const loadFavorites = async () => {
    loading.value = true
    try {
      const res = await xhsApi.getFavorites({
        folder_id: currentFolderId.value || undefined,
        page: currentPage.value,
        page_size: pageSize.value
      })
      const data = (res as any)?.data || res
      favoritesList.value = data?.list || []
      total.value = data?.total || 0
    } catch (error) {
      console.error('加载收藏失败:', error)
    } finally {
      loading.value = false
    }
  }

  const loadFolders = async () => {
    try {
      const res = await xhsApi.getFolders()
      folders.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载文件夹失败:', error)
    }
  }

  const showFolderDialog = (folder?: XHSFolder) => {
    editingFolder.value = folder || null
    folderForm.name = folder?.name || ''
    folderDialogVisible.value = true
  }

  const saveFolder = async () => {
    if (!folderForm.name.trim()) {
      ElMessage.warning('请输入文件夹名称')
      return
    }
    saving.value = true
    try {
      if (editingFolder.value) {
        await xhsApi.updateFolder(editingFolder.value.id, { name: folderForm.name })
      } else {
        await xhsApi.createFolder(folderForm.name)
      }
      ElMessage.success('保存成功')
      folderDialogVisible.value = false
      loadFolders()
    } catch (error: any) {
      ElMessage.error(error.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  const handleCommand = (command: string, item: XHSSummary) => {
    if (command === 'move') {
      movingItem.value = item
      targetFolderId.value = item.folder_id
      moveDialogVisible.value = true
    } else if (command === 'unfavorite') {
      unfavorite(item)
    }
  }

  const unfavorite = async (item: XHSSummary) => {
    try {
      await xhsApi.updateFavorite(item.id, { is_favorite: false })
      ElMessage.success('已取消收藏')
      loadFavorites()
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    }
  }

  const moveToFolder = async () => {
    if (!movingItem.value) return
    try {
      await xhsApi.updateFavorite(movingItem.value.id, {
        is_favorite: true,
        folder_id: targetFolderId.value
      })
      ElMessage.success('移动成功')
      moveDialogVisible.value = false
      loadFavorites()
    } catch (error: any) {
      ElMessage.error(error.message || '移动失败')
    }
  }

  onMounted(() => {
    loadFolders()
    loadFavorites()
  })
</script>

<style scoped lang="scss">
  .xhs-favorites-page {
    padding: 0;
  }

  .toolbar-card {
    margin-bottom: 20px;
  }

  .folder-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .mr-1 {
    margin-right: 4px;
  }

  .mb-4 {
    margin-bottom: 16px;
  }

  .favorite-card {
    height: 100%;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .card-date {
        font-size: 12px;
        color: var(--art-gray-500);
      }
    }

    .card-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--art-gray-900);
      margin: 0 0 10px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .card-content {
      font-size: 13px;
      color: var(--art-gray-600);
      line-height: 1.6;
      margin-bottom: 12px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .card-tags {
      margin-bottom: 0;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }

  .folder-radio-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
</style>
