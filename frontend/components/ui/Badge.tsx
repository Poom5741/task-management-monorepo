import * as React from 'react'
import { cn } from '../../lib/utils'

interface BadgeProps {
  children: React.ReactNode
  variant?: 'success' | 'warning' | 'info' | 'error' | 'neutral'
  size?: 'sm' | 'md'
  className?: string
}

const Badge: React.FC<BadgeProps> = ({ 
  children, 
  variant = 'neutral', 
  size = 'md',
  className 
}) => {
  const variants = {
    success: 'badge-success',
    warning: 'badge-warning',
    info: 'badge-info',
    error: 'badge-error',
    neutral: 'badge-neutral',
  }
  
  const sizes = {
    sm: 'text-xs px-2 py-0.5',
    md: 'text-xs px-2.5 py-1',
  }

  return (
    <span className={cn(variants[variant], sizes[size], className)}>
      {children}
    </span>
  )
}

export { Badge }
