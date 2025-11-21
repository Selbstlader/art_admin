# AI 学习系统 - 三级导航实现指南

## ✅ 后端已完成

### 1. 新增接口
- `GET /api/learning/grades?subject_id=xxx` - 获取某学科下的年级列表及课程数量
- `GET /api/learning/material/list?subject_id=xxx&grade=xxx&page=1&limit=12` - 支持按年级筛选

### 2. 数据结构
```json
// 年级列表响应
{
  "subject_id": 1,
  "grades": [
    { "grade": "初一", "count": 5 },
    { "grade": "初二", "count": 3 }
  ]
}
```

## 📋 前端需要实现的三级结构

### 第一级：学科分类页面 (`/learning/list`)
**左侧**：学科菜单（语文、数学、英语等）
**右侧**：年级卡片列表（初一、初二、大一、大二等）
- 显示年级名称
- 显示该年级下的课程数量
- 点击年级卡片 → 跳转到第二级

### 第二级：课程列表页面 (`/learning/courses?subject_id=1&grade=初一`)
**顶部**：搜索框（可快捷查询课程）
**内容**：课程卡片列表
- 分类展示：
  - AI生成的课程
  - 收藏的课程
- 点击课程卡片 → 跳转到第三级

### 第三级：课程详情页面 (`/learning/detail/:id`)
显示完整的课程内容

## 🔧 前端实现步骤

### 步骤1：修改 `list/index.vue` - 显示年级列表

```vue
<template>
  <div class="learning-list-page">
    <ElCard>
      <div class="content-wrapper">
        <!-- 左侧：学科分类 -->
        <div class="subject-sidebar">
          <ElMenu @select="handleSubjectChange">
            <ElMenuItem v-for="subject in subjects" :key="subject.id">
              {{ subject.icon }} {{ subject.name }}
            </ElMenuItem>
          </ElMenu>
        </div>

        <!-- 右侧：年级卡片 -->
        <div class="grade-content">
          <div class="grade-grid">
            <ElCard
              v-for="grade in grades"
              :key="grade.grade"
              @click="goToCourses(grade.grade)"
            >
              <h3>{{ grade.grade }}</h3>
              <p>{{ grade.count }} 门课程</p>
            </ElCard>
          </div>
        </div>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { subjectApi, gradeApi } from '@/api/learning'

const router = useRouter()
const subjects = ref([])
const grades = ref([])
const activeSubjectId = ref(null)

const loadSubjects = async () => {
  const res = await subjectApi.getSubjectList()
  subjects.value = res
  if (res.length > 0) {
    activeSubjectId.value = res[0].id
    loadGrades()
  }
}

const loadGrades = async () => {
  const res = await gradeApi.getGradesBySubject(activeSubjectId.value)
  grades.value = res.grades
}

const handleSubjectChange = (subjectId) => {
  activeSubjectId.value = subjectId
  loadGrades()
}

const goToCourses = (grade) => {
  router.push({
    path: '/learning/courses',
    query: { subject_id: activeSubjectId.value, grade }
  })
}

onMounted(() => {
  loadSubjects()
})
</script>
```

### 步骤2：创建 `courses/index.vue` - 课程列表页面

```vue
<template>
  <div class="courses-page">
    <ElCard>
      <!-- 搜索框 -->
      <ElInput v-model="searchKeyword" placeholder="搜索课程" />

      <!-- AI生成的课程 -->
      <div class="section">
        <h2>AI生成的课程</h2>
        <div class="course-grid">
          <ElCard v-for="course in aiCourses" :key="course.id" @click="goToDetail(course.id)">
            <h3>{{ course.title }}</h3>
            <p>{{ course.summary }}</p>
          </ElCard>
        </div>
      </div>

      <!-- 收藏的课程 -->
      <div class="section">
        <h2>收藏的课程</h2>
        <div class="course-grid">
          <ElCard v-for="course in favoriteCourses" :key="course.id" @click="goToDetail(course.id)">
            <h3>{{ course.title }}</h3>
            <p>{{ course.summary }}</p>
          </ElCard>
        </div>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { learningMaterialApi } from '@/api/learning'

const route = useRoute()
const router = useRouter()
const searchKeyword = ref('')
const materials = ref([])

const aiCourses = computed(() => {
  return materials.value.filter(m => 
    m.title.includes(searchKeyword.value)
  )
})

const favoriteCourses = computed(() => {
  // 这里需要后端支持收藏功能
  return []
})

const loadCourses = async () => {
  const res = await learningMaterialApi.getMaterialList({
    subject_id: route.query.subject_id,
    grade: route.query.grade,
    page: 1,
    limit: 100
  })
  materials.value = res.list
}

const goToDetail = (id) => {
  router.push(`/learning/detail/${id}`)
}

onMounted(() => {
  loadCourses()
})
</script>
```

### 步骤3：创建 `detail/index.vue` - 课程详情页面

```vue
<template>
  <div class="detail-page">
    <ElCard v-if="material">
      <h1>{{ material.title }}</h1>
      <div v-html="formatContent(material.content)" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { learningMaterialApi } from '@/api/learning'

const route = useRoute()
const material = ref(null)

const loadDetail = async () => {
  const res = await learningMaterialApi.getMaterialDetail(route.params.id)
  material.value = res
}

const formatContent = (content) => {
  // 格式化内容
  return content
}

onMounted(() => {
  loadDetail()
})
</script>
```

### 步骤4：配置路由

在 `src/router/modules/learning.ts` 中添加：

```typescript
{
  path: 'courses',
  name: 'LearningCourses',
  component: '/learning/courses/index',
  meta: {
    title: '课程列表',
    keepAlive: true
  }
},
{
  path: 'detail/:id',
  name: 'LearningDetail',
  component: '/learning/detail/index',
  meta: {
    title: '课程详情',
    keepAlive: false
  }
}
```

## 🎯 总结

现在的结构是：
1. `/learning/list` - 学科分类 + 年级列表
2. `/learning/courses?subject_id=1&grade=初一` - 课程列表（带搜索）
3. `/learning/detail/:id` - 课程详情

后端已经完成，前端需要按照上面的步骤实现三个页面。
