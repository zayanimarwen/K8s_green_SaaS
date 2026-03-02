import { gradeColor } from '@/lib/utils'

interface Props { grade: string; size?: 'sm' | 'md' | 'lg' }

export function GradeBadge({ grade, size = 'md' }: Props) {
  const sizes = { sm: 'text-xs px-2 py-0.5', md: 'text-sm px-3 py-1', lg: 'text-lg px-4 py-1.5' }
  return (
    <span
      className={`inline-block font-bold rounded-full ${sizes[size]}`}
      style={{ backgroundColor: gradeColor(grade) + '22', color: gradeColor(grade), border: `1px solid ${gradeColor(grade)}` }}
    >
      {grade}
    </span>
  )
}
