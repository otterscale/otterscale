import { tv, type VariantProps } from 'tailwind-variants';

export const typographyVariants = tv({
	variants: {
		variant: {
			h1: 'scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl',
			h2: 'scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0',
			h3: 'scroll-m-20 text-2xl font-semibold tracking-tight',
			h4: 'scroll-m-20 text-xl font-semibold tracking-tight',
			h5: 'scroll-m-20 text-lg font-semibold tracking-tight',
			h6: 'scroll-m-20 text-sm font-semibold tracking-tight',

			p: 'leading-7 [&:not(:first-child)]:mt-6',
			lead: 'text-xl text-muted-foreground',
			large: 'text-lg font-semibold',
			small: 'text-sm font-medium leading-none',
			muted: 'text-sm text-muted-foreground',

			inline_code: 'rounded bg-muted px-2 py-1 font-mono text-sm',
			pre: 'overflow-x-auto rounded-lg bg-muted p-4 font-mono text-sm',

			blockquote: 'mt-6 border-l-2 border-muted-foreground pl-6 italic text-muted-foreground',
			caption: 'text-xs text-muted-foreground',

			ul: 'my-6 ml-6 list-disc [&>li]:mt-2',
			ol: 'my-6 ml-6 list-decimal [&>li]:mt-2'
		}
	}
});

export type TypographyVariant = VariantProps<typeof typographyVariants>['variant'];
