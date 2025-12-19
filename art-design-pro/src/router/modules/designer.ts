import { AppRouteRecord } from '@/types/router'

/***
 * Designer Assistant Module Routes
 * 工装设计师AI辅助系统路由配置
 * Requirements: 4.2, 11.1, 11.2
 ***/
export const designerRoutes: AppRouteRecord = {
  path: '/designer',
  name: 'Designer',
  component: '/index/index',
  meta: {
    title: '设计师助手',
    icon: '&#xe7b0;',
    keepAlive: true
  },
  children: [
    {
      path: 'project',
      name: 'DesignerProject',
      component: '/designer-assistant/project/ProjectList',
      meta: {
        title: '项目管理',
        keepAlive: true
      }
    },
    {
      path: 'project/create',
      name: 'DesignerProjectCreate',
      component: '/designer-assistant/project/ProjectCreate',
      meta: {
        title: '  ',
        keepAlive: false,
        hidden: true
      }
    },
    {
      path: 'project/:id',
      name: 'DesignerProjectDetail',
      component: '/designer-assistant/project/ProjectDetail',
      meta: {
        title: '项目详情',
        keepAlive: false,
        hidden: true
      }
    },
    {
      path: 'project/edit/:id',
      name: 'DesignerProjectEdit',
      component: '/designer-assistant/project/ProjectCreate',
      meta: {
        title: '编辑项目',
        keepAlive: false,
        hidden: true
      }
    },
    {
      path: 'document',
      name: 'DesignerDocument',
      component: '/designer-assistant/document/DocumentUpload',
      meta: {
        title: '文档分析',
        keepAlive: true
      }
    },
    {
      path: 'compare',
      name: 'DesignerCompare',
      component: '/designer-assistant/design-compare/CompareUpload',
      meta: {
        title: '设计比对',
        keepAlive: true
      }
    },
    {
      path: 'cad',
      name: 'DesignerCad',
      component: '/designer-assistant/cad-viewer/CadUpload',
      meta: {
        title: 'CAD预览',
        keepAlive: true
      }
    },
    {
      path: 'material',
      name: 'DesignerMaterial',
      component: '/designer-assistant/material/MaterialList',
      meta: {
        title: '材料库',
        keepAlive: true
      }
    },
    {
      path: 'cost',
      name: 'DesignerCost',
      component: '/designer-assistant/cost/CostConfig',
      meta: {
        title: '成本估算',
        keepAlive: true
      }
    },
    {
      path: 'chat',
      name: 'DesignerChat',
      component: '/designer-assistant/chat/DesignerChat',
      meta: {
        title: 'AI对话',
        keepAlive: true
      }
    },
    {
      path: 'compliance',
      name: 'DesignerCompliance',
      component: '/designer-assistant/compliance/ComplianceCheck',
      meta: {
        title: '合规检查',
        keepAlive: true
      }
    },
    {
      path: 'cad-generation',
      name: 'DesignerCadGeneration',
      component: '/designer-assistant/cad-generation/index',
      meta: {
        title: 'AI生成CAD',
        keepAlive: true
      }
    },
    {
      path: 'version-compare',
      name: 'DesignerVersionCompare',
      component: '/designer-assistant/version-compare/VersionList',
      meta: {
        title: '设计版本',
        keepAlive: true
      }
    },
    {
      path: 'version-compare/diff',
      name: 'DesignerVersionDiff',
      component: '/designer-assistant/version-compare/VersionDiff',
      meta: {
        title: '版本差异',
        keepAlive: false,
        hidden: true
      }
    },
    {
      path: 'construction-annotation',
      name: 'DesignerConstructionAnnotation',
      component: '/designer-assistant/construction-annotation/index',
      meta: {
        title: '施工图标注',
        keepAlive: true
      }
    },
    {
      path: 'construction-annotation/edit',
      name: 'DesignerConstructionAnnotationEdit',
      component: '/designer-assistant/construction-annotation/EditPage',
      meta: {
        title: '编辑标注',
        keepAlive: false,
        hidden: true
      }
    }
  ]
}

export default designerRoutes
