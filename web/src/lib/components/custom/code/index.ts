import { tv, type VariantProps } from 'tailwind-variants';

import CopyButton from './code-copy-button.svelte';
import Root from './code.svelte';
import type { CodeCopyButtonProps, CodeRootProps } from './types';

export const codeVariants = tv({
	base: 'not-prose relative h-full overflow-auto rounded-lg border',
	variants: {
		variant: {
			default: 'border-border bg-card',
			secondary: 'bg-secondary/50 border-transparent'
		}
	}
});

export type CodeVariant = VariantProps<typeof codeVariants>['variant'];

export {
	Root,
	CopyButton,
	type CodeRootProps as RootProps,
	type CodeCopyButtonProps as CopyButtonProps
};
