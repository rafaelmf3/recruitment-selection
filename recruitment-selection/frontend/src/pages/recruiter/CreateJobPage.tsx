import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Plus, Trash2, ArrowLeft, ArrowUp, ArrowDown } from 'lucide-react'
import { jobsService } from '@/services/jobs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

interface StageField {
  name: string
  order_index: number
}

const defaultStages: StageField[] = [
  { name: 'Triagem', order_index: 1 },
  { name: 'Entrevista RH', order_index: 2 },
  { name: 'Entrevista Tecnica', order_index: 3 },
  { name: 'Proposta', order_index: 4 },
]

function reindex(stages: StageField[]): StageField[] {
  return stages.map((s, i) => ({ ...s, order_index: i + 1 }))
}

export default function CreateJobPage() {
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const [company, setCompany] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [location, setLocation] = useState('')
  const [salaryMin, setSalaryMin] = useState('')
  const [salaryMax, setSalaryMax] = useState('')
  const [stages, setStages] = useState<StageField[]>(defaultStages)

  function addStage() {
    setStages((prev) => reindex([...prev, { name: '', order_index: 0 }]))
  }

  function removeStage(index: number) {
    setStages((prev) => reindex(prev.filter((_, i) => i !== index)))
  }

  function moveStage(index: number, direction: 'up' | 'down') {
    setStages((prev) => {
      const next = [...prev]
      const target = direction === 'up' ? index - 1 : index + 1
      if (target < 0 || target >= next.length) return prev
      ;[next[index], next[target]] = [next[target], next[index]]
      return reindex(next)
    })
  }

  function updateStageName(index: number, name: string) {
    setStages((prev) => prev.map((s, i) => (i === index ? { ...s, name } : s)))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')

    const min = parseFloat(salaryMin)
    const max = parseFloat(salaryMax)
    if (isNaN(min) || isNaN(max) || min > max) {
      setError('Intervalo salarial inválido.')
      return
    }

    setIsLoading(true)
    try {
      const job = await jobsService.create({
        company: company.trim() || undefined,
        title,
        description,
        location,
        salary_min: min,
        salary_max: max,
        stages: stages.filter((s) => s.name.trim()),
      })
      navigate(`/recruiter/jobs/${job.id}/pipeline`)
    } catch {
      setError('Erro ao criar vaga. Verifique os dados e tente novamente.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/recruiter/jobs">
            <ArrowLeft className="w-4 h-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Nova Vaga</h1>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Informações da Vaga</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {error && <p className="text-sm text-destructive">{error}</p>}

            <div className="space-y-2">
              <Label htmlFor="company">Empresa</Label>
              <Input
                id="company"
                placeholder="Ex: Acme Corp"
                value={company}
                onChange={(e) => setCompany(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="title">Título</Label>
              <Input
                id="title"
                required
                placeholder="Ex: Desenvolvedor Backend Sênior"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Descrição</Label>
              <Textarea
                id="description"
                required
                rows={5}
                placeholder="Descreva as responsabilidades, requisitos e benefícios..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="location">Localidade</Label>
              <Input
                id="location"
                required
                placeholder="Ex: São Paulo, SP (Remoto)"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="salary_min">Salário Mínimo (R$)</Label>
                <Input
                  id="salary_min"
                  type="number"
                  required
                  min={0}
                  step={1}
                  placeholder="5000"
                  value={salaryMin}
                  onChange={(e) => setSalaryMin(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="salary_max">Salário Máximo (R$)</Label>
                <Input
                  id="salary_max"
                  type="number"
                  required
                  min={0}
                  step={1}
                  placeholder="8000"
                  value={salaryMax}
                  onChange={(e) => setSalaryMax(e.target.value)}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Etapas do Processo</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {stages.map((stage, idx) => (
              <div key={idx} className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground w-5 text-right shrink-0">
                  {idx + 1}.
                </span>
                <Input
                  value={stage.name}
                  onChange={(e) => updateStageName(idx, e.target.value)}
                  placeholder={`Etapa ${idx + 1}`}
                  className="flex-1"
                />
                {/* Move up */}
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => moveStage(idx, 'up')}
                  disabled={idx === 0}
                  title="Mover para cima"
                >
                  <ArrowUp className="w-4 h-4 text-muted-foreground" />
                </Button>
                {/* Move down */}
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => moveStage(idx, 'down')}
                  disabled={idx === stages.length - 1}
                  title="Mover para baixo"
                >
                  <ArrowDown className="w-4 h-4 text-muted-foreground" />
                </Button>
                {/* Remove */}
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => removeStage(idx)}
                  disabled={stages.length <= 1}
                  title="Remover etapa"
                >
                  <Trash2 className="w-4 h-4 text-muted-foreground" />
                </Button>
              </div>
            ))}

            <Separator />

            <Button type="button" variant="outline" size="sm" onClick={addStage}>
              <Plus className="w-4 h-4 mr-2" />
              Adicionar Etapa
            </Button>
          </CardContent>
        </Card>

        <div className="flex justify-end gap-3">
          <Button type="button" variant="outline" asChild>
            <Link to="/recruiter/jobs">Cancelar</Link>
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Criando...' : 'Criar Vaga'}
          </Button>
        </div>
      </form>
    </div>
  )
}
