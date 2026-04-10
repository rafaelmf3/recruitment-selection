import { cn } from '@/lib/utils'
import type { ApplicationStatus, JobStatus } from '@/types'

const jobStatusConfig: Record<JobStatus, { label: string; className: string }> = {
  open: { label: 'Aberta', className: 'bg-green-100 text-green-800' },
  paused: { label: 'Pausada', className: 'bg-yellow-100 text-yellow-800' },
  closed: { label: 'Encerrada', className: 'bg-gray-100 text-gray-600' },
  cancelled: { label: 'Cancelada', className: 'bg-red-100 text-red-700' },
}

const applicationStatusConfig: Record<ApplicationStatus, { label: string; className: string }> = {
  pending: { label: 'Pendente', className: 'bg-blue-100 text-blue-800' },
  in_progress: { label: 'Em andamento', className: 'bg-purple-100 text-purple-800' },
  accepted: { label: 'Aceita', className: 'bg-green-100 text-green-800' },
  rejected: { label: 'Recusada', className: 'bg-red-100 text-red-700' },
  withdrawn: { label: 'Retirada', className: 'bg-gray-100 text-gray-600' },
}

interface JobStatusBadgeProps {
  type: 'job'
  status: JobStatus
}

interface ApplicationStatusBadgeProps {
  type: 'application'
  status: ApplicationStatus
}

type StatusBadgeProps = JobStatusBadgeProps | ApplicationStatusBadgeProps

export default function StatusBadge({ type, status }: StatusBadgeProps) {
  const config =
    type === 'job'
      ? jobStatusConfig[status as JobStatus]
      : applicationStatusConfig[status as ApplicationStatus]

  if (!config) return null

  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        config.className
      )}
    >
      {config.label}
    </span>
  )
}
