<script lang="ts">
	import { useCodeCopyButton } from './code.svelte.js';
	import type { CodeCopyButtonProps } from './types.js';

	import { CopyButton } from '$lib/components/custom/copy-button';
	import { cn } from '$lib/utils';

	let {
		ref = $bindable(null),
		variant = 'ghost',
		size = 'icon',
		class: className,
		...rest
	}: CodeCopyButtonProps = $props();

	const copyButton = useCodeCopyButton();
</script>

<CopyButton class={cn('absolute top-2 right-2', className)} text={copyButton.code} {variant} {size} {...rest} />
