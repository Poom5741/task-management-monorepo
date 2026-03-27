import * as React from 'react'
import { cn } from '../../lib/utils'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, label, error, id, ...props }, ref) => {
    const inputId = id || React.useId()
    
    return React.createElement(
      'div',
      { className: 'w-full' },
      label && React.createElement(
        'label',
        { htmlFor: inputId, className: 'block text-sm font-medium text-gray-700 mb-1' },
        label
      ),
      React.createElement('input', {
        ref,
        id: inputId,
        className: cn(
          'w-full px-3 py-2 border rounded-lg shadow-sm transition-colors',
          'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
          error ? 'border-red-500' : 'border-gray-300',
          'disabled:bg-gray-100 disabled:cursor-not-allowed',
          className
        ),
        ...props,
      }),
      error && React.createElement(
        'p',
        { className: 'mt-1 text-sm text-red-600' },
        error
      )
    )
  }
)

Input.displayName = 'Input'

export { Input }