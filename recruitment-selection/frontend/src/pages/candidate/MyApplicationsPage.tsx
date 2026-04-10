import { useEffect, useState } from 'react'
import { MapPin, DollarSign, ChevronDown, ChevronUp, FileText, File, LogOut } from 'lucide-react'
import { applicationsService } from '@/services/applications'
import type { Application } from '@/types'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import StatusBadge from '@/components/StatusBadge'
import PipelineBar from '@/components/PipelineBar'

function formatAppliedAt(isoString: string): string {
  const date = new Date(isoString)
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `Candidatou-se em ${day}/${month}/${year} às ${hours}:${minutes}`
}

interface ApplicationCardProps {
  app: Application
  onUpdate: (updated: Application) => void
}

function ApplicationCard({ app, onUpdate }: ApplicationCardProps) {
  const [expanded, setExpanded] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [withdrawing, setWithdrawing] = useState(false)
  const [withdrawError, setWithdrawError] = useState('')

  const isWithdrawn = app.status === 'withdrawn'

  async function handleWithdraw() {
    setWithdrawing(true)
    setWithdrawError('')
    try {
      const updated = await applicationsService.withdraw(app.id)
      onUpdate(updated)
      setConfirmOpen(false)
    } catch {
      setWithdrawError('Erro ao retirar candidatura. Tente novamente.')
    } finally {
      setWithdrawing(false)
    }
  }

  return (
    <>
      <Card className={isWithdrawn ? 'opacity-60' : ''}>
        <CardHeader className="pb-2">
          <div className="flex items-start justify-between gap-2 flex-wrap">
            <div className="space-y-0.5">
              <CardTitle className="text-base">{app.job?.title ?? 'Vaga'}</CardTitle>
              {app.job?.company && (
                <p className="text-sm font-medium text-foreground">{app.job.company}</p>
              )}
              <CardDescription className="flex items-center gap-3 flex-wrap mt-1">
                {app.job?.location && (
                  <span className="flex items-center gap-1">
                    <MapPin className="w-3 h-3" />
                    {app.job.location}
                  </span>
                )}
                {app.job?.salary_min != null && app.job?.salary_max != null && (
                  <span className="flex items-center gap-1">
                    <DollarSign className="w-3 h-3" />
                    R$ {app.job.salary_min.toLocaleString('pt-BR')} &ndash; R${' '}
                    {app.job.salary_max.toLocaleString('pt-BR')}
                  </span>
                )}
              </CardDescription>
              <p className="text-xs text-muted-foreground pt-0.5">
                {formatAppliedAt(app.created_at)}
              </p>
            </div>
            <div className="flex flex-col items-end gap-1 shrink-0">
              <StatusBadge type="application" status={app.status} />
              {app.job && <StatusBadge type="job" status={app.job.status} />}
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-3">
          {/* Job description snippet */}
          {app.job?.description && (
            <p className="text-sm text-muted-foreground line-clamp-2">{app.job.description}</p>
          )}

          {/* Pipeline progress */}
          {app.job?.stages && app.job.stages.length > 0 && (
            <PipelineBar stages={app.job.stages} currentStageId={app.current_stage_id} />
          )}
          {app.current_stage && (
            <p className="text-sm text-muted-foreground">
              Etapa atual:{' '}
              <span className="font-medium text-foreground">{app.current_stage.name}</span>
            </p>
          )}

          {/* Actions row */}
          <div className="flex items-center justify-between gap-2 pt-1">
            <Button
              variant="ghost"
              size="sm"
              className="gap-1 text-muted-foreground hover:text-foreground"
              onClick={() => setExpanded((v) => !v)}
            >
              {expanded ? (
                <>
                  <ChevronUp className="w-4 h-4" />
                  Ocultar candidatura
                </>
              ) : (
                <>
                  <ChevronDown className="w-4 h-4" />
                  Ver dados da candidatura
                </>
              )}
            </Button>

            {!isWithdrawn && (
              <Button
                variant="ghost"
                size="sm"
                className="gap-1 text-destructive hover:text-destructive hover:bg-destructive/10"
                onClick={() => { setWithdrawError(''); setConfirmOpen(true) }}
              >
                <LogOut className="w-4 h-4" />
                Retirar candidatura
              </Button>
            )}
          </div>

          {/* Expanded: cover letter + CV */}
          {expanded && (
            <div className="space-y-3 pt-2 border-t">
              <div className="space-y-1">
                <div className="flex items-center gap-1.5 text-sm font-medium">
                  <FileText className="w-4 h-4 text-muted-foreground" />
                  Carta de Apresentação
                </div>
                {app.cover_letter ? (
                  <p className="text-sm text-muted-foreground whitespace-pre-line pl-5">
                    {app.cover_letter}
                  </p>
                ) : (
                  <p className="text-sm text-muted-foreground italic pl-5">Não informada.</p>
                )}
              </div>

              <div className="space-y-1">
                <div className="flex items-center gap-1.5 text-sm font-medium">
                  <File className="w-4 h-4 text-muted-foreground" />
                  Currículo
                </div>
                {app.cv_url ? (
                  <a
                    href={app.cv_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-primary hover:underline pl-5 block"
                  >
                    Ver currículo enviado
                  </a>
                ) : (
                  <p className="text-sm text-muted-foreground italic pl-5">Não enviado.</p>
                )}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Withdrawal confirmation dialog */}
      <Dialog open={confirmOpen} onOpenChange={(open) => !withdrawing && setConfirmOpen(open)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Retirar candidatura?</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Você tem certeza que deseja retirar sua candidatura para a vaga{' '}
            <strong className="text-foreground">{app.job?.title ?? 'esta vaga'}</strong>? O
            recrutador verá o status atualizado e esta ação não pode ser desfeita.
          </p>
          {withdrawError && <p className="text-sm text-destructive">{withdrawError}</p>}
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setConfirmOpen(false)}
              disabled={withdrawing}
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleWithdraw}
              disabled={withdrawing}
            >
              {withdrawing ? 'Retirando...' : 'Sim, retirar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}

export default function MyApplicationsPage() {
  const [applications, setApplications] = useState<Application[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    applicationsService
      .myApplications()
      .then(setApplications)
      .catch(() => setError('Erro ao carregar candidaturas.'))
      .finally(() => setIsLoading(false))
  }, [])

  function handleUpdate(updated: Application) {
    setApplications((prev) => prev.map((a) => (a.id === updated.id ? updated : a)))
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Minhas Candidaturas</h1>

      {isLoading && <p className="text-muted-foreground">Carregando...</p>}
      {error && <p className="text-destructive">{error}</p>}

      {!isLoading && !error && applications.length === 0 && (
        <p className="text-muted-foreground">Você ainda não se candidatou a nenhuma vaga.</p>
      )}

      <div className="space-y-3">
        {applications.map((app) => (
          <ApplicationCard key={app.id} app={app} onUpdate={handleUpdate} />
        ))}
      </div>
    </div>
  )
}
