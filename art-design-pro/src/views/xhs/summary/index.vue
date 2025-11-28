<template>
  <div class="xhs-summary-page">
    <div class="summary-container">
      <!-- 左侧功能区 -->
      <div class="left-panel" :style="{ width: leftWidth + '%' }">
        <!-- 输入模块 -->
        <ElCard shadow="never" class="input-card">
          <div class="card-header">
            <img src="https://www.xiaohongshu.com/favicon.ico" class="xhs-icon" alt="小红书" />
            <span class="card-title">笔记输入</span>
          </div>
          <ElTabs v-model="inputMode" class="input-tabs">
            <ElTabPane label="链接输入" name="link">
              <ElInput
                v-model="noteUrl"
                placeholder="请粘贴小红书笔记链接，如：https://www.xiaohongshu.com/explore/..."
                clearable
                class="mb-3"
              >
                <template #prefix>
                  <ElIcon><Link /></ElIcon>
                </template>
              </ElInput>
              <div class="input-hint">
                <ElIcon><InfoFilled /></ElIcon>
                <span>提示：支持小红书分享链接或网页链接</span>
              </div>
              <div class="link-examples">
                <span class="example-label">示例：</span>
                <code>https://www.xiaohongshu.com/explore/xxx</code>
              </div>
            </ElTabPane>
            <ElTabPane label="内容输入" name="content">
              <ElInput v-model="noteTitle" placeholder="笔记标题（可选）" class="mb-3" />
              <ElInput
                v-model="noteContent"
                type="textarea"
                :rows="8"
                placeholder="请粘贴小红书笔记的文字内容..."
                maxlength="10000"
                show-word-limit
              />
              <div class="input-hint">
                <ElIcon><InfoFilled /></ElIcon>
                <span>提示：从小红书复制笔记文字内容粘贴至此</span>
              </div>
            </ElTabPane>
          </ElTabs>
        </ElCard>

        <!-- 配置模块 -->
        <ElCard shadow="never" class="config-card">
          <div class="card-header">
            <ElIcon>
              <Setting />
            </ElIcon>
            <span class="card-title">总结配置</span>
          </div>
          <ElForm label-position="top" size="default">
            <ElFormItem label="总结风格">
              <ElSelect v-model="summaryStyle" style="width: 100%">
                <ElOption
                  v-for="(config, key) in SummaryStyleConfig"
                  :key="key"
                  :label="config.label"
                  :value="key"
                >
                  <div class="style-option">
                    <span class="style-label">{{ config.label }}</span>
                    <span class="style-desc">{{ config.desc }}</span>
                  </div>
                </ElOption>
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="字数限制">
              <ElSelect v-model="maxLength" style="width: 100%">
                <ElOption label="短（100字内）" :value="100" />
                <ElOption label="中（100-300字）" :value="300" />
                <ElOption label="长（300字以上）" :value="500" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem>
              <div class="switch-item">
                <span>专业分析</span>
                <ElSwitch v-model="enableAnalysis" />
              </div>
              <div class="switch-desc">开启后将进行垂直领域深度解析</div>
            </ElFormItem>
            <ElFormItem>
              <div class="switch-item">
                <span>保存到历史</span>
                <ElSwitch v-model="saveHistory" />
              </div>
            </ElFormItem>
          </ElForm>
        </ElCard>

        <!-- 操作模块 -->
        <ElCard shadow="never" class="action-card">
          <ElButton
            type="primary"
            size="large"
            :icon="MagicStick"
            :loading="loading"
            :disabled="inputMode === 'link' ? !noteUrl.trim() : !noteContent.trim()"
            style="width: 100%"
            @click="handleSummarize"
          >
            {{ loading ? '总结中...' : '开始总结' }}
          </ElButton>
          <ElButton
            v-if="summaryResult"
            size="large"
            style="width: 100%; margin-top: 12px"
            @click="handleReset"
          >
            重新总结
          </ElButton>
        </ElCard>
      </div>

      <!-- 分割线 -->
      <div class="divider" @mousedown="startResize"></div>

      <!-- 右侧结果区 -->
      <div class="right-panel" :style="{ width: 100 - leftWidth + '%' }">
        <!-- 空状态 -->
        <div v-if="!summaryResult && !loading" class="empty-state">
          <div class="empty-icon">
            <ElIcon :size="64">
              <MagicStick />
            </ElIcon>
          </div>
          <h3>等待输入笔记内容</h3>
          <p>粘贴笔记内容后，AI将为您快速生成总结及专业分析</p>
          <p class="hint">支持美妆、数码、旅游、美食等多种品类内容</p>
        </div>

        <!-- 加载状态 -->
        <div v-else-if="loading" class="loading-state">
          <ElIcon class="loading-icon" :size="48">
            <Loading />
          </ElIcon>
          <p class="loading-text">{{ loadingText }}</p>
        </div>

        <!-- 结果展示 -->
        <div v-else-if="summaryResult" class="result-container">
          <!-- 顶部标签栏 -->
          <ElCard shadow="never" class="result-header">
            <div class="header-info">
              <h3 class="result-title">{{ getResultTitle(summaryResult) }}</h3>
              <div class="result-meta">
                <ElTag :type="getSentimentType(summaryResult.sentiment)" size="small">
                  {{ SentimentConfig[summaryResult.sentiment]?.icon }}
                  {{ SentimentConfig[summaryResult.sentiment]?.label }}
                </ElTag>
                <span class="time">{{ formatTime(new Date()) }}</span>
              </div>
            </div>
            <ElTabs v-model="activeTab" class="result-tabs">
              <ElTabPane label="综合总结" name="summary" />
              <ElTabPane v-if="summaryResult.analysis" label="专业分析" name="analysis" />
            </ElTabs>
          </ElCard>

          <!-- 内容区 -->
          <ElCard shadow="never" class="result-content">
            <!-- 综合总结 -->
            <div v-show="activeTab === 'summary'" class="tab-content">
              <!-- 快速信息栏 -->
              <div class="quick-info-bar" v-if="summaryResult.reading_time || summaryResult.content_type || summaryResult.target_audience">
                <span v-if="summaryResult.reading_time" class="info-item">
                  <span class="info-icon">⏱️</span>{{ summaryResult.reading_time }}
                </span>
                <span v-if="summaryResult.content_type" class="info-item">
                  <span class="info-icon">📂</span>{{ summaryResult.content_type }}
                </span>
                <span v-if="summaryResult.target_audience" class="info-item">
                  <span class="info-icon">👥</span>{{ summaryResult.target_audience }}
                </span>
              </div>

              <!-- 核心观点 -->
              <div class="section">
                <div class="section-title">
                  <span class="icon">📌</span>
                  <span>核心观点</span>
                </div>
                <div class="section-body summary-text">
                  {{ getResultSummary(summaryResult) }}
                </div>
              </div>

              <!-- 快速事实 -->
              <div class="section" v-if="summaryResult.quick_facts && Object.keys(summaryResult.quick_facts).length">
                <div class="section-title">
                  <span class="icon">📋</span>
                  <span>快速事实</span>
                </div>
                <div class="section-body">
                  <div class="quick-facts-grid">
                    <div v-for="(value, key) in summaryResult.quick_facts" :key="key" class="fact-item">
                      <span class="fact-label">{{ key }}</span>
                      <span class="fact-value">{{ value }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 关键信息 -->
              <div class="section" v-if="summaryResult.key_points?.length">
                <div class="section-title">
                  <span class="icon">🔍</span>
                  <span>关键信息</span>
                </div>
                <div class="section-body">
                  <ul class="key-points">
                    <li v-for="(point, index) in summaryResult.key_points" :key="index">
                      <span class="point-marker">{{ index + 1 }}</span>
                      <span class="point-text">{{ point }}</span>
                    </li>
                  </ul>
                </div>
              </div>

              <!-- 亮点与优势 -->
              <div class="section" v-if="summaryResult.highlights?.length">
                <div class="section-title highlight-title">
                  <span class="icon">✨</span>
                  <span>亮点与优势</span>
                </div>
                <div class="section-body">
                  <div class="highlight-list">
                    <div v-for="(item, index) in summaryResult.highlights" :key="index" class="highlight-item">
                      <span class="highlight-icon">👍</span>
                      <span>{{ item }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 注意事项 -->
              <div class="section" v-if="summaryResult.concerns?.length">
                <div class="section-title warning-title">
                  <span class="icon">⚠️</span>
                  <span>注意事项</span>
                </div>
                <div class="section-body">
                  <div class="concern-list">
                    <div v-for="(item, index) in summaryResult.concerns" :key="index" class="concern-item">
                      {{ item }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- 可执行建议 -->
              <div class="section" v-if="summaryResult.action_items?.length">
                <div class="section-title">
                  <span class="icon">✅</span>
                  <span>可执行建议</span>
                </div>
                <div class="section-body">
                  <ul class="action-list">
                    <li v-for="(item, index) in summaryResult.action_items" :key="index">
                      {{ item }}
                    </li>
                  </ul>
                </div>
              </div>

              <!-- 可信度说明 -->
              <div class="section" v-if="summaryResult.credibility_note">
                <div class="section-title">
                  <span class="icon">🔒</span>
                  <span>可信度评估</span>
                </div>
                <div class="section-body credibility-note">
                  {{ summaryResult.credibility_note }}
                </div>
              </div>

              <!-- 内容标签 -->
              <div class="section" v-if="summaryResult.tags?.length">
                <div class="section-title">
                  <span class="icon">🏷️</span>
                  <span>内容标签</span>
                </div>
                <div class="section-body">
                  <ElTag
                    v-for="tag in summaryResult.tags"
                    :key="tag"
                    class="mr-2 mb-2"
                    effect="plain"
                  >
                    {{ tag }}
                  </ElTag>
                </div>
              </div>

              <!-- 相关话题 -->
              <div class="section" v-if="summaryResult.related_topics?.length">
                <div class="section-title">
                  <span class="icon">🔗</span>
                  <span>相关话题</span>
                </div>
                <div class="section-body">
                  <ElTag
                    v-for="topic in summaryResult.related_topics"
                    :key="topic"
                    class="mr-2 mb-2"
                    type="info"
                    effect="plain"
                  >
                    #{{ topic }}
                  </ElTag>
                </div>
              </div>
            </div>

            <!-- 专业分析 -->
            <div
              v-show="activeTab === 'analysis'"
              class="tab-content"
              v-if="summaryResult.analysis"
            >
              <div class="analysis-card conclusion-card">
                <div class="analysis-category">
                  <ElTag type="primary">{{
                    getCategoryLabel(summaryResult.analysis.category)
                  }}</ElTag>
                </div>
                <div class="analysis-conclusion">
                  {{ summaryResult.analysis.main_conclusion }}
                </div>
              </div>

              <div
                class="section"
                v-if="Object.keys(summaryResult.analysis.detailed_data || {}).length"
              >
                <div class="section-title">
                  <span class="icon">📊</span>
                  <span>详细数据</span>
                </div>
                <div class="section-body">
                  <div class="data-grid">
                    <div
                      v-for="(value, key) in summaryResult.analysis.detailed_data"
                      :key="key"
                      class="data-item"
                    >
                      <span class="data-label">{{ key }}</span>
                      <span class="data-value">{{ value }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="section" v-if="summaryResult.analysis.risk_warnings?.length">
                <div class="section-title warning">
                  <span class="icon">⚠️</span>
                  <span>风险提示</span>
                </div>
                <div class="section-body">
                  <div class="warning-list">
                    <div
                      v-for="(warning, index) in summaryResult.analysis.risk_warnings"
                      :key="index"
                      class="warning-item"
                    >
                      {{ warning }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="section" v-if="summaryResult.analysis.recommendations?.length">
                <div class="section-title">
                  <span class="icon">💡</span>
                  <span>建议</span>
                </div>
                <div class="section-body">
                  <ul class="recommendations">
                    <li v-for="(rec, index) in summaryResult.analysis.recommendations" :key="index">
                      {{ rec }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 操作栏 -->
          <ElCard shadow="never" class="result-actions">
            <div class="action-buttons">
              <ElButton :icon="DocumentCopy" @click="copyResult">复制总结</ElButton>
              <ElButton
                v-if="savedId"
                :icon="isFavorite ? StarFilled : Star"
                @click="toggleFavorite"
              >
                {{ isFavorite ? '取消收藏' : '收藏' }}
              </ElButton>
              <ElButton v-if="savedId" :icon="Delete" type="danger" plain @click="deleteResult">
                删除
              </ElButton>
            </div>
            <div class="feedback-buttons">
              <ElButton text @click="handleFeedback('satisfied')">👍 满意</ElButton>
              <ElButton text @click="handleFeedback('unsatisfied')">👎 不满意</ElButton>
            </div>
          </ElCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onUnmounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    InfoFilled,
    Setting,
    MagicStick,
    Loading,
    DocumentCopy,
    Star,
    StarFilled,
    Delete,
    Link
  } from '@element-plus/icons-vue'
  import {
    xhsApi,
    SummaryStyleConfig,
    SentimentConfig,
    type SummaryStyle,
    type Sentiment,
    type SummaryResult,
    type XHSSummary
  } from '@/api/xhs'
  import dayjs from 'dayjs'

  defineOptions({ name: 'XHSSummary' })

  // 布局
  const leftWidth = ref(35)
  let isResizing = false

  const startResize = (e: MouseEvent) => {
    isResizing = true
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
  }

  const handleResize = (e: MouseEvent) => {
    if (!isResizing) return
    const container = document.querySelector('.summary-container') as HTMLElement
    if (!container) return
    const rect = container.getBoundingClientRect()
    const newWidth = ((e.clientX - rect.left) / rect.width) * 100
    leftWidth.value = Math.max(25, Math.min(50, newWidth))
  }

  const stopResize = () => {
    isResizing = false
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
  }

  onUnmounted(() => {
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
  })

  // 输入模式
  const inputMode = ref<'link' | 'content'>('link')
  const noteUrl = ref('')

  // 表单数据
  const noteTitle = ref('')
  const noteContent = ref('')
  const summaryStyle = ref<SummaryStyle>('concise')
  const maxLength = ref(300)
  const enableAnalysis = ref(true)
  const saveHistory = ref(true)

  // 状态
  const loading = ref(false)
  const loadingText = ref('正在解析内容...')
  const activeTab = ref('summary')
  const summaryResult = ref<SummaryResult | XHSSummary | null>(null)
  const savedId = ref<number | null>(null)
  const isFavorite = ref(false)

  // 加载动画文本
  const loadingTexts = ['正在解析内容...', 'AI正在提炼核心信息...', '正在生成总结...']
  let loadingTextIndex = 0
  let loadingTimer: number | null = null

  const startLoadingAnimation = () => {
    loadingTextIndex = 0
    loadingText.value = loadingTexts[0]
    loadingTimer = window.setInterval(() => {
      loadingTextIndex = (loadingTextIndex + 1) % loadingTexts.length
      loadingText.value = loadingTexts[loadingTextIndex]
    }, 2000)
  }

  const stopLoadingAnimation = () => {
    if (loadingTimer) {
      clearInterval(loadingTimer)
      loadingTimer = null
    }
  }

  // 方法
  const handleSummarize = async () => {
    // 根据输入模式验证
    if (inputMode.value === 'link') {
      if (!noteUrl.value.trim()) {
        ElMessage.warning('请输入小红书链接')
        return
      }
      // 简单验证链接格式
      if (!noteUrl.value.includes('xiaohongshu.com') && !noteUrl.value.includes('xhslink.com')) {
        ElMessage.warning('请输入有效的小红书链接')
        return
      }
    } else {
      if (!noteContent.value.trim()) {
        ElMessage.warning('请输入笔记内容')
        return
      }
    }

    loading.value = true
    startLoadingAnimation()
    summaryResult.value = null
    savedId.value = null
    isFavorite.value = false

    try {
      const requestData: any = {
        style: summaryStyle.value,
        max_length: maxLength.value,
        enable_analysis: enableAnalysis.value,
        save_history: saveHistory.value
      }

      // 根据输入模式设置不同参数
      if (inputMode.value === 'link') {
        requestData.original_url = noteUrl.value
      } else {
        requestData.note_title = noteTitle.value
        requestData.content = noteContent.value
      }

      const res = await xhsApi.summarize(requestData)

      const data = (res as any)?.data || res
      summaryResult.value = data

      // 如果保存了历史记录，获取ID
      if (saveHistory.value && data.id) {
        savedId.value = data.id
        isFavorite.value = data.is_favorite || false
      }

      activeTab.value = 'summary'
      ElMessage.success('总结生成成功')
    } catch (error: any) {
      ElMessage.error(error.message || '总结生成失败')
    } finally {
      loading.value = false
      stopLoadingAnimation()
    }
  }

  const handleReset = () => {
    summaryResult.value = null
    savedId.value = null
    isFavorite.value = false
    activeTab.value = 'summary'
    noteUrl.value = ''
    noteTitle.value = ''
    noteContent.value = ''
  }

  const copyResult = () => {
    if (!summaryResult.value) return

    const text = `【${getResultTitle(summaryResult.value)}】\n\n${getResultSummary(summaryResult.value)}\n\n关键信息：\n${summaryResult.value.key_points?.map((p, i) => `${i + 1}. ${p}`).join('\n') || ''}`

    navigator.clipboard
      .writeText(text)
      .then(() => {
        ElMessage.success('已复制到剪贴板')
      })
      .catch(() => {
        ElMessage.error('复制失败')
      })
  }

  const toggleFavorite = async () => {
    if (!savedId.value) return

    try {
      await xhsApi.updateFavorite(savedId.value, {
        is_favorite: !isFavorite.value
      })
      isFavorite.value = !isFavorite.value
      ElMessage.success(isFavorite.value ? '已收藏' : '已取消收藏')
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    }
  }

  const deleteResult = async () => {
    if (!savedId.value) return

    try {
      await ElMessageBox.confirm('确定要删除这条总结吗？', '删除确认', { type: 'warning' })
      await xhsApi.deleteSummary(savedId.value)
      ElMessage.success('删除成功')
      handleReset()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  const handleFeedback = (type: 'satisfied' | 'unsatisfied') => {
    if (type === 'satisfied') {
      ElMessage.success('感谢您的反馈！')
    } else {
      ElMessage.info('感谢反馈，我们会继续优化')
    }
  }

  // 工具函数
  const formatTime = (date: Date) => dayjs(date).format('YYYY-MM-DD HH:mm')

  // 获取结果标题（兼容两种返回类型）
  const getResultTitle = (result: any) => {
    return result?.title || result?.summary_title || '总结结果'
  }

  // 获取结果摘要（兼容两种返回类型）
  const getResultSummary = (result: any) => {
    return result?.summary || result?.summary_content || ''
  }

  const getSentimentType = (sentiment: Sentiment) => {
    const map: Record<Sentiment, '' | 'success' | 'warning' | 'info' | 'danger'> = {
      positive: 'success',
      neutral: 'info',
      negative: 'warning'
    }
    return map[sentiment] || 'info'
  }

  const getCategoryLabel = (category: string) => {
    const map: Record<string, string> = {
      beauty: '美妆护肤',
      digital: '数码产品',
      food: '美食',
      travel: '旅游',
      lifestyle: '生活方式',
      other: '其他'
    }
    return map[category] || category
  }
</script>

<style scoped lang="scss">
  .xhs-summary-page {
    height: 100%;
    padding: 20px;
  }

  .summary-container {
    display: flex;
    height: calc(100vh - 140px);
    min-height: 600px;
    background: var(--art-main-bg-color);
    border-radius: var(--custom-radius);
    overflow: hidden;
  }

  // 左侧面板
  .left-panel {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 20px;
    overflow-y: auto;
    background: var(--art-sider-bg-color);
  }

  // 分割线
  .divider {
    width: 4px;
    background: var(--art-border-color);
    cursor: col-resize;
    transition: background 0.3s;

    &:hover {
      background: rgb(var(--art-primary));
    }
  }

  // 右侧面板
  .right-panel {
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow-y: auto;
  }

  // 卡片样式
  :deep(.el-card) {
    border: none;
    border-radius: calc(var(--custom-radius) + 4px);
    box-shadow: var(--art-root-card-box-shadow);
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--art-border-color);

    .xhs-icon {
      width: 24px;
      height: 24px;
    }

    .card-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--art-gray-900);
    }
  }

  .input-hint {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-top: 8px;
    font-size: 12px;
    color: var(--art-gray-500);
  }

  // 输入模式Tab
  .input-tabs {
    :deep(.el-tabs__header) {
      margin-bottom: 16px;
    }

    :deep(.el-tabs__content) {
      overflow: visible;
    }
  }

  .link-examples {
    margin-top: 12px;
    padding: 12px;
    background: var(--art-main-bg-color);
    border-radius: 8px;
    font-size: 12px;

    .example-label {
      color: var(--art-gray-500);
      margin-right: 8px;
    }

    code {
      color: var(--art-gray-600);
      background: var(--art-sider-bg-color);
      padding: 2px 6px;
      border-radius: 4px;
      font-family: monospace;
      word-break: break-all;
    }
  }

  .mb-3 {
    margin-bottom: 12px;
  }

  .mr-2 {
    margin-right: 8px;
  }

  .mb-2 {
    margin-bottom: 8px;
  }

  // 风格选项
  .style-option {
    display: flex;
    flex-direction: column;

    .style-label {
      font-size: 14px;
    }

    .style-desc {
      font-size: 12px;
      color: var(--art-gray-500);
    }
  }

  // 开关项
  .switch-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }

  .switch-desc {
    font-size: 12px;
    color: var(--art-gray-500);
    margin-top: 4px;
  }

  // 空状态
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;

    .empty-icon {
      color: var(--art-gray-300);
      margin-bottom: 24px;
    }

    h3 {
      font-size: 20px;
      color: var(--art-gray-800);
      margin-bottom: 12px;
    }

    p {
      font-size: 14px;
      color: var(--art-gray-500);
      margin-bottom: 8px;
    }

    .hint {
      font-size: 12px;
      color: var(--art-gray-400);
    }
  }

  // 加载状态
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;

    .loading-icon {
      color: rgb(var(--art-primary));
      animation: rotate 1.5s linear infinite;
    }

    .loading-text {
      margin-top: 16px;
      font-size: 16px;
      color: var(--art-gray-600);
    }
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }

    to {
      transform: rotate(360deg);
    }
  }

  // 结果容器
  .result-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
  }

  .result-header {
    flex-shrink: 0;

    .header-info {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;
    }

    .result-title {
      font-size: 18px;
      font-weight: 600;
      color: var(--art-gray-900);
      margin: 0;
    }

    .result-meta {
      display: flex;
      align-items: center;
      gap: 12px;

      .time {
        font-size: 12px;
        color: var(--art-gray-500);
      }
    }

    .result-tabs {
      margin-bottom: -20px;

      :deep(.el-tabs__header) {
        margin-bottom: 0;
      }
    }
  }

  .result-content {
    flex: 1;
    overflow-y: auto;
  }

  .tab-content {
    padding: 8px 0;
  }

  // 快速信息栏
  .quick-info-bar {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    padding: 12px 16px;
    background: var(--art-main-bg-color);
    border-radius: 8px;
    margin-bottom: 20px;

    .info-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--art-gray-600);

      .info-icon {
        font-size: 14px;
      }
    }
  }

  // 快速事实
  .quick-facts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;

    .fact-item {
      display: flex;
      flex-direction: column;
      padding: 10px 12px;
      background: var(--art-main-bg-color);
      border-radius: 8px;
      border-left: 3px solid rgb(var(--art-primary));

      .fact-label {
        font-size: 12px;
        color: var(--art-gray-500);
        margin-bottom: 4px;
      }

      .fact-value {
        font-size: 14px;
        font-weight: 500;
        color: var(--art-gray-800);
      }
    }
  }

  // 亮点列表
  .highlight-list {
    .highlight-item {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      padding: 10px 12px;
      background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(16, 185, 129, 0.05));
      border-radius: 8px;
      margin-bottom: 8px;
      font-size: 14px;
      color: var(--art-gray-700);

      &:last-child {
        margin-bottom: 0;
      }

      .highlight-icon {
        flex-shrink: 0;
      }
    }
  }

  .highlight-title {
    color: #10b981;
  }

  // 注意事项
  .concern-list {
    .concern-item {
      padding: 10px 12px;
      background: #fef3c7;
      border-left: 3px solid #f59e0b;
      border-radius: 0 8px 8px 0;
      margin-bottom: 8px;
      font-size: 14px;
      color: #92400e;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .warning-title {
    color: #f59e0b;
  }

  // 可执行建议
  .action-list {
    list-style: none;
    padding: 0;
    margin: 0;

    li {
      position: relative;
      padding: 8px 0 8px 24px;
      font-size: 14px;
      line-height: 1.6;
      color: var(--art-gray-700);

      &::before {
        content: '✓';
        position: absolute;
        left: 0;
        color: #10b981;
        font-weight: bold;
      }
    }
  }

  // 可信度说明
  .credibility-note {
    padding: 12px 16px;
    background: var(--art-main-bg-color);
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
    color: var(--art-gray-600);
    border-left: 3px solid var(--art-gray-400);
  }

  .section {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    font-weight: 600;
    color: var(--art-gray-800);
    margin-bottom: 12px;

    .icon {
      font-size: 18px;
    }

    &.warning {
      color: #f59e0b;
    }
  }

  .section-body {
    padding-left: 26px;
  }

  .summary-text {
    font-size: 15px;
    line-height: 1.8;
    color: var(--art-gray-700);
  }

  .key-points {
    list-style: none;
    padding: 0;
    margin: 0;

    li {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }
    }

    .point-marker {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      background: rgb(var(--art-primary));
      color: #fff;
      border-radius: 50%;
      font-size: 12px;
      flex-shrink: 0;
    }

    .point-text {
      font-size: 14px;
      line-height: 24px;
      color: var(--art-gray-700);
    }
  }

  // 专业分析
  .analysis-card {
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 20px;
  }

  .conclusion-card {
    background: linear-gradient(
      135deg,
      rgba(var(--art-primary), 0.1),
      rgba(var(--art-primary), 0.05)
    );
    border: 1px solid rgba(var(--art-primary), 0.2);

    .analysis-category {
      margin-bottom: 12px;
    }

    .analysis-conclusion {
      font-size: 15px;
      line-height: 1.8;
      color: var(--art-gray-800);
    }
  }

  .data-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 12px;
  }

  .data-item {
    display: flex;
    flex-direction: column;
    padding: 12px;
    background: var(--art-sider-bg-color);
    border-radius: 8px;

    .data-label {
      font-size: 12px;
      color: var(--art-gray-500);
      margin-bottom: 4px;
    }

    .data-value {
      font-size: 14px;
      font-weight: 500;
      color: var(--art-gray-800);
    }
  }

  .warning-list {
    .warning-item {
      padding: 12px;
      background: #fef3c7;
      border-left: 3px solid #f59e0b;
      border-radius: 0 8px 8px 0;
      margin-bottom: 8px;
      font-size: 14px;
      color: #92400e;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .recommendations {
    list-style: disc;
    padding-left: 20px;
    margin: 0;

    li {
      font-size: 14px;
      line-height: 1.8;
      color: var(--art-gray-700);
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  // 操作栏
  .result-actions {
    flex-shrink: 0;

    :deep(.el-card__body) {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .action-buttons {
      display: flex;
      gap: 12px;
    }

    .feedback-buttons {
      display: flex;
      gap: 8px;
    }
  }

  // 响应式
  @media screen and (max-width: 1200px) {
    .summary-container {
      flex-direction: column;
      height: auto;
    }

    .left-panel,
    .right-panel {
      width: 100% !important;
    }

    .divider {
      display: none;
    }

    .right-panel {
      min-height: 500px;
    }
  }
</style>
