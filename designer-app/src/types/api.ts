/**
 * API 相关类型定义
 * API related type definitions
 */

/*** 统一响应格式 - Unified response format ***/
export interface ApiResponse<T = unknown> {
  code: number;
  msg?: string;
  message?: string;
  data?: T;
}

/*** 分页响应格式 - Paginated response format ***/
export interface PageResponse<T> {
  records: T[];
  total: number;
  current: number;
  size: number;
}

/*** 分页请求参数 - Pagination request params ***/
export interface PageParams {
  current: number;
  size: number;
}

/*** 登录响应 - Login response ***/
export interface LoginResponse {
  token: string;
  accessToken?: string; // 兼容两种字段名
  refreshToken: string;
  expiresIn?: number;
  user?: UserInfo;
}

/*** Token 响应 - Token response ***/
export interface TokenResponse {
  token: string;
  accessToken?: string; // 兼容两种字段名
  refreshToken: string;
  expiresIn?: number;
}

/*** 用户信息 - User info ***/
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  role: string;
}

/*** 项目列表参数 - Project list params ***/
export interface ProjectListParams extends PageParams {
  name?: string;
  status?: string;
  style?: string;
}

/*** 项目基础信息 - Project basic info ***/
export interface Project {
  id: number;
  name: string;
  description: string;
  area: number;
  budget: number;
  style: string;
  status: "draft" | "in_progress" | "completed" | "archived";
  createdAt: string;
  updatedAt: string;
}

/*** 项目详情 - Project detail ***/
export interface ProjectDetail extends Project {
  documentCount: number;
  designCount: number;
  lastAnalysisTime: string;
  costEstimate?: CostEstimate;
}

/*** 成本估算 - Cost estimate ***/
export interface CostEstimate {
  projectId: number;
  materialCost: number;      // 材料费总计
  laborCost: number;         // 人工费
  equipmentCost: number;     // 设备费
  managementCost: number;    // 管理费
  customCosts: CustomCost[]; // 自定义费用
  customCostTotal: number;   // 自定义费用总计
  totalCost: number;         // 总成本
  budgetLimit: number;       // 预算限额
  budgetExceeded: boolean;   // 是否超预算
  exceededAmount: number;    // 超出金额
  materialItems?: MaterialItem[]; // 材料明细列表
}

/*** 自定义费用项 - Custom cost item ***/
export interface CustomCost {
  name: string;
  amount: number;
}

/*** 材料明细项 - Material item ***/
export interface MaterialItem {
  id: number;
  projectId: number;
  materialId?: number;
  name: string;
  category: string;
  specification: string;
  unit: string;
  unitPrice: number;
  quantity: number;
  totalPrice: number;
  brand: string;
  supplier: string;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

/*** 成本明细项（兼容旧结构） - Cost item ***/
export interface CostItem {
  category: string;
  name: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
}

/*** 文档信息 - Document info ***/
export interface Document {
  id: number;
  projectId: number;
  fileName: string;
  filePath: string;  // 文件存储路径 - File storage path
  fileType: string;
  fileSize: number;
  analysisStatus: "pending" | "analyzing" | "completed" | "failed";
  createdAt: string;
  errorMessage: string;
  keywords?: string[];
  summary?: string;
  projectName?: string;
  extractedArea?: number;
  extractedBudget?: number;
  extractedStyle?: string;
  functionalZones?: string[];
}

/*** 文档摘要 - Document summary ***/
export interface DocumentSummary {
  overview: string;
  requirements: string;
  special: string;
}

/*** 材料分类 - Material category ***/
export interface MaterialCategory {
  name: string;
  count: number;
  icon: string;
}

/*** 材料信息 - Material info ***/
export interface Material {
  id: number;
  name: string;
  category: string;
  brand: string;
  price: number;
  unit: string;
  imageUrl: string;
}

/*** 材料详情 - Material detail ***/
export interface MaterialDetail extends Material {
  specification: string;
  supplier: string;
  applicableScenes: string[];
  images: string[];
}

/*** 材料列表参数 - Material list params ***/
export interface MaterialListParams extends PageParams {
  category?: string;
  keyword?: string;
}

/*** 聊天请求 - Chat request ***/
export interface ChatRequest {
  query: string;
  projectId?: number;
  sessionId?: string;
}

/*** 聊天响应 - Chat response ***/
export interface ChatResponse {
  answer: string;
  sessionId: string;
  messageId: number;
  tokensUsed: number;
}

/*** 聊天消息 - Chat message ***/
export interface ChatMessage {
  id: number;
  role: "user" | "assistant";
  content: string;
  createdAt: string;
}

/*** 聊天会话 - Chat session ***/
export interface ChatSession {
  id: string;
  projectId?: number;
  title: string;
  createdAt: string;
  updatedAt: string;
}

/*** 流式回调 - Stream callbacks ***/
export interface StreamCallbacks {
  onChunk: (content: string) => void;
  onDone: (result: { sessionId: string; messageId: number }) => void;
  onError: (error: string) => void;
}

/*** 通知信息 - Notification info ***/
export interface Notification {
  id: number;
  title: string;
  content: string;
  type: "project" | "system" | "approval";
  relatedId?: number;
  isRead: boolean;
  createdAt: string;
}

/*** 设计图信息 - Design image info ***/
export interface DesignImage {
  id: number;
  projectId: number;
  name: string;
  description?: string;
  imageUrl: string;
  thumbnailUrl?: string;
  fileSize: number;
  width?: number;
  height?: number;
  createdAt: string;
  updatedAt: string;
}

/*** 设计图列表参数 - Design image list params ***/
export interface DesignImageListParams extends PageParams {
  projectId: number;
}
