<!-- 通知组件 -->
<template>
  <div
    class="notice"
    v-show="visible"
    :style="{
      transform: show ? 'scaleY(1)' : 'scaleY(0.9)',
      opacity: show ? 1 : 0
    }"
    @click.stop=""
  >
    <div class="header">
      <span class="text">{{ $t('notice.title') }}</span>
      <span class="btn" @click="markAllRead">{{ $t('notice.btnRead') }}</span>
    </div>

    <ul class="bar">
      <li
        v-for="(item, index) in barList"
        :key="index"
        :class="{ active: barActiveIndex === index }"
        @click="changeBar(index)"
      >
        {{ item.name }} ({{ item.num }})
      </li>
    </ul>

    <div class="content">
      <div class="scroll">
        <!-- 通知 -->
        <ul class="notice-list" v-show="barActiveIndex === 0">
          <li
            v-for="(item, index) in noticeList"
            :key="index"
            @click="handleNoticeClick(item)"
            :class="{ clickable: item.relatedId && item.relatedType }"
          >
            <div
              class="icon"
              :style="{ background: getNoticeStyle(item.type).backgroundColor + '!important' }"
            >
              <i
                class="iconfont-sys"
                :style="{ color: getNoticeStyle(item.type).iconColor + '!important' }"
                v-html="getNoticeStyle(item.type).icon"
              >
              </i>
            </div>
            <div class="text">
              <h4>{{ item.title }}</h4>
              <p>{{ item.time }}</p>
            </div>
          </li>
        </ul>

        <!-- 消息 -->
        <ul class="user-list" v-show="barActiveIndex === 1">
          <li v-for="(item, index) in msgList" :key="index">
            <div class="avatar">
              <img :src="item.avatar" />
            </div>
            <div class="text">
              <h4>{{ item.title }}</h4>
              <p>{{ item.time }}</p>
            </div>
          </li>
        </ul>

        <!-- 待办 -->
        <ul class="base" v-show="barActiveIndex === 2">
          <li v-for="(item, index) in pendingList" :key="index">
            <h4>{{ item.title }}</h4>
            <p>{{ item.time }}</p>
          </li>
        </ul>

        <!-- 空状态 -->
        <div class="empty-tips" v-show="currentTabIsEmpty">
          <i class="iconfont-sys">&#xe8d7;</i>
          <p>{{ $t('notice.text[0]') }}{{ barList[barActiveIndex].name }}</p>
        </div>
      </div>

      <div class="btn-wrapper">
        <ElButton class="view-all" @click="handleViewAll" v-ripple>
          {{ $t('notice.viewAll') }}
        </ElButton>
      </div>
    </div>

    <div style="height: 100px"></div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, watch, type Ref, type ComputedRef } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRouter } from 'vue-router'
  import AppConfig from '@/config'
  import {
    useNotificationStore,
    type NotificationDisplayItem
  } from '@/store/modules/notification'
  import type { NotificationRelatedType } from '@/api/notification'

  // 导入头像图片
  // import avatar1 from '@/assets/img/avatar/avatar1.webp'
  // import avatar2 from '@/assets/img/avatar/avatar2.webp'
  // import avatar3 from '@/assets/img/avatar/avatar3.webp'
  // import avatar4 from '@/assets/img/avatar/avatar4.webp'
  // import avatar5 from '@/assets/img/avatar/avatar5.webp'
  // import avatar6 from '@/assets/img/avatar/avatar6.webp'

  defineOptions({ name: 'ArtNotification' })

  interface NoticeItem {
    /** ID */
    id?: number
    /** 标题 */
    title: string
    /** 时间 */
    time: string
    /** 类型 */
    type: NoticeType
    /** 关联ID */
    relatedId?: number
    /** 关联类型 */
    relatedType?: NotificationRelatedType
  }

  interface MessageItem {
    /** 标题 */
    title: string
    /** 时间 */
    time: string
    /** 头像 */
    avatar: string
  }

  interface PendingItem {
    /** 标题 */
    title: string
    /** 时间 */
    time: string
  }

  interface BarItem {
    /** 名称 */
    name: ComputedRef<string>
    /** 数量 */
    num: number
  }

  interface NoticeStyle {
    /** 图标 */
    icon: string
    /** 图标颜色 */
    iconColor: string
    /** 背景颜色 */
    backgroundColor: string
  }

  type NoticeType = 'email' | 'message' | 'collection' | 'user' | 'notice'

  const { t } = useI18n()
  const router = useRouter()

  const props = defineProps<{
    value: boolean
  }>()

  const emit = defineEmits<{
    (e: 'close'): void
  }>()

  const show = ref(false)
  const visible = ref(false)
  const barActiveIndex = ref(0)

  const useNotificationData = () => {
    // 引入通知 store
    const notificationStore = useNotificationStore()

    // 通知数据 - 从 store 获取
    const noticeList = computed<NoticeItem[]>(() => {
      return notificationStore.notifications.map((n: NotificationDisplayItem) => ({
        id: n.id,
        title: n.title,
        time: n.time,
        type: n.type,
        relatedId: n.relatedId,
        relatedType: n.relatedType
      }))
    })

    // 消息数据
    const msgList = ref<MessageItem[]>([])

    // 待办数据
    const pendingList = ref<PendingItem[]>([])

    // 标签栏数据
    const barList = computed<BarItem[]>(() => [
      {
        name: computed(() => t('notice.bar[0]')),
        num: noticeList.value.length
      },
      {
        name: computed(() => t('notice.bar[1]')),
        num: msgList.value.length
      },
      {
        name: computed(() => t('notice.bar[2]')),
        num: pendingList.value.length
      }
    ])

    // 加载通知数据 / Load notification data
    const loadNotifications = () => {
      notificationStore.fetchNotifications()
    }

    // 标记全部已读 / Mark all as read
    const markAllRead = () => {
      notificationStore.markAllAsRead()
    }

    // 标记单条已读 / Mark single as read
    const markAsRead = (id: number) => {
      notificationStore.markAsRead(id)
    }

    return {
      noticeList,
      msgList,
      pendingList,
      barList,
      loadNotifications,
      markAllRead,
      markAsRead
    }
  }

  // 样式管理
  const useNotificationStyles = () => {
    const noticeStyleMap: Record<NoticeType, NoticeStyle> = {
      email: {
        icon: '&#xe72e;',
        iconColor: 'rgb(var(--art-warning))',
        backgroundColor: 'rgb(var(--art-bg-warning))'
      },
      message: {
        icon: '&#xe747;',
        iconColor: 'rgb(var(--art-success))',
        backgroundColor: 'rgb(var(--art-bg-success))'
      },
      collection: {
        icon: '&#xe714;',
        iconColor: 'rgb(var(--art-danger))',
        backgroundColor: 'rgb(var(--art-bg-danger))'
      },
      user: {
        icon: '&#xe608;',
        iconColor: 'rgb(var(--art-info))',
        backgroundColor: 'rgb(var(--art-bg-info))'
      },
      notice: {
        icon: '&#xe6c2;',
        iconColor: 'rgb(var(--art-primary))',
        backgroundColor: 'rgb(var(--art-bg-primary))'
      }
    }

    const getRandomColor = (): string => {
      const index = Math.floor(Math.random() * AppConfig.systemMainColor.length)
      return AppConfig.systemMainColor[index]
    }

    const getNoticeStyle = (type: NoticeType): NoticeStyle => {
      const defaultStyle: NoticeStyle = {
        icon: '&#xe747;',
        iconColor: '#FFFFFF',
        backgroundColor: getRandomColor()
      }

      return noticeStyleMap[type] || defaultStyle
    }

    return {
      getNoticeStyle
    }
  }

  // 动画管理
  const useNotificationAnimation = () => {
    const showNotice = (open: boolean) => {
      if (open) {
        visible.value = open
        setTimeout(() => {
          show.value = open
        }, 5)
      } else {
        show.value = open
        setTimeout(() => {
          visible.value = open
        }, 350)
      }
    }

    return {
      showNotice
    }
  }

  // 标签页管理
  const useTabManagement = (
    noticeList: Ref<NoticeItem[]> | ComputedRef<NoticeItem[]>,
    msgList: Ref<MessageItem[]>,
    pendingList: Ref<PendingItem[]>,
    businessHandlers: {
      handleNoticeAll: () => void
      handleMsgAll: () => void
      handlePendingAll: () => void
    }
  ) => {
    const changeBar = (index: number) => {
      barActiveIndex.value = index
    }

    // 检查当前标签页是否为空
    const currentTabIsEmpty = computed(() => {
      const tabDataMap = [noticeList.value, msgList.value, pendingList.value]

      const currentData = tabDataMap[barActiveIndex.value]
      return currentData && currentData.length === 0
    })

    const handleViewAll = () => {
      // 查看全部处理器映射
      const viewAllHandlers: Record<number, () => void> = {
        0: businessHandlers.handleNoticeAll,
        1: businessHandlers.handleMsgAll,
        2: businessHandlers.handlePendingAll
      }

      const handler = viewAllHandlers[barActiveIndex.value]
      handler?.()
    }

    return {
      changeBar,
      currentTabIsEmpty,
      handleViewAll
    }
  }

  // 业务逻辑处理
  const useBusinessLogic = () => {
    const handleNoticeAll = () => {
      // 处理查看全部通知
      console.log('查看全部通知')
    }

    const handleMsgAll = () => {
      // 处理查看全部消息
      console.log('查看全部消息')
    }

    const handlePendingAll = () => {
      // 处理查看全部待办
      console.log('查看全部待办')
    }

    /*** 点击通知跳转到对应页面 / Navigate to related page on click ***/
    const handleNoticeClick = (item: NoticeItem) => {
      // 标记为已读
      if (item.id) {
        markAsRead(item.id)
      }

      // 关闭通知面板
      emit('close')

      // 根据关联类型跳转
      if (!item.relatedType || !item.relatedId) {
        return
      }

      // 版本对比需要特殊处理，跳转到对比详情页面
      if (item.relatedType === 'version_compare') {
        router.push({
          path: '/designer/designer-assistant/version-compare/VersionDiff',
          query: { compareId: item.relatedId }
        })
        return
      }

      const routeMap: Record<NotificationRelatedType, string> = {
        version_compare: '/designer/designer-assistant/version-compare/VersionDiff',
        render: '/designer/designer-assistant/cad-viewer/CadUpload',
        cad_generation: '/designer/designer-assistant/cad-generation/CadGenerationList',
        document_analysis: '/designer/designer-assistant/document/DocumentUpload',
        design_suggestion: '/designer/designer-assistant/suggestion/SuggestionList',
        general: ''
      }

      const basePath = routeMap[item.relatedType]
      if (basePath) {
        router.push({
          path: basePath,
          query: { id: item.relatedId }
        })
      }
    }

    return {
      handleNoticeAll,
      handleMsgAll,
      handlePendingAll,
      handleNoticeClick
    }
  }

  // 组合所有逻辑
  const { noticeList, msgList, pendingList, barList, loadNotifications, markAllRead, markAsRead } =
    useNotificationData()
  const { getNoticeStyle } = useNotificationStyles()
  const { showNotice } = useNotificationAnimation()
  const { handleNoticeAll, handleMsgAll, handlePendingAll, handleNoticeClick } = useBusinessLogic()
  const { changeBar, currentTabIsEmpty, handleViewAll } = useTabManagement(
    noticeList,
    msgList,
    pendingList,
    { handleNoticeAll, handleMsgAll, handlePendingAll }
  )

  // 监听属性变化
  watch(
    () => props.value,
    (newValue) => {
      showNotice(newValue)
      // 打开通知面板时加载数据 / Load data when opening notification panel
      if (newValue) {
        loadNotifications()
      }
    }
  )
</script>

<style lang="scss" scoped>
  @use './style';
</style>
