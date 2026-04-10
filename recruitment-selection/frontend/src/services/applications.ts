import { api } from './api'
import type { Application, ApplicationStatus } from '@/types'

export const applicationsService = {
  // Candidate: apply to a job (multipart form)
  apply: async (jobId: string, coverLetter: string, cv?: File): Promise<Application> => {
    const form = new FormData()
    form.append('cover_letter', coverLetter)
    if (cv) form.append('cv', cv)

    const res = await api.post<Application>(`/jobs/${jobId}/apply`, form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return res.data
  },

  // Candidate: list own applications
  myApplications: async (): Promise<Application[]> => {
    const res = await api.get<Application[]>('/applications')
    return res.data
  },

  // Candidate: withdraw from a job application
  withdraw: async (applicationId: string): Promise<Application> => {
    const res = await api.patch<Application>(`/applications/${applicationId}/withdraw`)
    return res.data
  },

  // Recruiter: list applications for a job
  jobApplications: async (jobId: string): Promise<Application[]> => {
    const res = await api.get<Application[]>(`/recruiter/jobs/${jobId}/applications`)
    return res.data
  },

  // Recruiter: advance application to next stage
  advanceStage: async (applicationId: string): Promise<Application> => {
    const res = await api.patch<Application>(`/recruiter/applications/${applicationId}/stage`)
    return res.data
  },

  // Recruiter: update application status (accepted/rejected)
  updateStatus: async (applicationId: string, status: ApplicationStatus): Promise<Application> => {
    const res = await api.patch<Application>(`/recruiter/applications/${applicationId}/status`, {
      status,
    })
    return res.data
  },
}
