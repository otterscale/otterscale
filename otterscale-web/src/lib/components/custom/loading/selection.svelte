<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { cn } from '$lib/utils';

	let { class: className }: { class?: string } = $props();
</script>

<Skeleton class={cn('bg-muted h-[30px]', className)} />
