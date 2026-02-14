import type { ButtonProps } from '$lib/components/ui/button';
import type { UseClipboard } from '$lib/hooks/use-clipboard.svelte';

type CopyProps = ButtonProps & {
	text: string;
	onCopy?: (status: UseClipboard['status']) => void;
};

export type { CopyProps };
