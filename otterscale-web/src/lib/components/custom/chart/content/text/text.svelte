<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		value,
		unit,
		class: className,
		textClass
	}: {
		value?: string;
		unit?: string;
		class?: string;
		textClass?: string;
	} = $props();
</script>

<div class={cn('flex aspect-auto h-[250px] w-full items-center justify-center', className)}>
	<p class={cn('fill-foreground text-4xl! font-bold', textClass)}>
		{value || 'No Data'}
		{unit || ''}
	</p>
</div>
