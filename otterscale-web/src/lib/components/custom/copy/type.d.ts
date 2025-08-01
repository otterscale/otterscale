import type { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
import type { WithChildren, WithoutChildren } from 'bits-ui';
import type { Snippet } from 'svelte';
import type { HTMLAttributes } from 'svelte/elements';
import type { ButtonProps } from '$lib/components/ui/button';

type CopyProps = ButtonProps & {
    text: string;
    onCopy?: (status: UseClipboard['status']) => void;
}

export type { CopyProps };
