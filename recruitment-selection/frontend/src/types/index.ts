// ---- Enums ------------------------------------------------------------------

export type UserRole = 'recruiter' | 'candidate'

export type JobStatus = 'open' | 'paused' | 'closed' | 'cancelled'

export type ApplicationStatus = 'pending' | 'in_progress' | 'accepted' | 'rejected' | 'withdrawn'

// ---- Domain models ----------------------------------------------------------

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  created_at: string
}

export interface JobStage {
  id: string
  job_id: string
  name: string
  order_index: number
}

export interface Job {
  id: string
  company?: string
  title: string
  description: string
  location: string
  salary_min?: number | null
  salary_max?: number | null
  status: JobStatus
  recruiter_id: string
  recruiter?: User
  stages?: JobStage[]
  created_at: string
  updated_at: string
}

export interface Application {
  id: string
  job_id: string
  candidate_id: string
  cover_letter: string
  cv_url?: string
  status: ApplicationStatus
  current_stage_id?: string
  current_stage?: JobStage
  job?: Job
  candidate?: User
  created_at: string
  updated_at: string
}

// ---- Auth DTOs --------------------------------------------------------------

export interface RegisterRequest {
  name: string
  email: string
  password: string
  role: UserRole
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}

// ---- Job DTOs ---------------------------------------------------------------

export interface CreateJobRequest {
  company?: string
  title: string
  description: string
  location: string
  salary_min: number
  salary_max: number
  stages?: Array<{ name: string; order_index: number }>
}

export interface UpdateJobRequest {
  company?: string
  title?: string
  description?: string
  location?: string
  salary_min?: number
  salary_max?: number
  status?: JobStatus
  stages?: Array<{ name: string; order_index: number }>
}

export interface ListJobsParams {
  q?: string
  status?: JobStatus
  page?: number
  limit?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  pages: number
}

// ---- Application DTOs -------------------------------------------------------

export interface ApplyRequest {
  cover_letter: string
  cv?: File
}
