<template>
  <div class="learning-list-page art-full-height">
    <ElCard shadow="never" class="list-container">
      <div class="content-wrapper">
        <!-- 左侧：学科分类 -->
        <div class="subject-sidebar">
          <div class="sidebar-title">课程分类</div>
          <ElMenu
            :default-active="activeSubjectId?.toString()"
            @select="handleSubjectChange"
            class="subject-menu"
          >
            <ElMenuItem index="all">
              <ElIcon><List /></ElIcon>
              <span>全部学科</span>
            </ElMenuItem>
            <ElMenuItem
              v-for="subject in subjects"
              :key="subject.id"
              :index="subject.id.toString()"
            >
              <span class="subject-icon">{{ subject.icon }}</span>
              <span>{{ subject.name }}</span>
            </ElMenuItem>
          </ElMenu>
        </div>

        <!-- 右侧：年级卡片列表 -->
        <div class="grade-content">
          <div class="grade-header">
            <h2>{{ currentSubjectName }}</h2>
            <p>选择年级查看对应课程</p>
          </div>

          <div v-if="loading" class="loading-wrapper">
            <ElSkeleton :rows="5" animated />
          </div>

          <ElEmpty v-else-if="grades.length === 0" description="暂无课程，快去生成吧！">
            <ElButton type="primary" @click="$router.push('/learning/learning/my')"
              >立即生成</ElButton
            >
          </ElEmpty>

          <div v-else class="grade-grid">
            <ElCard
              v-for="grade in grades"
              :key="grade.grade"
              shadow="hover"
              class="grade-card"
              @click="goToCourses(grade.grade)"
            >
              <div class="grade-icon">📚</div>
              <div class="grade-name">{{ grade.grade }}</div>
              <div class="grade-count">{{ grade.count }} 门课程</div>
              <ElIcon class="arrow-icon"><ArrowRight /></ElIcon>
            </ElCard>
          </div>
        </div>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { List, ArrowRight } from '@element-plus/icons-vue'
  import { subjectApi, gradeApi } from '@/api/learning'

  const router = useRouter()
  const loading = ref(false)
  const subjects = ref<any[]>([])
  const grades = ref<any[]>([])
  const activeSubjectId = ref<number | null>(null)

  const currentSubjectName = computed(() => {
    if (!activeSubjectId.value) return '全部学科'
    const subject = subjects.value.find((s) => s.id === activeSubjectId.value)
    return subject ? subject.name : '全部学科'
  })

  const loadSubjects = async () => {
    try {
      const res: any = await subjectApi.getSubjectList()
      subjects.value = res || []
      if (subjects.value.length > 0) {
        activeSubjectId.value = subjects.value[0].id
        loadGrades()
      }
    } catch (error) {
      console.error('获取学科列表失败:', error)
    }
  }

  const loadGrades = async () => {
    try {
      loading.value = true
      const res: any = await gradeApi.getGradesBySubject(activeSubjectId.value || undefined)
      grades.value = res.grades || []
    } catch (error) {
      console.error('获取年级列表失败:', error)
      grades.value = []
    } finally {
      loading.value = false
    }
  }

  const handleSubjectChange = (index: string) => {
    if (index === 'all') {
      activeSubjectId.value = null
    } else {
      activeSubjectId.value = parseInt(index)
    }
    loadGrades()
  }

  const goToCourses = (grade: string) => {
    router.push({
      path: '/learning/learning/courses',
      query: {
        subject_id: activeSubjectId.value || '',
        grade: grade
      }
    })
  }

  onMounted(() => {
    loadSubjects()
  })
</script>

<style scoped lang="scss">
  .learning-list-page {
    .list-container {
      height: 100%;
      display: flex;
      flex-direction: column;
    }

    .content-wrapper {
      display: flex;
      gap: 20px;
      flex: 1;
      overflow: hidden;
    }

    .subject-sidebar {
      width: 240px;
      flex-shrink: 0;
      border-right: 1px solid var(--el-border-color-light);
      padding-right: 20px;

      .sidebar-title {
        font-size: 16px;
        font-weight: 600;
        padding: 12px 16px;
        color: var(--el-text-color-primary);
      }

      .subject-menu {
        border: none;

        .subject-icon {
          font-size: 18px;
          margin-right: 8px;
        }
      }
    }

    .grade-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;

      .grade-header {
        margin-bottom: 24px;

        h2 {
          font-size: 24px;
          margin: 0 0 8px 0;
          color: var(--el-text-color-primary);
        }

        p {
          margin: 0;
          color: var(--el-text-color-secondary);
        }
      }

      .loading-wrapper {
        padding: 20px;
      }

      .grade-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 20px;
        overflow-y: auto;
        padding: 4px;

        .grade-card {
          cursor: pointer;
          transition: all 0.3s;
          position: relative;
          text-align: center;
          padding: 32px 20px;

          &:hover {
            transform: translateY(-8px);
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);

            .arrow-icon {
              opacity: 1;
              transform: translateX(4px);
            }
          }

          .grade-icon {
            font-size: 48px;
            margin-bottom: 16px;
          }

          .grade-name {
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 8px;
            color: var(--el-text-color-primary);
          }

          .grade-count {
            font-size: 14px;
            color: var(--el-text-color-secondary);
          }

          .arrow-icon {
            position: absolute;
            right: 16px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 20px;
            color: var(--el-color-primary);
            opacity: 0;
            transition: all 0.3s;
          }
        }
      }
    }
  }
</style>
