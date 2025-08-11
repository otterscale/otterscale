<script lang="ts" module>
	import * as Accordion from '$lib/components/ui/accordion';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let {
		icon,
		name,
		type,
		value,
		children
	}: { icon: string; name: string; type: string; value: string; children?: Snippet } = $props();
</script>

<Accordion.Item {value}>
	<Accordion.Trigger
		class="bg-muted rounded-none p-6 hover:no-underline data-[state=closed]:rounded-lg data-[state=open]:rounded-t-lg"
	>
		<div class="flex flex-1 items-center gap-4">
			<div class="relative hidden h-16 w-16 overflow-hidden rounded-md border sm:block">
				<Icon {icon} class="h-full w-full object-cover p-2" />
			</div>
			<div class="flex-1 space-y-1 text-left">
				<h3 class="text-xl font-semibold">{name}</h3>
				<Badge variant="outline" class="text-xs">{type}</Badge>
			</div>
		</div>
	</Accordion.Trigger>

	<Accordion.Content class="space-y-2 rounded-b-lg border border-t-0 p-6">
		{@render children?.()}
	</Accordion.Content>
</Accordion.Item>
