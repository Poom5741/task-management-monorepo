import * as React from 'react'
import { cn } from '../../lib/utils'

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string
}

const Skeleton: React.FC<SkeletonProps> = ({ className, ...props }) => {
  return (
    <div
      className={cn('skeleton', className)}
      {...props}
    />
  )
}

interface SkeletonCardProps {
  className?: string
}

const SkeletonCard: React.FC<SkeletonCardProps> = ({ className }) => {
  return (
    <div className={cn('card p-6 space-y-4', className)}>
      <div className="flex items-center space-x-4">
        <Skeleton className="h-12 w-12 rounded-full" />
        <div className="space-y-2 flex-1">
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-3 w-1/2" />
        </div>
      </div>
      <Skeleton className="h-20 w-full" />
      <div className="flex items-center justify-between">
        <Skeleton className="h-8 w-20" />
        <Skeleton className="h-8 w-24" />
      </div>
    </div>
  )
}

interface SkeletonTextProps {
  lines?: number
  className?: string
}

const SkeletonText: React.FC<SkeletonTextProps> = ({ lines = 3, className }) => {
  return (
    <div className={cn('space-y-2', className)}>
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton 
          key={i} 
          className={cn(
            'h-4',
            i === lines - 1 ? 'w-2/3' : 'w-full'
          )} 
        />
      ))}
    </div>
  )
}

export { Skeleton, SkeletonCard, SkeletonText }
