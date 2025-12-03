export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  department_id?: number;
  department?: Department;
  role_id: number;
  role?: Role;
  manager_id?: number;
  manager?: User;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Department {
  id: number;
  name: string;
  description: string;
  parent_id?: number;
  parent?: Department;
  level: number;
  created_at: string;
  updated_at: string;
}

export interface Role {
  id: number;
  name: string;
  description: string;
  permissions: string;
  created_at: string;
  updated_at: string;
}

export interface KPIIndicator {
  id: number;
  name: string;
  description: string;
  category: string;
  data_type: string;
  unit: string;
  calculation_type: string;
  formula: string;
  applicable_roles: string;
  is_published: boolean;
  created_by: number;
  creator?: User;
  created_at: string;
  updated_at: string;
}

export interface AssessmentPlan {
  id: number;
  name: string;
  description: string;
  cycle_type: string;
  start_date: string;
  end_date: string;
  employee_id: number;
  employee?: User;
  manager_id: number;
  manager?: User;
  status: string;
  confirmed_at?: string;
  items?: AssessmentPlanItem[];
  created_at: string;
  updated_at: string;
}

export interface AssessmentPlanItem {
  id: number;
  plan_id: number;
  indicator_id: number;
  indicator?: KPIIndicator;
  weight: number;
  target_value: number;
  challenge_value: number;
  baseline_value: number;
  created_at: string;
  updated_at: string;
}

export interface KPIProgress {
  id: number;
  plan_item_id: number;
  plan_item?: AssessmentPlanItem;
  current_value: number;
  completion_rate: number;
  updated_by: number;
  updater?: User;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface PerformanceReview {
  id: number;
  plan_id: number;
  plan?: AssessmentPlan;
  employee_id: number;
  employee?: User;
  reviewer_id: number;
  reviewer?: User;
  self_score?: number;
  self_comment: string;
  manager_score?: number;
  manager_comment: string;
  final_score?: number;
  performance_level: string;
  status: string;
  submitted_at?: string;
  items?: PerformanceReviewItem[];
  created_at: string;
  updated_at: string;
}

export interface PerformanceReviewItem {
  id: number;
  review_id: number;
  plan_item_id: number;
  plan_item?: AssessmentPlanItem;
  actual_value: number;
  completion_rate: number;
  auto_score: number;
  adjusted_score?: number;
  adjustment_reason: string;
  created_at: string;
  updated_at: string;
}

export interface Review360Campaign {
  id: number;
  name: string;
  description: string;
  start_date: string;
  end_date: string;
  status: string;
  created_by: number;
  creator?: User;
  created_at: string;
  updated_at: string;
}

export interface Review360Assignment {
  id: number;
  campaign_id: number;
  campaign?: Review360Campaign;
  reviewee_id: number;
  reviewee?: User;
  reviewer_id: number;
  reviewer?: User;
  relationship: string;
  status: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Review360Response {
  id: number;
  assignment_id: number;
  assignment?: Review360Assignment;
  dimension: string;
  score: number;
  comment: string;
  created_at: string;
  updated_at: string;
}

export interface AuditLog {
  id: number;
  user_id: number;
  user?: User;
  action: string;
  resource: string;
  resource_id: number;
  ip_address: string;
  user_agent: string;
  details: string;
  status: string;
  created_at: string;
  updated_at: string;
}
