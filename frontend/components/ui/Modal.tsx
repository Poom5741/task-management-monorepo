import * as React from 'react'
import { cn } from '../../lib/utils'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  children: React.ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl'
  showCloseButton?: boolean
}

function Modal({ 
  isOpen, 
  onClose, 
  title, 
  description,
  children, 
  size = 'md',
  showCloseButton = true
}: ModalProps) {
  React.useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
      document.body.style.overflow = 'hidden'
    }

    return () => {
      document.removeEventListener('keydown', handleEscape)
      document.body.style.overflow = 'unset'
    }
  }, [isOpen, onClose])

  if (!isOpen) return null

  const sizes = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-2xl',
  }

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center animate-fade-in"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div 
        className="fixed inset-0 bg-gray-900/60 backdrop-blur-sm transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />
      <div 
        className={cn(
          'relative bg-white rounded-2xl shadow-xl w-full mx-4 animate-slide-up',
          sizes[size]
        )}
      >
        <div className="flex items-start justify-between p-6 border-b border-gray-100">
          <div>
            <h2 
              id="modal-title" 
              className="text-lg font-semibold text-gray-900"
            >
              {title}
            </h2>
            {description && (
              <p className="mt-1 text-sm text-gray-500">{description}</p>
            )}
          </div>
          {showCloseButton && (
            <button
              onClick={onClose}
              className="ml-4 text-gray-400 hover:text-gray-600 transition-colors rounded-lg p-1 hover:bg-gray-100"
              aria-label="Close modal"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  )
}

export { Modal }
