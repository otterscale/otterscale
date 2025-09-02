<script lang="ts" module>
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		title,
		subtitle,
		class: className,
	}: {
		title: string;
		subtitle?: string;
		class?: string;
	} = $props();
</script>

<div class={cn('flex w-full items-start gap-2 text-sm', className)}>
	<div class="grid gap-2">
		<div class="flex items-center gap-2 leading-none font-medium">
			{title}
		</div>
		{#if subtitle}
			<div class="text-muted-foreground flex items-center gap-2 leading-none">
				{subtitle}
			</div>
		{/if}
	</div>
</div>
