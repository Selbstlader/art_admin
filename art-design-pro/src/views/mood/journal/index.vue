<template>
  <div class="journal-page">
    <!-- 搜索和操作栏 -->
    <ElCard shadow="never" class="mb-4">
      <div class="toolbar">
        <div class="search-area">
          <ElInput v-model="searchKeyword" placeholder="搜索日记..." clearable style="width: 250px" @keyup.enter="handleSearch">
            <template #prefix>
              <ElIcon><Search /></ElIcon>
            </template>
          </ElInput>
          <ElSelect v-model="selectedTag" placeholder="选择标签" clearable style="width: 150px" @change="handleSearch">
            <ElOption v-for="tag in allTags" :key="tag" :label="tag" :value="tag" />
          </ElSelect>
          <ElDatePicker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" @change="handleSearch" />
        </div>
        <ElButton type="primary" :icon="Plus" @click="openEditor()">写日记</ElButton>
      </div>
    </ElCard>

    <!-- 日记列表 -->
    <div v-loading="loading" class="journal-list">
      <ElRow :gutter="16">
        <ElCol v-for="journal in journalList" :key="journal.id" :xs="24" :sm="12" :lg="8" class="mb-4">
          <ElCard shadow="hover" class="journal-card" @click="viewJournal(journal)">
            <div class="journal-header">
              <span class="journal-date">{{ formatDate(journal.created_at) }}</span>
              <div class="journal-mood" v-if="journal.mood_type">
                <span>{{ MoodTypeConfig[journal.mood_type]?.icon }}</span>
              </div>
            </div>
            <div class="journal-title">{{ journal.title || '无标题' }}</div>
            <div class="journal-content">{{ getPreview(journal.content) }}</div>
            <div class="journal-tags" v-if="journal.tags?.length">
              <ElTag v-for="tag in journal.tags.slice(0, 3)" :key="tag" size="small" class="mr-1">{{ tag }}</ElTag>
              <span v-if="journal.tags.length > 3" class="more-tags">+{{ journal.tags.length - 3 }}</span>
            </div>
            <div class="journal-footer">
              <span class="word-count">{{ journal.content?.length || 0 }} 字</span>
              <div class="actions" @click.stop>
                <ElButton link type="primary" size="small" @click="openEditor(journal)">编辑</ElButton>
                <ElButton link type="danger" size="small" @click="deleteJournal(journal)">删除</ElButton>
              </div>
            </div>
          </ElCard>
        </ElCol>
      </ElRow>

      <ElEmpty v-if="!loading && journalList.length === 0" description="还没有日记，开始记录吧" />

      <div v-if="total > pageSize" class="pagination-wrapper">
        <ElPagination v-model:current-page="currentPage" :page-size="pageSize" :total="total" layout="prev, pager, next" @current-change="loadJournals" />
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <ElDialog v-model="editorVisible" :title="editingJournal ? '编辑日记' : '写日记'" width="800px" :close-on-click-modal="false">
      <ElForm :model="journalForm" label-width="80px">
        <ElFormItem label="标题">
          <ElInput v-model="journalForm.title" placeholder="给日记起个标题（可选）" />
        </ElFormItem>
        <ElFormItem label="内容">
          <ElInput v-model="journalForm.content" type="textarea" :rows="12" placeholder="今天发生了什么..." maxlength="10000" show-word-limit />
        </ElFormItem>
        <ElFormItem label="当前心情">
          <div class="mood-selector">
            <div v-for="(config, type) in MoodTypeConfig" :key="type" class="mood-option" :class="{ active: journalForm.mood_type === type }" @click="journalForm.mood_type = type as MoodType">
              <span class="mood-icon">{{ config.icon }}</span>
              <span class="mood-label">{{ config.label }}</span>
            </div>
          </div>
        </ElFormItem>
        <ElFormItem label="标签">
          <div class="tag-input-area">
            <ElTag v-for="tag in journalForm.tags" :key="tag" closable @close="removeTag(tag)" class="mr-2 mb-2">{{ tag }}</ElTag>
            <ElInput v-if="tagInputVisible" ref="tagInputRef" v-model="tagInputValue" size="small" style="width: 100px" @keyup.enter="addTag" @blur="addTag" />
            <ElButton v-else size="small" @click="showTagInput">+ 添加标签</ElButton>
          </div>
          <div class="quick-tags">
            <span class="quick-label">快捷标签:</span>
            <ElTag v-for="tag in quickTags" :key="tag" :type="journalForm.tags.includes(tag) ? 'primary' : 'info'" effect="plain" class="mr-2 cursor-pointer" @click="toggleTag(tag)">{{ tag }}</ElTag>
          </div>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="editorVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveJournal" :loading="saving">保存</ElButton>
      </template>
    </ElDialog>

    <!-- 查看弹窗 -->
    <ElDialog v-model="viewVisible" :title="viewingJournal?.title || '日记详情'" width="700px">
      <div v-if="viewingJournal" class="journal-detail">
        <div class="detail-meta">
          <span class="detail-date">{{ formatDateTime(viewingJournal.created_at) }}</span>
          <span v-if="viewingJournal.mood_type" class="detail-mood">{{ MoodTypeConfig[viewingJournal.mood_type]?.icon }} {{ MoodTypeConfig[viewingJournal.mood_type]?.label }}</span>
        </div>
        <div class="detail-content">{{ viewingJournal.content }}</div>
        <div class="detail-tags" v-if="viewingJournal.tags?.length">
          <ElTag v-for="tag in viewingJournal.tags" :key="tag" size="small" class="mr-2">{{ tag }}</ElTag>
        </div>
      </div>
      <template #footer>
        <ElButton @click="viewVisible = false">关闭</ElButton>
        <ElButton type="primary" @click="openEditor(viewingJournal!); viewVisible = false">编辑</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Search, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { journalApi, MoodTypeConfig, type JournalEntry, type MoodType } from '@/api/mood'
import dayjs from 'dayjs'

defineOptions({ name: 'Journal' })

// 搜索条件
const searchKeyword = ref('')
const selectedTag = ref('')
const dateRange = ref<[string, string] | null>(null)

// 分页
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 数据
const journalList = ref<JournalEntry[]>([])
const allTags = ref<string[]>([])
const loading = ref(false)

// 编辑器
const editorVisible = ref(false)
const editingJournal = ref<JournalEntry | null>(null)
const saving = ref(false)
const journalForm = reactive({
  title: '',
  content: '',
  mood_type: '' as MoodType | '',
  tags: [] as string[]
})

// 标签输入
const tagInputVisible = ref(false)
const tagInputValue = ref('')
const tagInputRef = ref()
const quickTags = ['工作', '生活', '学习', '旅行', '美食', '运动', '阅读', '感悟']

// 查看
const viewVisible = ref(false)
const viewingJournal = ref<JournalEntry | null>(null)

// 方法
const formatDate = (date: string) => dayjs(date).format('MM月DD日')
const formatDateTime = (date: string) => dayjs(date).format('YYYY年MM月DD日 HH:mm')

const getPreview = (content: string) => {
  if (!content) return ''
  return content.length > 100 ? content.slice(0, 100) + '...' : content
}

const handleSearch = () => {
  currentPage.value = 1
  loadJournals()
}

const loadJournals = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (selectedTag.value) params.tag = selectedTag.value
    if (dateRange.value) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }

    const res = await journalApi.getList(params)
    const data = (res as any)?.data || res
    journalList.value = data?.records || data || []
    total.value = data?.total || journalList.value.length
  } catch (error) {
    console.error('加载日记失败:', error)
  } finally {
    loading.value = false
  }
}

const loadTags = async () => {
  try {
    const res = await journalApi.getTags()
    allTags.value = (res as any)?.data || res || []
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

const openEditor = (journal?: JournalEntry | null) => {
  editingJournal.value = journal || null
  if (journal) {
    journalForm.title = journal.title
    journalForm.content = journal.content
    journalForm.mood_type = journal.mood_type || ''
    journalForm.tags = [...(journal.tags || [])]
  } else {
    journalForm.title = ''
    journalForm.content = ''
    journalForm.mood_type = ''
    journalForm.tags = []
  }
  editorVisible.value = true
}

const saveJournal = async () => {
  if (!journalForm.content.trim()) {
    ElMessage.warning('请输入日记内容')
    return
  }

  saving.value = true
  try {
    const data = {
      title: journalForm.title,
      content: journalForm.content,
      mood_type: journalForm.mood_type || undefined,
      tags: journalForm.tags
    }

    if (editingJournal.value) {
      await journalApi.update({ id: editingJournal.value.id, ...data })
      ElMessage.success('日记已更新')
    } else {
      await journalApi.create(data)
      ElMessage.success('日记已保存')
    }

    editorVisible.value = false
    loadJournals()
    loadTags()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const deleteJournal = async (journal: JournalEntry) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇日记吗？', '删除确认', { type: 'warning' })
    await journalApi.delete(journal.id)
    ElMessage.success('删除成功')
    loadJournals()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const viewJournal = (journal: JournalEntry) => {
  viewingJournal.value = journal
  viewVisible.value = true
}

const showTagInput = () => {
  tagInputVisible.value = true
  nextTick(() => tagInputRef.value?.focus())
}

const addTag = () => {
  const tag = tagInputValue.value.trim()
  if (tag && !journalForm.tags.includes(tag)) {
    journalForm.tags.push(tag)
  }
  tagInputVisible.value = false
  tagInputValue.value = ''
}

const removeTag = (tag: string) => {
  const index = journalForm.tags.indexOf(tag)
  if (index > -1) journalForm.tags.splice(index, 1)
}

const toggleTag = (tag: string) => {
  const index = journalForm.tags.indexOf(tag)
  if (index > -1) {
    journalForm.tags.splice(index, 1)
  } else {
    journalForm.tags.push(tag)
  }
}

onMounted(() => {
  loadJournals()
  loadTags()
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/variables.scss' as *;

.journal-page {
  --card-spacing: 20px;
}

.mr-1 {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}

.mb-2 {
  margin-bottom: 8px;
}

// 统一卡片样式
:deep(.el-card) {
  margin-bottom: var(--card-spacing);
  background: var(--art-main-bg-color);
  border-radius: calc(var(--custom-radius) + 4px) !important;
  border: none;
  box-shadow: var(--art-root-card-box-shadow);
}

// 工具栏
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;

  .search-area {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }
}

// 日记卡片
.journal-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;

  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--art-box-shadow);
  }

  .journal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    .journal-date {
      font-size: 13px;
      color: var(--art-gray-500);
    }

    .journal-mood {
      font-size: 20px;
    }
  }

  .journal-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--art-gray-900);
    margin-bottom: 10px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .journal-content {
    font-size: 14px;
    color: var(--art-gray-600);
    line-height: 1.7;
    margin-bottom: 12px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .journal-tags {
    margin-bottom: 12px;

    .more-tags {
      font-size: 12px;
      color: var(--art-gray-500);
    }
  }

  .journal-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid var(--art-border-color);

    .word-count {
      font-size: 12px;
      color: var(--art-gray-500);
    }
  }
}

// 分页
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
  padding: 20px;
  background: var(--art-main-bg-color);
  border-radius: calc(var(--custom-radius) + 4px);
}

// 情绪选择器
.mood-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  .mood-option {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 12px 16px;
    border: 2px solid var(--art-border-color);
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      border-color: rgb(var(--art-primary));
      background: rgb(var(--art-bg-primary));
    }

    &.active {
      border-color: rgb(var(--art-primary));
      background: rgb(var(--art-bg-primary));
    }

    .mood-icon {
      font-size: 24px;
      margin-bottom: 4px;
    }

    .mood-label {
      font-size: 12px;
      color: var(--art-gray-600);
    }
  }
}

// 标签输入
.tag-input-area {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 8px;
}

.quick-tags {
  display: flex;
  align-items: center;
  flex-wrap: wrap;

  .quick-label {
    font-size: 12px;
    color: var(--art-gray-500);
    margin-right: 8px;
  }

  .cursor-pointer {
    cursor: pointer;
  }
}

// 日记详情
.journal-detail {
  .detail-meta {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--art-border-color);

    .detail-date {
      color: var(--art-gray-500);
    }

    .detail-mood {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .detail-content {
    font-size: 15px;
    line-height: 1.8;
    color: var(--art-gray-800);
    white-space: pre-wrap;
    margin-bottom: 20px;
  }

  .detail-tags {
    padding-top: 16px;
    border-top: 1px solid var(--art-border-color);
  }
}

// 响应式
@media screen and (max-width: $device-phone) {
  .journal-page {
    --card-spacing: 15px;
  }
}
</style>
