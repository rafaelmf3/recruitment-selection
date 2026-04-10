import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, ChevronRight, ChevronDown, ChevronUp, Pencil, Check, X, Search } from 'lucide-react'
import { jobsService } from '@/services/jobs'
import type { Job } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import StatusBadge from '@/components/StatusBadge'

interface JobCardProps {
  job: Job
  onUpdate: (updated: Job) => void
}

function JobCard({ job, onUpdate }: JobCardProps) {
  const [expanded, setExpanded] = useState(false)
  const [editingSalary, setEditingSalary] = useState(false)
  const [salaryMin, setSalaryMin] = useState('')
  const [salaryMax, setSalaryMax] = useState('')
  const [salaryError, setSalaryError] = useState('')
  const [salaryLoading, setSalaryLoading] = useState(false)

  function startEditSalary() {
    setSalaryMin(job.salary_min != null ? String(job.salary_min) : '')
    setSalaryMax(job.salary_max != null ? String(job.salary_max) : '')
    setSalaryError('')
    setEditingSalary(true)
  }

  async function saveSalary() {
    const min = parseFloat(salaryMin)
    const max = parseFloat(salaryMax)
    if (isNaN(min) || isNaN(max) || min < 0 || min > max) {
      setSalaryError('Intervalo salarial inválido.')
      return
    }
    setSalaryLoading(true)
    try {
      const updated = await jobsService.update(job.id, { salary_min: min, salary_max: max })
      onUpdate(updated)
      setEditingSalary(false)
    } catch {
      setSalaryError('Erro ao atualizar salário.')
    } finally {
      setSalaryLoading(false)
    }
  }

  return (
    <Card>
      <CardContent className="py-4 space-y-3">
        {/* Top row: title, status, pipeline link */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h2 className="font-semibold truncate">{job.title}</h2>
              <StatusBadge type="job" status={job.status} />
            </div>
            <p className="text-sm text-muted-foreground mt-0.5">
              {job.company && <span className="font-medium text-foreground">{job.company} &middot; </span>}
              {job.location}
            </p>
          </div>
          <div className="flex items-center gap-1 shrink-0">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setExpanded((v) => !v)}
              title={expanded ? 'Recolher descrição' : 'Expandir descrição'}
            >
              {expanded ? (
                <ChevronUp className="w-4 h-4" />
              ) : (
                <ChevronDown className="w-4 h-4" />
              )}
            </Button>
            <Button variant="ghost" size="sm" asChild>
              <Link to={`/recruiter/jobs/${job.id}/pipeline`}>
                Pipeline
                <ChevronRight className="w-4 h-4 ml-1" />
              </Link>
            </Button>
          </div>
        </div>

        {/* Salary row */}
        {editingSalary ? (
          <div className="space-y-1">
            <div className="flex items-center gap-2 flex-wrap">
              <div className="flex items-center gap-1">
                <span className="text-sm text-muted-foreground">R$</span>
                <Input
                  type="number"
                  min={0}
                  step={1}
                  className="w-28 h-8 text-sm"
                  placeholder="Mínimo"
                  value={salaryMin}
                  onChange={(e) => setSalaryMin(e.target.value)}
                />
              </div>
              <span className="text-muted-foreground">–</span>
              <div className="flex items-center gap-1">
                <span className="text-sm text-muted-foreground">R$</span>
                <Input
                  type="number"
                  min={0}
                  step={1}
                  className="w-28 h-8 text-sm"
                  placeholder="Máximo"
                  value={salaryMax}
                  onChange={(e) => setSalaryMax(e.target.value)}
                />
              </div>
              <Button size="sm" className="h-8" onClick={saveSalary} disabled={salaryLoading}>
                <Check className="w-3.5 h-3.5 mr-1" />
                Salvar
              </Button>
              <Button
                size="sm"
                variant="ghost"
                className="h-8"
                onClick={() => setEditingSalary(false)}
              >
                <X className="w-3.5 h-3.5" />
              </Button>
            </div>
            {salaryError && <p className="text-xs text-destructive">{salaryError}</p>}
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <p className="text-sm text-muted-foreground">
              <span className="font-medium text-foreground">Faixa salarial:</span>{' '}
              {job.salary_min != null && job.salary_max != null
                ? `R$ ${job.salary_min.toLocaleString('pt-BR')} – R$ ${job.salary_max.toLocaleString('pt-BR')}`
                : 'Não informado'}
            </p>
            {job.status !== 'closed' && job.status !== 'cancelled' && (
              <Button
                variant="ghost"
                size="sm"
                className="h-6 px-2"
                onClick={startEditSalary}
                title="Editar faixa salarial"
              >
                <Pencil className="w-3 h-3" />
              </Button>
            )}
          </div>
        )}

        {/* Expandable description */}
        {expanded && (
          <div className="pt-1 border-t">
            <p className="text-sm text-muted-foreground whitespace-pre-line">
              {job.description || 'Sem descrição.'}
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

const STATUS_OPTIONS = [
  { value: '', label: 'Todos os status' },
  { value: 'open', label: 'Aberta' },
  { value: 'paused', label: 'Pausada' },
  { value: 'closed', label: 'Encerrada' },
  { value: 'cancelled', label: 'Cancelada' },
]

export default function MyJobsPage() {
  const [jobs, setJobs] = useState<Job[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState('')
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  function fetchJobs(q: string, st: string) {
    setIsLoading(true)
    setError('')
    jobsService
      .myJobs({ q: q || undefined, status: st || undefined })
      .then(setJobs)
      .catch(() => setError('Erro ao carregar vagas.'))
      .finally(() => setIsLoading(false))
  }

  // Initial load
  useEffect(() => {
    fetchJobs('', '')
  }, [])

  function handleSearchChange(value: string) {
    setSearch(value)
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      fetchJobs(value, status)
    }, 350)
  }

  function handleStatusChange(value: string) {
    setStatus(value)
    fetchJobs(search, value)
  }

  function handleJobUpdate(updated: Job) {
    setJobs((prev) => prev.map((j) => (j.id === updated.id ? updated : j)))
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Minhas Vagas</h1>
        <Button asChild>
          <Link to="/recruiter/jobs/new">
            <Plus className="w-4 h-4 mr-2" />
            Nova Vaga
          </Link>
        </Button>
      </div>

      {/* Search + filter bar */}
      <div className="flex gap-2 flex-wrap">
        <div className="relative flex-1 min-w-48">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground pointer-events-none" />
          <Input
            className="pl-8"
            placeholder="Buscar por nome ou empresa..."
            value={search}
            onChange={(e) => handleSearchChange(e.target.value)}
          />
        </div>
        <select
          className="border rounded-md px-3 py-2 text-sm bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          value={status}
          onChange={(e) => handleStatusChange(e.target.value)}
        >
          {STATUS_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      </div>

      {isLoading && <p className="text-muted-foreground">Carregando...</p>}
      {error && <p className="text-destructive">{error}</p>}

      {!isLoading && !error && jobs.length === 0 && (
        <Card>
          <CardContent className="py-12 text-center text-muted-foreground">
            {search || status
              ? 'Nenhuma vaga encontrada para os filtros informados.'
              : <>
                  Nenhuma vaga cadastrada ainda.{' '}
                  <Link to="/recruiter/jobs/new" className="text-primary hover:underline">
                    Crie a primeira!
                  </Link>
                </>
            }
          </CardContent>
        </Card>
      )}

      <div className="space-y-3">
        {jobs.map((job) => (
          <JobCard key={job.id} job={job} onUpdate={handleJobUpdate} />
        ))}
      </div>
    </div>
  )
}
