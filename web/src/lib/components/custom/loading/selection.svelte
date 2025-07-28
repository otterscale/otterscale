<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { class: className }: { class?: string } = $props();
</script>

<Skeleton class={cn('bg-muted h-[20px]', className)} />
