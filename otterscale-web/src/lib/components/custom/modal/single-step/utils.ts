import { buttonVariants } from '$lib/components/ui/button';
import { cn } from '$lib/utils';
import { tv } from 'tailwind-variants';

const triggerVariants = tv({
    base: 'flex items-center gap-1 h-full w-full disabled:pointer-events-auto disabled:cursor-not-allowed disabled:text-muted-foreground',
    variants: {
        variant: {
            default: cn(buttonVariants({ variant: 'default', size: 'default' }), 'w-fit'),
            creative: '',
            destructive: 'text-destructive **:text-destructive'
        }
    },
    defaultVariants: {
        variant: 'default'
    }
});

export { triggerVariants }
