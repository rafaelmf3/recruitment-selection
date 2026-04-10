import { api } from './api'
import type {
  Job,
  CreateJobRequest,
  UpdateJobRequest,
  ListJobsParams,
  PaginatedResponse,
} from '@/types'

export const jobsService = {
  list: async (params?: ListJobsParams): Promise<PaginatedResponse<Job>> => {
    const res = await api.get<PaginatedResponse<Job>>('/jobs', { params })
    return res.data
  },

  get: async (id: string): Promise<Job> => {
    const res = await api.get<Job>(`/jobs/${id}`)
    return res.data
  },

  create: async (data: CreateJobRequest): Promise<Job> => {
    const res = await api.post<Job>('/recruiter/jobs', data)
    return res.data
  },

  update: async (id: string, data: UpdateJobRequest): Promise<Job> => {
    const res = await api.put<Job>(`/recruiter/jobs/${id}`, data)
    return res.data
  },

  myJobs: async (): Promise<Job[]> => {
    const res = await api.get<Job[]>('/recruiter/jobs')
    return res.data
  },
}
