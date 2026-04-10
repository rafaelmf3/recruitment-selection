import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { JobStage } from '@/types'

interface PipelineBarProps {
  stages: JobStage[]
  currentStageId?: string
  /** When true, renders a more compact version without numbers */
  compact?: boolean
}

export default function PipelineBar({ stages, currentStageId, compact }: PipelineBarProps) {
  if (!stages || stages.length === 0) return null

  const sorted = [...stages].sort((a, b) => a.order_index - b.order_index)
  const currentIndex = currentStageId
    ? sorted.findIndex((s) => s.id === currentStageId)
    : -1

  return (
    <div className="flex items-center gap-0 overflow-x-auto pb-1">
      {sorted.map((stage, idx) => {
        const isCompleted = currentIndex >= 0 && idx < currentIndex
        const isCurrent = idx === currentIndex
        const isPending = currentIndex === -1 || idx > currentIndex

        return (
          <div key={stage.id} className="flex items-center shrink-0">
            {/* Stage pill */}
            <div
              className={cn(
                'flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap transition-colors',
                isCompleted && 'bg-green-500 text-white',
                isCurrent && 'bg-primary text-primary-foreground ring-2 ring-primary/30',
                isPending && 'bg-muted text-muted-foreground',
              )}
            >
              {isCompleted && !compact && <Check className="w-3 h-3 shrink-0" />}
              {!compact && (
                <span className={cn(
                  'inline-flex items-center justify-center w-4 h-4 rounded-full text-[10px] font-bold mr-0.5',
                  isCompleted && 'bg-white/20',
                  isCurrent && 'bg-white/20',
                  isPending && 'bg-foreground/10',
                )}>
                  {idx + 1}
                </span>
              )}
              {stage.name}
            </div>

            {/* Connector line */}
            {idx < sorted.length - 1 && (
              <div
                className={cn(
                  'h-0.5 w-5 shrink-0',
                  isCompleted ? 'bg-green-500' : 'bg-muted'
                )}
              />
            )}
          </div>
        )
      })}
    </div>
  )
}
