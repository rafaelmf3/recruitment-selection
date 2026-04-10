import { useEffect, useState, useCallback } from 'react'
import { Search, MapPin, DollarSign, EyeOff, Eye, CalendarCheck } from 'lucide-react'
import { jobsService } from '@/services/jobs'
import { applicationsService } from '@/services/applications'
import type { Job, PaginatedResponse } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import StatusBadge from '@/components/StatusBadge'

const PAGE_SIZE = 10

// Map jobId → ISO applied_at string
type AppliedMap = Map<string, string>

function formatAppliedAt(isoString: string): string {
  const date = new Date(isoString)
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `Candidatou-se ${day}/${month}/${year} às ${hours}:${minutes}`
}

export default function BrowseJobsPage() {
  const [result, setResult] = useState<PaginatedResponse<Job> | null>(null)
  const [query, setQuery] = useState('')
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  // Apply dialog state
  const [applyJob, setApplyJob] = useState<Job | null>(null)
  const [coverLetter, setCoverLetter] = useState('')
  const [cvFile, setCvFile] = useState<File | null>(null)
  const [isApplying, setIsApplying] = useState(false)
  const [applyError, setApplyError] = useState('')

  // Applied jobs: jobId → applied_at
  const [appliedMap, setAppliedMap] = useState<AppliedMap>(new Map())
  const [hideApplied, setHideApplied] = useState(false)

  // Load existing applications once on mount to build the appliedMap.
  // Errors are intentionally ignored — the jobs list still renders; applied
  // badges simply won't appear until the next successful load.
  useEffect(() => {
    applicationsService.myApplications()
      .then((apps) => {
        const map: AppliedMap = new Map()
        apps.forEach((app) => {
          if (app.job_id) map.set(app.job_id, app.created_at)
        })
        setAppliedMap(map)
      })
      .catch(() => {/* non-critical — appliedMap stays empty */})
  }, [])

  const fetchJobs = useCallback(async () => {
    setIsLoading(true)
    setError('')
    try {
      const data = await jobsService.list({ q: query || undefined, status: 'open', page, limit: PAGE_SIZE })
      setResult(data)
    } catch {
      setError('Erro ao carregar vagas.')
    } finally {
      setIsLoading(false)
    }
  }, [query, page])

  useEffect(() => {
    fetchJobs()
  }, [fetchJobs])

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    setPage(1)
    fetchJobs()
  }

  async function handleApply(e: React.FormEvent) {
    e.preventDefault()
    if (!applyJob) return
    setIsApplying(true)
    setApplyError('')
    try {
      const app = await applicationsService.apply(applyJob.id, coverLetter, cvFile ?? undefined)
      setAppliedMap((prev) => {
        const next = new Map(prev)
        next.set(applyJob.id, app.created_at)
        return next
      })
      setApplyJob(null)
      setCoverLetter('')
      setCvFile(null)
    } catch {
      setApplyError('Erro ao se candidatar. Você já pode ter se candidatado a esta vaga.')
    } finally {
      setIsApplying(false)
    }
  }

  const visibleJobs = hideApplied
    ? (result?.data ?? []).filter((job) => !appliedMap.has(job.id))
    : (result?.data ?? [])

  const appliedCount = (result?.data ?? []).filter((job) => appliedMap.has(job.id)).length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-2">
        <h1 className="text-2xl font-bold">Vagas Disponíveis</h1>
        {appliedCount > 0 && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => setHideApplied((v) => !v)}
          >
            {hideApplied ? (
              <>
                <Eye className="w-4 h-4 mr-2" />
                Mostrar candidaturas ({appliedCount})
              </>
            ) : (
              <>
                <EyeOff className="w-4 h-4 mr-2" />
                Ocultar candidaturas ({appliedCount})
              </>
            )}
          </Button>
        )}
      </div>

      <form onSubmit={handleSearch} className="flex gap-2">
        <Input
          placeholder="Buscar vagas..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-1"
        />
        <Button type="submit" size="icon" variant="outline">
          <Search className="w-4 h-4" />
        </Button>
      </form>

      {isLoading && <p className="text-muted-foreground">Carregando...</p>}
      {error && <p className="text-destructive">{error}</p>}

      {!isLoading && result && visibleJobs.length === 0 && (
        <p className="text-muted-foreground">
          {hideApplied && appliedCount > 0
            ? 'Você já se candidatou a todas as vagas desta página.'
            : 'Nenhuma vaga encontrada.'}
        </p>
      )}

      <div className="space-y-3">
        {visibleJobs.map((job) => {
          const appliedAt = appliedMap.get(job.id)
          return (
            <Card key={job.id} className={appliedAt ? 'border-green-500/40' : ''}>
              <CardHeader className="pb-2">
                <div className="flex items-start justify-between gap-2 flex-wrap">
                  <div>
                    <CardTitle className="text-base">{job.title}</CardTitle>
                    {job.company && (
                      <p className="text-sm font-medium text-foreground mt-0.5">{job.company}</p>
                    )}
                    <CardDescription className="flex items-center gap-3 mt-1 flex-wrap">
                      <span className="flex items-center gap-1">
                        <MapPin className="w-3 h-3" />
                        {job.location}
                      </span>
                      {job.salary_min != null && job.salary_max != null && (
                        <span className="flex items-center gap-1">
                          <DollarSign className="w-3 h-3" />
                          R$ {job.salary_min.toLocaleString('pt-BR')} &ndash; R${' '}
                          {job.salary_max.toLocaleString('pt-BR')}
                        </span>
                      )}
                    </CardDescription>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <StatusBadge type="job" status={job.status} />
                    {appliedAt && (
                      <span className="flex items-center gap-1 text-xs text-green-600 font-medium">
                        <CalendarCheck className="w-3 h-3" />
                        {formatAppliedAt(appliedAt)}
                      </span>
                    )}
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground line-clamp-2 mb-3">{job.description}</p>
                <Button
                  size="sm"
                  variant={appliedAt ? 'secondary' : 'default'}
                  disabled={!!appliedAt}
                  onClick={() => {
                    setApplyJob(job)
                    setCoverLetter('')
                    setCvFile(null)
                    setApplyError('')
                  }}
                >
                  {appliedAt ? 'Candidatura enviada' : 'Candidatar-se'}
                </Button>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Pagination */}
      {result && result.pages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
          >
            Anterior
          </Button>
          <span className="text-sm text-muted-foreground">
            {page} / {result.pages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page >= result.pages}
            onClick={() => setPage((p) => p + 1)}
          >
            Próximo
          </Button>
        </div>
      )}

      {/* Apply dialog */}
      <Dialog open={!!applyJob} onOpenChange={(open) => !open && setApplyJob(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Candidatura: {applyJob?.title}</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleApply} className="space-y-4 mt-2">
            {applyError && <p className="text-sm text-destructive">{applyError}</p>}

            <div className="space-y-2">
              <Label htmlFor="cover_letter">Carta de Apresentação</Label>
              <Textarea
                id="cover_letter"
                required
                rows={5}
                placeholder="Fale sobre você e por que se interessa por esta vaga..."
                value={coverLetter}
                onChange={(e) => setCoverLetter(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="cv">
                Currículo (PDF) <span className="text-destructive">*</span>
              </Label>
              <Input
                id="cv"
                type="file"
                accept=".pdf"
                required
                onChange={(e) => setCvFile(e.target.files?.[0] ?? null)}
              />
              <p className="text-xs text-muted-foreground">Apenas arquivos PDF são aceitos.</p>
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setApplyJob(null)}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isApplying}>
                {isApplying ? 'Enviando...' : 'Enviar Candidatura'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
