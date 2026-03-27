# TaskFlow Design System

## Overview
A modern, professional design system for task management SaaS with micro-interactions and accessible patterns.

## Color Palette

### Primary Colors
- **Primary**: `#0D9488` (Teal 600)
- **Primary Light**: `#14B8A6` (Teal 500)
- **Primary Dark**: `#0F766E` (Teal 700)
- **Background**: `#F0FDFA` (Teal 50)
- **Text Primary**: `#134E4A` (Teal 900)
- **Text Secondary**: `#115E59` (Teal 800)

### Accent Colors
- **CTA**: `#F97316` (Orange 500)
- **CTA Hover**: `#EA580C` (Orange 600)
- **Success**: `#10B981` (Emerald 500)
- **Warning**: `#F59E0B` (Amber 500)
- **Error**: `#EF4444` (Red 500)

### Neutral Colors
- **White**: `#FFFFFF`
- **Gray 50**: `#F9FAFB`
- **Gray 100**: `#F3F4F6`
- **Gray 200**: `#E5E7EB`
- **Gray 300**: `#D1D5DB`
- **Gray 400**: `#9CA3AF`
- **Gray 500**: `#6B7280`
- **Gray 600**: `#4B5563`
- **Gray 700**: `#374151`
- **Gray 800**: `#1F2937`
- **Gray 900`: `#111827`

## Typography

### Font Family
- **Primary**: Plus Jakarta Sans (300, 400, 500, 600, 700)
- **Import**: `https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700&display=swap`

### Type Scale
- **Heading 1**: 2.25rem (36px), font-weight: 700
- **Heading 2**: 1.875rem (30px), font-weight: 600
- **Heading 3**: 1.5rem (24px), font-weight: 600
- **Heading 4**: 1.25rem (20px), font-weight: 600
- **Body Large**: 1.125rem (18px), font-weight: 400
- **Body**: 1rem (16px), font-weight: 400
- **Body Small**: 0.875rem (14px), font-weight: 400
- **Caption**: 0.75rem (12px), font-weight: 400

## Spacing Scale
- **xs**: 0.25rem (4px)
- **sm**: 0.5rem (8px)
- **md**: 1rem (16px)
- **lg**: 1.5rem (24px)
- **xl**: 2rem (32px)
- **2xl**: 3rem (48px)
- **3xl**: 4rem (64px)
- **4xl**: 6rem (96px)

## Border Radius
- **sm**: 0.375rem (6px)
- **md**: 0.5rem (8px)
- **lg**: 0.75rem (12px)
- **xl**: 1rem (16px)
- **2xl**: 1.5rem (24px)
- **full**: 9999px

## Shadows
- **sm**: `0 1px 2px 0 rgba(0, 0, 0, 0.05)`
- **md**: `0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)`
- **lg**: `0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)`
- **xl**: `0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)`
- **card**: `0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)`
- **card-hover**: `0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)`

## Animation Guidelines

### Timing
- **Micro-interactions**: 150-200ms
- **Modal/Panel transitions**: 250-300ms
- **Page transitions**: 300-400ms

### Easing
- **Enter**: `cubic-bezier(0, 0, 0.2, 1)` (ease-out)
- **Exit**: `cubic-bezier(0.4, 0, 1, 1)` (ease-in)
- **Default**: `cubic-bezier(0.4, 0, 0.2, 1)` (ease-in-out)

### Accessibility
- Respect `prefers-reduced-motion`
- No infinite animations on decorative elements
- Loading states: use spinners or skeletons

## Component Patterns

### Buttons
- **Primary**: bg-teal-600, text-white, hover:bg-teal-700
- **Secondary**: bg-white, border-teal-600, text-teal-600
- **CTA**: bg-orange-500, text-white, hover:bg-orange-600
- **Danger**: bg-red-500, text-white, hover:bg-red-600
- **Ghost**: bg-transparent, hover:bg-gray-100

### Cards
- **Default**: bg-white, rounded-xl, shadow-card
- **Hover**: shadow-card-hover, transform translateY(-2px)
- **Border**: border border-gray-200

### Inputs
- **Default**: border-gray-300, focus:border-teal-500, focus:ring-teal-500
- **Error**: border-red-500, focus:ring-red-500
- **Disabled**: bg-gray-100, text-gray-500

### Badges
- **Active/Success**: bg-emerald-100, text-emerald-700
- **Pending**: bg-amber-100, text-amber-700
- **Inactive**: bg-gray-100, text-gray-700

## Layout Patterns
- **Container max-width**: 1280px (max-w-7xl)
- **Page padding**: 1.5rem (px-6)
- **Section spacing**: 3rem (space-y-12)
- **Card grid**: grid-cols-1 md:grid-cols-2 lg:grid-cols-3

## Accessibility
- Minimum contrast ratio: 4.5:1
- Focus indicators: visible ring on interactive elements
- Cursor: pointer on all clickable elements
- Reduced motion support
