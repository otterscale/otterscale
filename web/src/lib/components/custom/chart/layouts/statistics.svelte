<script lang="ts" module>
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let { ref = $bindable(null), children }: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} =
		$props();
</script>

<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
	{@render children?.()}
</div>
