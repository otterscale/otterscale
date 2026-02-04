<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, attributes }: Props = $props();
</script>

<div class="flex w-full items-center justify-between" {...attributes}>
	{@render children()}
</div>
