import { AppRouteRecord } from '@/types/router'

/**
 * 工作流路由配置
 * 包含流程管理、审批工作台等功能模块
 */
export const workflowRoutes: AppRouteRecord = {
  path: '/workflow',
  name: 'Workflow',
  component: '/index/index',
  meta: {
    title: '工作流',
    icon: '&#xe6a1;',
    keepAlive: true
  },
  children: [
    // 审批工作台
    {
      path: 'workbench',
      name: 'WorkflowWorkbench',
      component: '/index/index',
      meta: {
        title: '审批工作台',
        keepAlive: true
      },
      children: [
        {
          path: 'todo',
          name: 'WorkflowTodo',
          component: '/workflow/workbench/todo',
          meta: {
            title: '我的待办',
            keepAlive: true
          }
        },
        {
          path: 'done',
          name: 'WorkflowDone',
          component: '/workflow/workbench/done',
          meta: {
            title: '我的已办',
            keepAlive: true
          }
        },
        {
          path: 'initiated',
          name: 'WorkflowInitiated',
          component: '/workflow/workbench/initiated',
          meta: {
            title: '我发起的',
            keepAlive: true
          }
        },
        {
          path: 'detail/:id',
          name: 'WorkflowDetail',
          component: '/workflow/workbench/detail',
          meta: {
            title: '审批详情',
            keepAlive: false,
            hidden: true
          }
        }
      ]
    },
    // 发起申请
    {
      path: 'apply',
      name: 'WorkflowApply',
      component: '/workflow/apply/index',
      meta: {
        title: '发起申请',
        keepAlive: true
      }
    },
    {
      path: 'apply/form/:defId',
      name: 'WorkflowApplyForm',
      component: '/workflow/apply/form',
      meta: {
        title: '填写申请',
        keepAlive: false,
        hidden: true
      }
    },
    // 流程管理
    {
      path: 'management',
      name: 'WorkflowManagement',
      component: '/index/index',
      meta: {
        title: '流程管理',
        keepAlive: true
      },
      children: [
        {
          path: 'process-def',
          name: 'ProcessDefList',
          component: '/workflow/process-def/index',
          meta: {
            title: '流程定义',
            keepAlive: true
          }
        },
        {
          path: 'process-def/edit/:id?',
          name: 'ProcessDefEdit',
          component: '/workflow/process-def/edit',
          meta: {
            title: '流程设计',
            keepAlive: false,
            hidden: true
          }
        },
        {
          path: 'form-template',
          name: 'FormTemplateList',
          component: '/workflow/form-template/index',
          meta: {
            title: '表单模板',
            keepAlive: true
          }
        },
        {
          path: 'form-template/edit/:id?',
          name: 'FormTemplateEdit',
          component: '/workflow/form-template/edit',
          meta: {
            title: '表单设计',
            keepAlive: false,
            hidden: true
          }
        }
      ]
    }
  ]
}

export default workflowRoutes
