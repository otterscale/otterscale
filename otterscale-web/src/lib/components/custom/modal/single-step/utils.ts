import { buttonVariants } from '$lib/components/ui/button';
import { cn } from '$lib/utils';
import { tv } from 'tailwind-variants';

const triggerVariants = tv({
    base: 'flex items-center gap-1 h-full w-full',
    variants: {
        variant: {
            default: cn(buttonVariants({ variant: 'default', size: 'default' }), 'w-fit'),
            creative: buttonVariants({ variant: 'ghost', size: 'sm' }),
            destructive: cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'text-destructive')
        }
    },
    defaultVariants: {
        variant: 'default'
    }
});

export { triggerVariants }