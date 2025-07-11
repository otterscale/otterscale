<script lang="ts" module>
	import { cn } from '$lib/utils.js';

	const content = 'No Data';
</script>

<script lang="ts">
	let { class: className }: { class?: string } = $props();
</script>

<p class={cn('flex h-full w-full items-center justify-center', className)}>
	{content}
</p>
