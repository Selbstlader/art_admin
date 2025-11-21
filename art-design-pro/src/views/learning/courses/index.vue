<template>
  <div class="courses-page art-full-height">
    <ElCard shadow="never">
      <!-- 页面头部 -->
      <div class="page-header">
        <ElPageHeader @back="$router.back()">
          <template #content>
            <span class="header-title">{{ subjectName }} - {{ grade }}</span>
          </template>
          <template #extra>
            <ElButton type="primary" @click="$router.push('/learning/learning/my')">
              <ElIcon><Plus /></ElIcon>
              生成新课程
            </ElButton>
          </template>
        </ElPageHeader>
      </div>

      <ElDivider />

      <!-- 搜索框 -->
      <div class="search-bar">
        <ElInput
          v-model="searchKeyword"
          placeholder="搜索课程标题或题材..."
          clearable
          size="large"
          style="max-width: 600px"
        >
          <template #prefix>
            <ElIcon><Search /></ElIcon>
          </template>
        </ElInput>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-wrapper">
        <ElSkeleton :rows="5" animated />
      </div>

      <!-- 空状态 -->
      <ElEmpty v-else-if="filteredCourses.length === 0" description="暂无课程">
        <ElButton type="primary" @click="$router.push('/learning/learning/my')">立即生成</ElButton>
      </ElEmpty>

      <!-- 课程列表 -->
      <div v-else class="courses-content">
        <!-- AI生成的课程 -->
        <div class="section">
          <div class="section-header">
            <h2>
              <ElIcon><MagicStick /></ElIcon>
              AI生成的课程
            </h2>
            <span class="count">{{ aiCourses.length }} 门</span>
          </div>
          <div class="course-grid">
            <ElCard
              v-for="course in aiCourses"
              :key="course.id"
              shadow="hover"
              class="course-card"
              @click="goToDetail(course.id)"
            >
              <div class="course-header">
                <div class="course-title">{{ course.title }}</div>
                <ElDropdown trigger="click" @command="(cmd) => handleCommand(cmd, course)">
                  <ElIcon class="more-icon" @click.stop><MoreFilled /></ElIcon>
                  <template #dropdown>
                    <ElDropdownMenu>
                      <ElDropdownItem command="view">查看详情</ElDropdownItem>
                      <ElDropdownItem command="favorite">
                        {{ course.is_favorite ? '取消收藏' : '收藏' }}
                      </ElDropdownItem>
                      <ElDropdownItem command="delete" divided>删除</ElDropdownItem>
                    </ElDropdownMenu>
                  </template>
                </ElDropdown>
              </div>

              <div class="course-info">
                <ElSpace wrap>
                  <ElTag size="small">{{ course.grade }}</ElTag>
                  <ElTag size="small" :type="getDifficultyType(course.difficulty)">
                    {{ getDifficultyText(course.difficulty) }}
                  </ElTag>
                  <ElTag size="small" type="info">
                    <ElIcon><Clock /></ElIcon>
                    {{ course.total_time }}分钟
                  </ElTag>
                </ElSpace>
              </div>

              <div class="course-summary">{{ course.summary || course.topic }}</div>

              <div class="course-footer">
                <div class="course-stats">
                  <span>
                    <ElIcon><View /></ElIcon>
                    {{ course.view_count || 0 }}
                  </span>
                  <span>
                    <ElIcon><Star /></ElIcon>
                    {{ course.favorite_count || 0 }}
                  </span>
                </div>
                <div class="course-time">{{ formatTime(course.created_at) }}</div>
              </div>
            </ElCard>
          </div>
        </div>

        <!-- 收藏的课程 -->
        <div v-if="favoriteCourses.length > 0" class="section">
          <div class="section-header">
            <h2>
              <ElIcon><Star /></ElIcon>
              收藏的课程
            </h2>
            <span class="count">{{ favoriteCourses.length }} 门</span>
          </div>
          <div class="course-grid">
            <ElCard
              v-for="course in favoriteCourses"
              :key="course.id"
              shadow="hover"
              class="course-card"
              @click="goToDetail(course.id)"
            >
              <!-- 同上，省略重复代码 -->
              <div class="course-title">{{ course.title }}</div>
            </ElCard>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <ElPagination
        v-if="filteredCourses.length > 0"
        v-model:current-page="pagination.current"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="loadCourses"
        @current-change="loadCourses"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Plus, Search, MagicStick, MoreFilled, Clock, View, Star } from '@element-plus/icons-vue'
  import { subjectApi, learningMaterialApi, type LearningMaterial } from '@/api/learning'
  import dayjs from 'dayjs'

  const route = useRoute()
  const router = useRouter()
  const loading = ref(false)
  const searchKeyword = ref('')
  const materials = ref<LearningMaterial[]>([])
  const subjectName = ref('')
  const grade = ref('')

  const pagination = reactive({
    current: 1,
    pageSize: 12,
    total: 0
  })

  // 过滤后的课程
  const filteredCourses = computed(() => {
    if (!searchKeyword.value) return materials.value
    const keyword = searchKeyword.value.toLowerCase()
    return materials.value.filter(
      (m) => m.title.toLowerCase().includes(keyword) || m.topic.toLowerCase().includes(keyword)
    )
  })

  // AI生成的课程（全部）
  const aiCourses = computed(() => {
    return filteredCourses.value
  })

  // 收藏的课程（需要后端支持收藏功能）
  const favoriteCourses = computed(() => {
    return filteredCourses.value.filter((m: any) => m.is_favorite)
  })

  const loadSubjectName = async () => {
    try {
      const res: any = await subjectApi.getSubjectList()
      const subject = res.find((s: any) => s.id === parseInt(route.query.subject_id as string))
      subjectName.value = subject ? subject.name : '未知学科'
    } catch (error) {
      console.error('获取学科信息失败:', error)
    }
  }

  const loadCourses = async () => {
    try {
      loading.value = true
      const params: any = {
        page: pagination.current,
        limit: pagination.pageSize
      }

      if (route.query.subject_id) {
        params.subject_id = parseInt(route.query.subject_id as string)
      }

      if (route.query.grade) {
        params.grade = route.query.grade
      }

      const res: any = await learningMaterialApi.getMaterialList(params)
      materials.value = res.list || []
      pagination.total = res.total || 0
    } catch (error) {
      ElMessage.error('加载课程失败')
    } finally {
      loading.value = false
    }
  }

  const goToDetail = (id: number) => {
    router.push(`/learning/learning/detail/${id}`)
  }

  const handleCommand = async (command: string, course: LearningMaterial) => {
    if (command === 'view') {
      goToDetail(course.id)
    } else if (command === 'favorite') {
      ElMessage.info('收藏功能开发中...')
    } else if (command === 'delete') {
      try {
        await ElMessageBox.confirm('确定删除这个课程吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        await learningMaterialApi.deleteMaterial(course.id)
        ElMessage.success('删除成功')
        loadCourses()
      } catch (error: any) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }
  }

  const getDifficultyType = (difficulty: number) => {
    const types: Record<number, any> = { 1: 'success', 2: 'warning', 3: 'danger' }
    return types[difficulty] || ''
  }

  const getDifficultyText = (difficulty: number) => {
    const texts: Record<number, string> = { 1: '⭐ 基础', 2: '⭐⭐ 进阶', 3: '⭐⭐⭐ 高级' }
    return texts[difficulty] || '未知'
  }

  const formatTime = (time: string) => {
    return dayjs(time).format('YYYY-MM-DD HH:mm')
  }

  onMounted(async () => {
    // 先从路由参数中获取年级信息
    if (route.query.grade) {
      grade.value = route.query.grade as string
    }
    await loadSubjectName()
    await loadCourses()
  })
</script>

<style scoped lang="scss">
  .courses-page {
    .page-header {
      margin-bottom: 16px;

      .header-title {
        font-size: 20px;
        font-weight: 600;
      }
    }

    .search-bar {
      margin: 24px 0;
    }

    .loading-wrapper {
      padding: 20px;
    }

    .courses-content {
      .section {
        margin-bottom: 40px;

        &:last-child {
          margin-bottom: 0;
        }

        .section-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 20px;

          h2 {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 18px;
            margin: 0;
            color: var(--el-text-color-primary);
          }

          .count {
            color: var(--el-text-color-secondary);
            font-size: 14px;
          }
        }

        .course-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
          gap: 16px;

          .course-card {
            cursor: pointer;
            transition: all 0.3s;
            height: fit-content;

            &:hover {
              transform: translateY(-4px);
              box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
            }

            .course-header {
              display: flex;
              justify-content: space-between;
              align-items: flex-start;
              margin-bottom: 12px;

              .course-title {
                font-size: 16px;
                font-weight: 600;
                line-height: 1.4;
                flex: 1;
                overflow: hidden;
                text-overflow: ellipsis;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
              }

              .more-icon {
                font-size: 18px;
                cursor: pointer;
                padding: 4px;
                margin-left: 8px;
                color: var(--el-text-color-secondary);

                &:hover {
                  color: var(--el-color-primary);
                }
              }
            }

            .course-info {
              margin-bottom: 12px;
            }

            .course-summary {
              font-size: 14px;
              color: var(--el-text-color-secondary);
              line-height: 1.6;
              margin-bottom: 16px;
              overflow: hidden;
              text-overflow: ellipsis;
              display: -webkit-box;
              -webkit-line-clamp: 2;
              -webkit-box-orient: vertical;
              min-height: 44px;
            }

            .course-footer {
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding-top: 12px;
              border-top: 1px solid var(--el-border-color-lighter);
              font-size: 13px;
              color: var(--el-text-color-secondary);

              .course-stats {
                display: flex;
                gap: 16px;

                span {
                  display: flex;
                  align-items: center;
                  gap: 4px;
                }
              }
            }
          }
        }
      }
    }

    .pagination {
      margin-top: 24px;
      justify-content: flex-end;
    }
  }
</style>
