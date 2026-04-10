import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, ChevronRight, FileText, Pencil, X, Check, Info, CalendarDays, ChevronDown, ChevronUp } from 'lucide-react'
import { jobsService } from '@/services/jobs'
import { applicationsService } from '@/services/applications'
import type { Job, Application, JobStage, JobStatus } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import StatusBadge from '@/components/StatusBadge'
import PipelineBar from '@/components/PipelineBar'

// ---- Job status config -------------------------------------------------------

const STATUS_TRANSITIONS: Record<JobStatus, JobStatus[]> = {
  open: ['paused', 'closed', 'cancelled'],
  paused: ['open', 'closed', 'cancelled'],
  closed: [],
  cancelled: [],
}

const STATUS_LABELS: Record<JobStatus, string> = {
  open: 'Aberta',
  paused: 'Pausada',
  closed: 'Encerrada',
  cancelled: 'Cancelada',
}

// ---- Confirmation dialog types -----------------------------------------------

type PendingAction =
  | { type: 'advance'; app: Application; current: JobStage | null; next: JobStage }
  | { type: 'accept'; app: Application }
  | { type: 'reject'; app: Application }
  | { type: 'job-status'; newStatus: JobStatus }

// ---- Stage helpers -----------------------------------------------------------

function formatDate(isoString: string): string {
  const d = new Date(isoString)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  return `${day}/${month}/${year} às ${hours}:${minutes}`
}

function getSortedStages(job: Job): JobStage[] {
  return [...(job.stages ?? [])].sort((a, b) => a.order_index - b.order_index)
}

function getNextStage(job: Job, currentStageId?: string): JobStage | null {
  const stages = getSortedStages(job)
  if (stages.length === 0) return null
  if (!currentStageId) return stages[0]
  const idx = stages.findIndex((s) => s.id === currentStageId)
  if (idx === -1 || idx + 1 >= stages.length) return null
  return stages[idx + 1]
}

function getCurrentStage(job: Job, currentStageId?: string): JobStage | null {
  if (!currentStageId) return null
  return job.stages?.find((s) => s.id === currentStageId) ?? null
}

// ---- Confirmation dialog text ------------------------------------------------

function getConfirmContent(action: PendingAction): { title: string; description: React.ReactNode; confirmLabel: string; destructive?: boolean } {
  switch (action.type) {
    case 'advance': {
      const candidateName = action.app.candidate?.name || action.app.candidate?.email || 'o candidato'
      const currentLabel = action.current ? `"${action.current.name}"` : 'sem etapa'
      const nextLabel = `"${action.next.name}"`
      return {
        title: 'Avançar etapa?',
        description: (
          <>
            Mover <strong>{candidateName}</strong> de {currentLabel} para{' '}
            <strong>{nextLabel}</strong>?
          </>
        ),
        confirmLabel: 'Avançar',
      }
    }
    case 'accept': {
      const candidateName = action.app.candidate?.name || action.app.candidate?.email || 'o candidato'
      return {
        title: 'Aprovar candidatura?',
        description: (
          <>
            Tem certeza que deseja <strong>aprovar</strong> a candidatura de{' '}
            <strong>{candidateName}</strong>? Esta ação refletirá para o candidato.
          </>
        ),
        confirmLabel: 'Aprovar',
      }
    }
    case 'reject': {
      const candidateName = action.app.candidate?.name || action.app.candidate?.email || 'o candidato'
      return {
        title: 'Reprovar candidatura?',
        description: (
          <>
            Tem certeza que deseja <strong>reprovar</strong> a candidatura de{' '}
            <strong>{candidateName}</strong>? Esta ação refletirá para o candidato.
          </>
        ),
        confirmLabel: 'Reprovar',
        destructive: true,
      }
    }
    case 'job-status':
      return {
        title: 'Alterar status da vaga?',
        description: (
          <>
            Deseja alterar o status da vaga para{' '}
            <strong>{STATUS_LABELS[action.newStatus]}</strong>?
          </>
        ),
        confirmLabel: 'Confirmar',
      }
  }
}

// ---- Main page ---------------------------------------------------------------

// ---- Application card --------------------------------------------------------

interface ApplicationCardProps {
  app: Application
  job: Job
  onAction: (action: PendingAction) => void
}

function ApplicationCard({ app, job, onAction }: ApplicationCardProps) {
  const [coverOpen, setCoverOpen] = useState(false)
  const isFinished = ['accepted', 'rejected', 'withdrawn'].includes(app.status)
  const nextStage = getNextStage(job, app.current_stage_id)
  const currentStage = getCurrentStage(job, app.current_stage_id)

  return (
    <Card className={app.status === 'withdrawn' ? 'opacity-60' : ''}>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div>
            <CardTitle className="text-base">
              {app.candidate?.name || app.candidate?.email || 'Candidato'}
            </CardTitle>
            <div className="flex items-center gap-3 mt-0.5 flex-wrap">
              {app.candidate?.name && app.candidate?.email && (
                <span className="text-xs text-muted-foreground">{app.candidate.email}</span>
              )}
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <CalendarDays className="w-3 h-3" />
                {formatDate(app.created_at)}
              </span>
            </div>
          </div>
          <StatusBadge type="application" status={app.status} />
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Pipeline */}
        {job.stages && job.stages.length > 0 && (
          <PipelineBar stages={job.stages} currentStageId={app.current_stage_id} />
        )}

        {/* Quick-access: CV + cover letter toggle */}
        <div className="flex items-center gap-2 flex-wrap">
          {app.cv_url ? (
            <a
              href={app.cv_url}
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border text-sm font-medium hover:bg-accent transition-colors"
            >
              <FileText className="w-4 h-4 text-primary" />
              Ver CV
            </a>
          ) : (
            <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border text-sm text-muted-foreground bg-muted/40 cursor-not-allowed">
              <FileText className="w-4 h-4" />
              Sem CV
            </span>
          )}
          <button
            type="button"
            onClick={() => setCoverOpen((v) => !v)}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border text-sm font-medium hover:bg-accent transition-colors"
          >
            {coverOpen ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
            Carta de Apresentação
          </button>
        </div>

        {/* Cover letter expanded */}
        {coverOpen && (
          <div className="rounded-md bg-muted/40 border p-3">
            <p className="text-sm text-muted-foreground whitespace-pre-line">
              {app.cover_letter || <em>Não informada.</em>}
            </p>
          </div>
        )}

        {/* Action buttons */}
        {!isFinished && (
          <div className="flex gap-2 flex-wrap pt-1 border-t">
            {nextStage && (
              <Button
                size="sm"
                variant="outline"
                onClick={() => onAction({ type: 'advance', app, current: currentStage, next: nextStage })}
              >
                Avançar Etapa
                <ChevronRight className="w-4 h-4 ml-1" />
              </Button>
            )}
            <Button
              size="sm"
              className="bg-green-600 hover:bg-green-700"
              onClick={() => onAction({ type: 'accept', app })}
            >
              Aprovar
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={() => onAction({ type: 'reject', app })}
            >
              Reprovar
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

// ---- Main page ---------------------------------------------------------------

export default function JobPipelinePage() {
  const { id } = useParams<{ id: string }>()
  const [job, setJob] = useState<Job | null>(null)
  const [applications, setApplications] = useState<Application[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [actionError, setActionError] = useState('')
  const [actionLoading, setActionLoading] = useState(false)

  // Pending confirmation
  const [pending, setPending] = useState<PendingAction | null>(null)
  const [pendingStatusValue, setPendingStatusValue] = useState<string>('')

  // Salary edit state
  const [editingSalary, setEditingSalary] = useState(false)
  const [salaryMin, setSalaryMin] = useState('')
  const [salaryMax, setSalaryMax] = useState('')
  const [salaryError, setSalaryError] = useState('')
  const [salaryLoading, setSalaryLoading] = useState(false)

  async function load() {
    if (!id) return
    setIsLoading(true)
    try {
      const [j, apps] = await Promise.all([
        jobsService.get(id),
        applicationsService.jobApplications(id),
      ])
      setJob(j)
      setApplications(apps)
    } catch {
      setError('Erro ao carregar dados da vaga.')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => { load() }, [id])

  // Intercept status select: open confirm dialog instead of applying immediately
  function handleStatusSelectChange(value: string) {
    if (!job) return
    setPendingStatusValue(value)
    setPending({ type: 'job-status', newStatus: value as JobStatus })
    setActionError('')
  }

  async function executeAction() {
    if (!pending) return
    setActionLoading(true)
    setActionError('')
    try {
      if (pending.type === 'advance') {
        const updated = await applicationsService.advanceStage(pending.app.id)
        setApplications((prev) => prev.map((a) => (a.id === updated.id ? updated : a)))
      } else if (pending.type === 'accept') {
        const updated = await applicationsService.updateStatus(pending.app.id, 'accepted')
        setApplications((prev) => prev.map((a) => (a.id === updated.id ? updated : a)))
      } else if (pending.type === 'reject') {
        const updated = await applicationsService.updateStatus(pending.app.id, 'rejected')
        setApplications((prev) => prev.map((a) => (a.id === updated.id ? updated : a)))
      } else if (pending.type === 'job-status') {
        const updated = await jobsService.update(job!.id, { status: pending.newStatus })
        setJob(updated)
        setPendingStatusValue('')
      }
      setPending(null)
    } catch {
      setActionError('Erro ao executar ação. Tente novamente.')
    } finally {
      setActionLoading(false)
    }
  }

  function cancelAction() {
    setPending(null)
    setPendingStatusValue('')
    setActionError('')
  }

  function startEditSalary() {
    if (!job) return
    setSalaryMin(job.salary_min != null ? String(job.salary_min) : '')
    setSalaryMax(job.salary_max != null ? String(job.salary_max) : '')
    setSalaryError('')
    setEditingSalary(true)
  }

  async function saveSalary() {
    if (!job) return
    const min = parseFloat(salaryMin)
    const max = parseFloat(salaryMax)
    if (isNaN(min) || isNaN(max) || min < 0 || min > max) {
      setSalaryError('Intervalo salarial inválido.')
      return
    }
    setSalaryLoading(true)
    try {
      const updated = await jobsService.update(job.id, { salary_min: min, salary_max: max })
      setJob(updated)
      setEditingSalary(false)
    } catch {
      setSalaryError('Erro ao atualizar salário.')
    } finally {
      setSalaryLoading(false)
    }
  }

  if (isLoading) return <p className="text-muted-foreground">Carregando...</p>
  if (error) return <p className="text-destructive">{error}</p>
  if (!job) return null

  const availableTransitions = STATUS_TRANSITIONS[job.status]
  const confirmContent = pending ? getConfirmContent(pending) : null

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3 flex-wrap">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/recruiter/jobs">
            <ArrowLeft className="w-4 h-4" />
          </Link>
        </Button>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <h1 className="text-2xl font-bold truncate">{job.title}</h1>
            <StatusBadge type="job" status={job.status} />
          </div>
          <p className="text-sm text-muted-foreground">
            {job.company && <span className="font-medium text-foreground">{job.company} &middot; </span>}
            {job.location}
          </p>
        </div>

        {availableTransitions.length > 0 && (
          <Select
            value={pendingStatusValue}
            onValueChange={handleStatusSelectChange}
          >
            <SelectTrigger className="w-44">
              <SelectValue placeholder="Alterar status" />
            </SelectTrigger>
            <SelectContent>
              {availableTransitions.map((s) => (
                <SelectItem key={s} value={s}>
                  {STATUS_LABELS[s]}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
      </div>

      {/* Salary */}
      <Card>
        <CardContent className="py-4">
          {editingSalary ? (
            <div className="space-y-2">
              <p className="text-sm font-medium">Editar faixa salarial</p>
              <div className="flex items-center gap-2 flex-wrap">
                <div className="flex items-center gap-1">
                  <span className="text-sm text-muted-foreground">R$</span>
                  <Input type="number" min={0} step={1} className="w-32" placeholder="Mínimo"
                    value={salaryMin} onChange={(e) => setSalaryMin(e.target.value)} />
                </div>
                <span className="text-muted-foreground">–</span>
                <div className="flex items-center gap-1">
                  <span className="text-sm text-muted-foreground">R$</span>
                  <Input type="number" min={0} step={1} className="w-32" placeholder="Máximo"
                    value={salaryMax} onChange={(e) => setSalaryMax(e.target.value)} />
                </div>
                <Button size="sm" onClick={saveSalary} disabled={salaryLoading}>
                  <Check className="w-4 h-4 mr-1" />Salvar
                </Button>
                <Button size="sm" variant="ghost" onClick={() => setEditingSalary(false)}>
                  <X className="w-4 h-4" />
                </Button>
              </div>
              {salaryError && <p className="text-sm text-destructive">{salaryError}</p>}
            </div>
          ) : (
            <div className="flex items-center gap-3">
              <div>
                <p className="text-xs text-muted-foreground uppercase tracking-wide">Faixa salarial</p>
                <p className="font-medium">
                  {job.salary_min != null && job.salary_max != null
                    ? `R$ ${job.salary_min.toLocaleString('pt-BR')} – R$ ${job.salary_max.toLocaleString('pt-BR')}`
                    : 'Não informado'}
                </p>
              </div>
              <Button variant="ghost" size="sm" onClick={startEditSalary} className="ml-auto">
                <Pencil className="w-4 h-4 mr-1" />Editar
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Pipeline stages (read-only) */}
      {job.stages && job.stages.length > 0 && (
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">Pipeline da Vaga</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-start gap-2 rounded-md bg-muted/50 border border-muted p-3 text-sm text-muted-foreground">
              <Info className="w-4 h-4 mt-0.5 shrink-0 text-blue-500" />
              <span>
                A pipeline é definida na criação da vaga e não pode ser alterada. Para usar uma pipeline diferente, crie uma nova vaga.
              </span>
            </div>
            <div className="flex flex-wrap gap-2">
              {getSortedStages(job).map((stage, idx, arr) => (
                <div key={stage.id} className="flex items-center gap-1">
                  <span className="inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium bg-background">
                    {idx + 1}. {stage.name}
                  </span>
                  {idx < arr.length - 1 && (
                    <ChevronRight className="w-4 h-4 text-muted-foreground" />
                  )}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Candidates */}
      <div>
        <h2 className="text-lg font-semibold mb-3">
          Candidaturas ({applications.length})
        </h2>

        {applications.length === 0 ? (
          <Card>
            <CardContent className="py-10 text-center text-muted-foreground">
              Nenhuma candidatura recebida ainda.
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {applications.map((app) => (
              <ApplicationCard
                key={app.id}
                app={app}
                job={job}
                onAction={(action) => { setActionError(''); setPending(action) }}
              />
            ))}
          </div>
        )}
      </div>

      {/* Confirmation dialog */}
      <Dialog open={!!pending} onOpenChange={(open) => { if (!open && !actionLoading) cancelAction() }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{confirmContent?.title}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">{confirmContent?.description}</p>
          {actionError && <p className="text-sm text-destructive">{actionError}</p>}
          <DialogFooter>
            <Button variant="outline" onClick={cancelAction} disabled={actionLoading}>
              Cancelar
            </Button>
            <Button
              variant={confirmContent?.destructive ? 'destructive' : 'default'}
              className={!confirmContent?.destructive && pending?.type === 'accept'
                ? 'bg-green-600 hover:bg-green-700' : undefined}
              onClick={executeAction}
              disabled={actionLoading}
            >
              {actionLoading ? 'Aguarde...' : confirmContent?.confirmLabel}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
