import * as React from 'react'
import { cn } from '../../lib/utils'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
  size?: 'sm' | 'md' | 'lg'
}

function Modal({ isOpen, onClose, title, children, size = 'md' }: ModalProps) {
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
  }

  return React.createElement(
    'div',
    { className: 'fixed inset-0 z-50 flex items-center justify-center' },
    React.createElement('div', {
      className: 'fixed inset-0 bg-black bg-opacity-50',
      onClick: onClose,
      'aria-hidden': 'true',
    }),
    React.createElement(
      'div',
      {
        className: cn(
          'relative bg-white rounded-lg shadow-xl w-full mx-4',
          sizes[size]
        ),
        role: 'dialog',
        'aria-modal': 'true',
        'aria-labelledby': 'modal-title',
      },
      React.createElement(
        'div',
        { className: 'flex items-center justify-between p-4 border-b' },
        React.createElement(
          'h2',
          { id: 'modal-title', className: 'text-lg font-semibold text-gray-900' },
          title
        ),
        React.createElement('button', {
          onClick: onClose,
          className: 'text-gray-400 hover:text-gray-600 transition-colors',
          'aria-label': 'Close modal',
        }, '✕')
      ),
      React.createElement('div', { className: 'p-4' }, children)
    )
  )
}

export { Modal }