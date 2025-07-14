<script lang="ts" module>
	import * as Card from '$lib/components/ui/card';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let {
		icon,
		title,
		description,
		action
	}: { icon: Snippet; title: Snippet; description: Snippet; action: Snippet } = $props();
</script>

<Card.Root>
	<Card.Content class="flex aspect-square flex-col items-center justify-between gap-4 border p-4">
		<div
			class="flex w-full flex-col items-center justify-center gap-4 p-4 text-3xl font-semibold leading-none tracking-tight"
		>
			{@render title()}
			<div class="text-muted-foreground 二p-4 w-full justify-start text-sm font-light">
				{@render description()}
			</div>
		</div>
		<div class="flex w-full items-center justify-center">
			<div class="w-48 p-4">
				{@render icon()}
			</div>
		</div>
		<div class="flex w-full items-center justify-center p-4">
			{@render action()}
		</div>
	</Card.Content>
</Card.Root>
