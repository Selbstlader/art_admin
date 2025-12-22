<template>
  <div class="material-list art-full-height">
    <!-- 搜索筛选区域 / Search and filter area -->
    <ArtSearchBar
      v-model="filterForm"
      :items="searchItems"
      @search="handleSearch"
      @reset="handleReset"
    />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 / Table Header -->
      <ArtTableHeader :loading="loading" @refresh="loadMaterialList">
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="handleCreate" v-ripple>新增材料</ElButton>
            <ElButton :disabled="selectedIds.length === 0" @click="handleBatchDelete" v-ripple>
              批量删除
            </ElButton>
            <ElButton @click="handleImport" v-ripple>批量导入</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 材料列表表格 / Material list table -->
      <ArtTable
        :loading="loading"
        :data="materialList"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>

    <!-- 批量导入弹窗 / Batch import dialog -->
    <el-dialog v-model="importDialogVisible" title="批量导入材料" width="900px" destroy-on-close>
      <div class="import-dialog">
        <!-- 步骤1: 上传文件 / Step 1: Upload file -->
        <div v-if="importStep === 1" class="import-step">
          <div class="import-tips">
            <h4>导入说明：</h4>
            <ul>
              <li>支持 Excel 文件格式（.xlsx, .xls）</li>
              <li>请按照模板格式填写数据，必填字段：材料名称、分类、单位、单价</li>
              <li>单次最多导入 500 条数据</li>
            </ul>
            <el-button type="primary" link @click="downloadTemplate">
              <el-icon><Download /></el-icon>
              下载导入模板
            </el-button>
          </div>
          <el-upload
            ref="uploadRef"
            class="import-upload"
            drag
            :auto-upload="false"
            :limit="1"
            accept=".xlsx,.xls"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
          >
            <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
            <template #tip>
              <div class="el-upload__tip">只能上传 xlsx/xls 文件</div>
            </template>
          </el-upload>
        </div>

        <!-- 步骤2: 预览数据 / Step 2: Preview data -->
        <div v-if="importStep === 2" class="import-step">
          <div class="preview-header">
            <span>共解析 {{ importPreviewData.length }} 条数据</span>
            <el-tag v-if="importErrors.length > 0" type="danger">
              {{ importErrors.length }} 条数据有错误
            </el-tag>
          </div>
          <el-table :data="importPreviewData" max-height="400" size="small" border>
            <el-table-column type="index" label="序号" width="60" />
            <el-table-column prop="name" label="材料名称" min-width="120" show-overflow-tooltip>
              <template #default="{ row, $index }">
                <span :class="{ 'error-cell': hasFieldError($index, 'name') }">
                  {{ row.name || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="category" label="分类" width="100">
              <template #default="{ row, $index }">
                <span :class="{ 'error-cell': hasFieldError($index, 'category') }">
                  {{ row.category || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="specification" label="规格" width="100" show-overflow-tooltip />
            <el-table-column prop="unit" label="单位" width="60">
              <template #default="{ row, $index }">
                <span :class="{ 'error-cell': hasFieldError($index, 'unit') }">
                  {{ row.unit || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="unitPrice" label="单价" width="80">
              <template #default="{ row, $index }">
                <span :class="{ 'error-cell': hasFieldError($index, 'unitPrice') }">
                  {{ row.unitPrice ?? '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="brand" label="品牌" width="80" show-overflow-tooltip />
            <el-table-column prop="supplier" label="供应商" width="100" show-overflow-tooltip />
            <el-table-column label="状态" width="70">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '启用' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="错误信息" min-width="150">
              <template #default="{ $index }">
                <span v-if="getRowErrors($index)" class="error-text">
                  {{ getRowErrors($index) }}
                </span>
                <span v-else class="success-text">✓ 正常</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 步骤3: 导入结果 / Step 3: Import result -->
        <div v-if="importStep === 3" class="import-step import-result">
          <el-result
            :icon="importResult.failCount === 0 ? 'success' : 'warning'"
            :title="importResult.failCount === 0 ? '导入成功' : '部分导入成功'"
          >
            <template #sub-title>
              <div class="result-stats">
                <p
                  >成功导入：<span class="success-count">{{ importResult.successCount }}</span>
                  条</p
                >
                <p v-if="importResult.failCount > 0">
                  导入失败：<span class="fail-count">{{ importResult.failCount }}</span> 条
                </p>
              </div>
            </template>
            <template #extra>
              <div v-if="importResult.failReasons.length > 0" class="fail-reasons">
                <h4>失败原因：</h4>
                <ul>
                  <li v-for="(reason, idx) in importResult.failReasons.slice(0, 10)" :key="idx">
                    {{ reason }}
                  </li>
                  <li v-if="importResult.failReasons.length > 10">
                    ... 还有 {{ importResult.failReasons.length - 10 }} 条错误
                  </li>
                </ul>
              </div>
            </template>
          </el-result>
        </div>
      </div>
      <template #footer>
        <template v-if="importStep === 1">
          <el-button @click="importDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            :disabled="!importFile"
            :loading="parseLoading"
            @click="parseImportFile"
          >
            解析文件
          </el-button>
        </template>
        <template v-else-if="importStep === 2">
          <el-button @click="importStep = 1">上一步</el-button>
          <el-button
            type="primary"
            :disabled="importErrors.length > 0 || importPreviewData.length === 0"
            :loading="importLoading"
            @click="submitImport"
          >
            确认导入 ({{ importPreviewData.length }} 条)
          </el-button>
        </template>
        <template v-else>
          <el-button type="primary" @click="closeImportDialog">完成</el-button>
        </template>
      </template>
    </el-dialog>

    <!-- 材料详情/编辑弹窗 / Material detail/edit dialog -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        :disabled="dialogMode === 'view'"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="材料名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入材料名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-autocomplete
                v-model="formData.category"
                :fetch-suggestions="queryCategorySearch"
                placeholder="输入或选择分类"
                style="width: 100%"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="规格" prop="specification">
              <el-input v-model="formData.specification" placeholder="请输入规格" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="单位" prop="unit">
              <el-autocomplete
                v-model="formData.unit"
                :fetch-suggestions="queryUnitSearch"
                placeholder="输入或选择单位"
                style="width: 100%"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="单价" prop="unitPrice">
              <el-input-number
                v-model="formData.unitPrice"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="品牌" prop="brand">
              <el-input v-model="formData.brand" placeholder="请输入品牌" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="供应商" prop="supplier">
              <el-input v-model="formData.supplier" placeholder="请输入供应商" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio value="active">启用</el-radio>
                <el-radio value="inactive">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="材料图片" prop="imageUrl">
          <div class="image-upload-wrapper">
            <el-upload
              class="image-uploader"
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleImageChange"
              accept="image/*"
              :disabled="dialogMode === 'view'"
            >
              <img v-if="formData.imageUrl" :src="formData.imageUrl" class="uploaded-image" />
              <div v-else class="upload-placeholder">
                <el-icon class="upload-icon"><Plus /></el-icon>
                <span>点击上传</span>
              </div>
            </el-upload>
            <el-button
              v-if="formData.imageUrl && dialogMode !== 'view'"
              type="danger"
              link
              size="small"
              class="remove-btn"
              @click="handleRemoveImage"
            >
              移除图片
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="适用场景" prop="applicableScenes">
          <el-input
            v-model="applicableScenesInput"
            placeholder="输入适用场景，多个用逗号分隔"
            style="width: 100%"
          />
          <div class="scene-hint">常用：办公空间、商业空间、工业厂房、酒店、餐饮、医疗、教育</div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入材料描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button
          v-if="dialogMode !== 'view'"
          type="primary"
          :loading="submitLoading"
          @click="handleSubmit"
        >
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted, h } from 'vue'
  import {
    ElMessage,
    ElMessageBox,
    ElImage,
    ElTag,
    type FormInstance,
    type FormRules,
    type UploadInstance,
    type UploadFile,
    type UploadRawFile
  } from 'element-plus'
  import { Plus, Download, UploadFilled } from '@element-plus/icons-vue'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import * as XLSX from 'xlsx'
  import {
    getMaterialList,
    getMaterialCategories,
    getMaterialBrands,
    createMaterial,
    updateMaterial,
    deleteMaterial,
    batchDeleteMaterials,
    batchImportMaterials,
    type MaterialResponse,
    type MaterialCategoryResponse,
    type MaterialBrandResponse,
    type CreateMaterialRequest
  } from '@/api/designer-material'
  import { uploadFile } from '@/api/file'
  import type { ColumnOption } from '@/types/component'

  /*** 筛选表单 / Filter form ***/
  const filterForm = ref({
    keyword: '',
    category: '',
    brand: '',
    minPrice: undefined as number | undefined,
    maxPrice: undefined as number | undefined,
    status: ''
  })

  /*** 分类和品牌数据 / Categories and brands data ***/
  const categories = ref<MaterialCategoryResponse[]>([])
  const brands = ref<MaterialBrandResponse[]>([])

  /*** 搜索配置 / Search config ***/
  const searchItems = computed(() => [
    {
      label: '关键字',
      key: 'keyword',
      type: 'input',
      placeholder: '搜索材料名称/描述',
      clearable: true
    },
    {
      label: '分类',
      key: 'category',
      type: 'select',
      props: {
        placeholder: '全部分类',
        options: categories.value.map((cat) => ({
          label: `${cat.category} (${cat.count})`,
          value: cat.category
        }))
      }
    },
    {
      label: '品牌',
      key: 'brand',
      type: 'select',
      props: {
        placeholder: '全部品牌',
        options: brands.value.map((b) => ({
          label: `${b.brand} (${b.count})`,
          value: b.brand
        }))
      }
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: {
        placeholder: '全部状态',
        options: [
          { label: '启用', value: 'active' },
          { label: '停用', value: 'inactive' }
        ]
      }
    }
  ])

  /*** 分页 / Pagination ***/
  const pagination = reactive({
    current: 1,
    size: 10,
    total: 0
  })

  /*** 数据 / Data ***/
  const loading = ref(false)
  const materialList = ref<MaterialResponse[]>([])
  const selectedIds = ref<number[]>([])

  /*** 表格列配置 / Table columns config ***/
  const columns = computed<ColumnOption[]>(() => [
    { type: 'selection', width: 50 },
    {
      prop: 'imageUrl',
      label: '图片',
      width: 80,
      formatter: (row: MaterialResponse) =>
        row.imageUrl
          ? h(ElImage, {
              src: row.imageUrl,
              previewSrcList: [row.imageUrl],
              fit: 'cover',
              style: 'width: 50px; height: 50px; border-radius: 4px'
            })
          : h('div', { class: 'no-image' }, '暂无')
    },
    { prop: 'name', label: '材料名称', minWidth: 150, showOverflowTooltip: true },
    { prop: 'category', label: '分类', width: 100 },
    { prop: 'brand', label: '品牌', width: 100, showOverflowTooltip: true },
    { prop: 'specification', label: '规格', width: 120, showOverflowTooltip: true },
    {
      prop: 'unitPrice',
      label: '单价',
      width: 120,
      formatter: (row: MaterialResponse) =>
        h('span', { class: 'price' }, `¥${row.unitPrice.toFixed(2)}/${row.unit}`)
    },
    { prop: 'supplier', label: '供应商', width: 120, showOverflowTooltip: true },
    {
      prop: 'status',
      label: '状态',
      width: 80,
      formatter: (row: MaterialResponse) =>
        h(ElTag, { type: row.status === 'active' ? 'success' : 'info', size: 'small' }, () =>
          row.status === 'active' ? '启用' : '停用'
        )
    },
    {
      prop: 'operation',
      label: '操作',
      width: 180,
      fixed: 'right',
      formatter: (row: MaterialResponse) =>
        h('div', { class: 'flex gap-1' }, [
          h(ArtButtonTable, { type: 'view', onClick: () => handleView(row) }),
          h(ArtButtonTable, { type: 'edit', onClick: () => handleEdit(row) }),
          h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
        ])
    }
  ])

  /*** 弹窗 / Dialog ***/
  const dialogVisible = ref(false)
  const dialogMode = ref<'create' | 'edit' | 'view'>('create')
  const dialogTitle = ref('新增材料')
  const submitLoading = ref(false)
  const formRef = ref<FormInstance>()

  const formData = reactive<CreateMaterialRequest & { id?: number }>({
    name: '',
    category: '',
    specification: '',
    unit: 'm²',
    unitPrice: 0,
    brand: '',
    supplier: '',
    description: '',
    imageUrl: '',
    applicableScenes: [],
    status: 'active'
  })

  const formRules: FormRules = {
    name: [{ required: true, message: '请输入材料名称', trigger: 'blur' }],
    category: [{ required: true, message: '请输入分类', trigger: 'blur' }],
    unit: [{ required: true, message: '请输入单位', trigger: 'blur' }],
    unitPrice: [{ required: true, message: '请输入单价', trigger: 'blur' }]
  }

  // 适用场景输入框值 / Applicable scenes input value
  const applicableScenesInput = ref('')

  /*** 批量导入相关 / Batch import related ***/
  const importDialogVisible = ref(false)
  const importStep = ref(1) // 1: 上传, 2: 预览, 3: 结果
  const importFile = ref<UploadRawFile | null>(null)
  const uploadRef = ref<UploadInstance>()
  const parseLoading = ref(false)
  const importLoading = ref(false)
  const importPreviewData = ref<CreateMaterialRequest[]>([])
  const importErrors = ref<{ row: number; field: string; message: string }[]>([])
  const importResult = ref({
    successCount: 0,
    failCount: 0,
    failReasons: [] as string[]
  })

  // 预设单位列表 / Preset unit list
  const presetUnits = ['m²', 'm', '个', 'kg', '套', '组', '块', '张', '卷', '桶', '件', '根']

  // 图片上传loading / Image upload loading
  const imageUploading = ref(false)

  // 处理图片选择 / Handle image change
  const handleImageChange = async (file: UploadFile) => {
    if (!file.raw) return

    // 验证文件类型 / Validate file type
    const isImage = file.raw.type.startsWith('image/')
    if (!isImage) {
      ElMessage.error('只能上传图片文件')
      return
    }

    // 验证文件大小（最大5MB）/ Validate file size (max 5MB)
    const isLt5M = file.raw.size / 1024 / 1024 < 5
    if (!isLt5M) {
      ElMessage.error('图片大小不能超过5MB')
      return
    }

    imageUploading.value = true
    try {
      const res: any = await uploadFile(file.raw, 'image')
      // API返回结构: { code, msg, data: { url, ... } }
      const url = res?.url || res?.data?.url
      if (url) {
        formData.imageUrl = url
        ElMessage.success('图片上传成功')
      } else {
        ElMessage.error('上传成功但未获取到图片地址')
      }
    } catch (error) {
      console.error('图片上传失败:', error)
      ElMessage.error('图片上传失败')
    } finally {
      imageUploading.value = false
    }
  }

  // 移除图片 / Remove image
  const handleRemoveImage = () => {
    formData.imageUrl = ''
  }

  // 分类自动完成 / Category autocomplete
  const queryCategorySearch = (queryString: string, cb: (results: { value: string }[]) => void) => {
    const results = categories.value
      .map((c) => ({ value: c.category }))
      .filter(
        (item) => !queryString || item.value.toLowerCase().includes(queryString.toLowerCase())
      )
    cb(results)
  }

  // 单位自动完成 / Unit autocomplete
  const queryUnitSearch = (queryString: string, cb: (results: { value: string }[]) => void) => {
    const results = presetUnits
      .map((u) => ({ value: u }))
      .filter(
        (item) => !queryString || item.value.toLowerCase().includes(queryString.toLowerCase())
      )
    cb(results)
  }

  /*** 方法 / Methods ***/

  // 加载材料列表 / Load material list
  const loadMaterialList = async () => {
    loading.value = true
    try {
      const res: any = await getMaterialList({
        current: pagination.current,
        size: pagination.size,
        keyword: filterForm.value.keyword || undefined,
        category: filterForm.value.category || undefined,
        brand: filterForm.value.brand || undefined,
        minPrice: filterForm.value.minPrice,
        maxPrice: filterForm.value.maxPrice,
        status: filterForm.value.status || undefined
      })
      materialList.value = res.data?.records || []
      pagination.total = res.data?.total || 0
    } catch (error) {
      console.error('加载材料列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 加载分类和品牌 / Load categories and brands
  const loadFilters = async () => {
    try {
      const [catRes, brandRes]: any[] = await Promise.all([
        getMaterialCategories(),
        getMaterialBrands()
      ])
      categories.value = catRes.data || []
      brands.value = brandRes.data || []
    } catch (error) {
      console.error('加载筛选数据失败:', error)
    }
  }

  // 搜索 / Search
  const handleSearch = () => {
    pagination.current = 1
    loadMaterialList()
  }

  // 重置 / Reset
  const handleReset = () => {
    filterForm.value = {
      keyword: '',
      category: '',
      brand: '',
      minPrice: undefined,
      maxPrice: undefined,
      status: ''
    }
    handleSearch()
  }

  // 分页变化 / Pagination change
  const handleSizeChange = (row: any) => {
    pagination.size = row
    loadMaterialList()
  }

  const handleCurrentChange = (row: any) => {
    pagination.current = row
    loadMaterialList()
  }

  // 选择变化 / Selection change
  const handleSelectionChange = (selection: MaterialResponse[]) => {
    selectedIds.value = selection.map((item) => item.id)
  }

  // 新增 / Create
  const handleCreate = () => {
    dialogMode.value = 'create'
    dialogTitle.value = '新增材料'
    resetForm()
    dialogVisible.value = true
  }

  // 查看 / View
  const handleView = (row: MaterialResponse) => {
    dialogMode.value = 'view'
    dialogTitle.value = '材料详情'
    Object.assign(formData, row)
    applicableScenesInput.value = (row.applicableScenes || []).join('、')
    dialogVisible.value = true
  }

  // 编辑 / Edit
  const handleEdit = (row: MaterialResponse) => {
    dialogMode.value = 'edit'
    dialogTitle.value = '编辑材料'
    Object.assign(formData, row)
    applicableScenesInput.value = (row.applicableScenes || []).join('、')
    dialogVisible.value = true
  }

  // 删除 / Delete
  const handleDelete = async (row: MaterialResponse) => {
    try {
      await ElMessageBox.confirm(`确定要删除材料"${row.name}"吗？`, '提示', {
        type: 'warning'
      })
      await deleteMaterial(row.id)
      ElMessage.success('删除成功')
      loadMaterialList()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
      }
    }
  }

  // 批量删除 / Batch delete
  const handleBatchDelete = async () => {
    if (selectedIds.value.length === 0) return
    try {
      await ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}个材料吗？`, '提示', {
        type: 'warning'
      })
      await batchDeleteMaterials(selectedIds.value)
      ElMessage.success('批量删除成功')
      loadMaterialList()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('批量删除失败:', error)
      }
    }
  }

  // 导入 / Import
  const handleImport = () => {
    importStep.value = 1
    importFile.value = null
    importPreviewData.value = []
    importErrors.value = []
    importResult.value = { successCount: 0, failCount: 0, failReasons: [] }
    importDialogVisible.value = true
  }

  // 下载导入模板 / Download import template
  const downloadTemplate = () => {
    const templateData = [
      {
        材料名称: '示例材料',
        分类: '地面材料',
        规格: '800x800mm',
        单位: 'm²',
        单价: 150,
        品牌: '示例品牌',
        供应商: '示例供应商',
        适用场景: '办公空间,商业空间',
        描述: '材料描述信息',
        状态: '启用'
      }
    ]
    const ws = XLSX.utils.json_to_sheet(templateData)
    // 设置列宽 / Set column width
    ws['!cols'] = [
      { wch: 15 },
      { wch: 12 },
      { wch: 15 },
      { wch: 8 },
      { wch: 10 },
      { wch: 12 },
      { wch: 15 },
      { wch: 25 },
      { wch: 30 },
      { wch: 8 }
    ]
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, '材料导入模板')
    XLSX.writeFile(wb, '材料导入模板.xlsx')
  }

  // 文件变化 / File change
  const handleFileChange = (file: UploadFile) => {
    importFile.value = file.raw || null
  }

  // 文件超出限制 / File exceed
  const handleExceed = () => {
    ElMessage.warning('只能上传一个文件，请先删除已选文件')
  }

  // 解析导入文件 / Parse import file
  const parseImportFile = async () => {
    if (!importFile.value) return

    parseLoading.value = true
    try {
      const data = await readExcelFile(importFile.value)
      const { materials, errors } = validateImportData(data)
      importPreviewData.value = materials
      importErrors.value = errors
      importStep.value = 2
    } catch (error) {
      ElMessage.error('文件解析失败，请检查文件格式')
      console.error('解析文件失败:', error)
    } finally {
      parseLoading.value = false
    }
  }

  // 读取Excel文件 / Read Excel file
  const readExcelFile = (file: File): Promise<Record<string, any>[]> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        try {
          const data = new Uint8Array(e.target?.result as ArrayBuffer)
          const workbook = XLSX.read(data, { type: 'array' })
          const sheetName = workbook.SheetNames[0]
          const worksheet = workbook.Sheets[sheetName]
          const jsonData = XLSX.utils.sheet_to_json(worksheet)
          resolve(jsonData as Record<string, any>[])
        } catch (err) {
          reject(err)
        }
      }
      reader.onerror = reject
      reader.readAsArrayBuffer(file)
    })
  }

  // 验证导入数据 / Validate import data
  const validateImportData = (
    data: Record<string, any>[]
  ): {
    materials: CreateMaterialRequest[]
    errors: { row: number; field: string; message: string }[]
  } => {
    const materials: CreateMaterialRequest[] = []
    const errors: { row: number; field: string; message: string }[] = []

    // 字段映射 / Field mapping
    const fieldMap: Record<string, string> = {
      材料名称: 'name',
      分类: 'category',
      规格: 'specification',
      单位: 'unit',
      单价: 'unitPrice',
      品牌: 'brand',
      供应商: 'supplier',
      适用场景: 'applicableScenes',
      描述: 'description',
      状态: 'status'
    }

    data.forEach((row, index) => {
      const material: CreateMaterialRequest = {
        name: '',
        category: '',
        unit: 'm²',
        unitPrice: 0,
        status: 'active'
      }

      // 映射字段 / Map fields
      Object.keys(row).forEach((key) => {
        const field = fieldMap[key]
        if (field) {
          let value = row[key]
          if (field === 'applicableScenes') {
            // 解析适用场景 / Parse applicable scenes
            material.applicableScenes = String(value || '')
              .split(/[,，、]/)
              .map((s) => s.trim())
              .filter((s) => s)
          } else if (field === 'unitPrice') {
            material.unitPrice = parseFloat(value) || 0
          } else if (field === 'status') {
            material.status = value === '停用' ? 'inactive' : 'active'
          } else {
            ;(material as any)[field] = String(value || '').trim()
          }
        }
      })

      // 验证必填字段 / Validate required fields
      if (!material.name) {
        errors.push({ row: index, field: 'name', message: '材料名称不能为空' })
      }
      if (!material.category) {
        errors.push({ row: index, field: 'category', message: '分类不能为空' })
      }
      if (!material.unit) {
        errors.push({ row: index, field: 'unit', message: '单位不能为空' })
      }
      if (material.unitPrice < 0) {
        errors.push({ row: index, field: 'unitPrice', message: '单价不能为负数' })
      }

      materials.push(material)
    })

    return { materials, errors }
  }

  // 检查字段是否有错误 / Check if field has error
  const hasFieldError = (rowIndex: number, field: string): boolean => {
    return importErrors.value.some((e) => e.row === rowIndex && e.field === field)
  }

  // 获取行错误信息 / Get row errors
  const getRowErrors = (rowIndex: number): string => {
    const rowErrors = importErrors.value.filter((e) => e.row === rowIndex)
    return rowErrors.map((e) => e.message).join('; ')
  }

  // 提交导入 / Submit import
  const submitImport = async () => {
    if (importPreviewData.value.length === 0) return

    importLoading.value = true
    try {
      const res: any = await batchImportMaterials(importPreviewData.value)
      importResult.value = {
        successCount: res.data?.successCount || 0,
        failCount: res.data?.failCount || 0,
        failReasons: res.data?.failReasons || []
      }
      importStep.value = 3
      // 刷新列表 / Refresh list
      loadMaterialList()
      loadFilters()
    } catch (error) {
      ElMessage.error('导入失败')
      console.error('导入失败:', error)
    } finally {
      importLoading.value = false
    }
  }

  // 关闭导入弹窗 / Close import dialog
  const closeImportDialog = () => {
    importDialogVisible.value = false
  }

  // 重置表单 / Reset form
  const resetForm = () => {
    formData.id = undefined
    formData.name = ''
    formData.category = ''
    formData.specification = ''
    formData.unit = 'm²'
    formData.unitPrice = 0
    formData.brand = ''
    formData.supplier = ''
    formData.description = ''
    formData.imageUrl = ''
    formData.applicableScenes = []
    formData.status = 'active'
    applicableScenesInput.value = ''
  }

  // 提交 / Submit
  const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate()

    // 解析适用场景输入 / Parse applicable scenes input
    const scenes = applicableScenesInput.value
      .split(/[,，、]/)
      .map((s) => s.trim())
      .filter((s) => s)
    formData.applicableScenes = scenes

    submitLoading.value = true
    try {
      if (dialogMode.value === 'create') {
        await createMaterial(formData)
        ElMessage.success('创建成功')
      } else {
        await updateMaterial({ ...formData, id: formData.id! })
        ElMessage.success('更新成功')
      }
      dialogVisible.value = false
      loadMaterialList()
      loadFilters()
    } catch (error) {
      console.error('提交失败:', error)
    } finally {
      submitLoading.value = false
    }
  }

  /*** 生命周期 / Lifecycle ***/
  onMounted(() => {
    loadMaterialList()
    loadFilters()
  })
</script>

<style scoped lang="scss">
  .material-list {
    .no-image {
      width: 50px;
      height: 50px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f7fa;
      border-radius: 4px;
      color: #909399;
      font-size: 12px;
    }

    .price {
      color: #f56c6c;
      font-weight: 500;
    }

    .scene-hint {
      margin-top: 4px;
      font-size: 12px;
      color: #909399;
    }

    /*** 图片上传样式 / Image upload styles ***/
    .image-upload-wrapper {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .image-uploader {
        :deep(.el-upload) {
          border: 1px dashed #d9d9d9;
          border-radius: 8px;
          cursor: pointer;
          overflow: hidden;
          transition: border-color 0.2s;

          &:hover {
            border-color: #409eff;
          }
        }
      }

      .uploaded-image {
        width: 120px;
        height: 120px;
        object-fit: cover;
        display: block;
      }

      .upload-placeholder {
        width: 120px;
        height: 120px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: #8c939d;
        background: #fafafa;

        .upload-icon {
          font-size: 28px;
          margin-bottom: 8px;
        }

        span {
          font-size: 12px;
        }
      }

      .remove-btn {
        align-self: flex-start;
      }
    }
  }

  /*** 批量导入弹窗样式 / Batch import dialog styles ***/
  .import-dialog {
    .import-step {
      min-height: 300px;
    }

    .import-tips {
      margin-bottom: 20px;
      padding: 16px;
      background: #f5f7fa;
      border-radius: 8px;

      h4 {
        margin: 0 0 12px;
        font-size: 14px;
        color: #303133;
      }

      ul {
        margin: 0 0 12px;
        padding-left: 20px;
        color: #606266;
        font-size: 13px;
        line-height: 1.8;
      }
    }

    .import-upload {
      :deep(.el-upload-dragger) {
        padding: 40px 20px;
      }
    }

    .preview-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;
      font-size: 14px;
      color: #606266;
    }

    .error-cell {
      color: #f56c6c;
      font-weight: 500;
    }

    .error-text {
      color: #f56c6c;
      font-size: 12px;
    }

    .success-text {
      color: #67c23a;
      font-size: 12px;
    }

    .import-result {
      display: flex;
      align-items: center;
      justify-content: center;

      .result-stats {
        font-size: 14px;
        color: #606266;

        .success-count {
          color: #67c23a;
          font-weight: 600;
          font-size: 18px;
        }

        .fail-count {
          color: #f56c6c;
          font-weight: 600;
          font-size: 18px;
        }
      }

      .fail-reasons {
        max-height: 200px;
        overflow-y: auto;
        text-align: left;
        padding: 12px;
        background: #fef0f0;
        border-radius: 4px;

        h4 {
          margin: 0 0 8px;
          font-size: 13px;
          color: #f56c6c;
        }

        ul {
          margin: 0;
          padding-left: 16px;
          font-size: 12px;
          color: #f56c6c;
          line-height: 1.6;
        }
      }
    }
  }
</style>
