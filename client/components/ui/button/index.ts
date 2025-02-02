import { cva, type VariantProps } from 'class-variance-authority'

export { default as Button } from './Button.vue'

export const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground font-bold rounded-none hover:bg-primary/90 transition-colors duration-200',
        destructive:
          'bg-red-700 text-white font-bold hover:bg-red-900 hover:text-red-100 transition-colors duration-200',
        outline:
          'bg-transparent text-black font-bold over:bg-black hover:text-white transition-colors duration-200',
        secondary:
          'bg-yellow-500 text-black font-bold hover:bg-yellow-700 hover:text-yellow-100 transition-colors duration-200',
        ghost:
          'bg-transparent text-primary font-bold rounded-none hover:bg-primary hover:text-primary-foreground transition-colors duration-200',
        link:
          'text-blue-600 underline decoration-4 underline-offset-4 hover:text-blue-800 hover:decoration-blue-800 transition-colors duration-200',
        // default:
        //   'bg-primary text-primary-foreground shadow hover:bg-primary/90',
        // destructive:
        //   'bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90',
        // outline:
        //   'border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground',
        // secondary:
        //   'bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80',
        // ghost: 'hover:bg-accent hover:text-accent-foreground',
        // link: 'text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-9 px-4 py-2',
        sm: 'h-8 rounded-md px-3 text-xs',
        lg: 'h-10 rounded-md px-8',
        icon: 'h-9 w-9',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export type ButtonVariants = VariantProps<typeof buttonVariants>
