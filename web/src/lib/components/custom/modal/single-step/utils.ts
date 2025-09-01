import { buttonVariants } from '$lib/components/ui/button';
import { cn } from '$lib/utils';
import { tv } from 'tailwind-variants';

const triggerVariants = tv({
	base: 'disabled:text-muted-foreground flex h-full w-full items-center gap-1 disabled:pointer-events-auto disabled:cursor-not-allowed',
	variants: {
		variant: {
			default: cn(buttonVariants({ variant: 'default', size: 'default' }), 'w-fit'),
			creative: '',
			destructive: 'text-destructive **:text-destructive',
		},
	},
	defaultVariants: {
		variant: 'default',
	},
});

export { triggerVariants };
