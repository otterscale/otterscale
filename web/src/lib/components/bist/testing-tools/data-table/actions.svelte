<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import Icon from '@iconify/svelte';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
	import type { FlexibleIOTest } from './types';

	let { flexibleIOTest }: { flexibleIOTest: FlexibleIOTest } = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="w-full">
				<Button variant="ghost" size="icon" class="float-right" {...props}>
					<Icon icon="ph:dots-three" />
				</Button>
			</div>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="[&_.action]:cursor-pointer [&_.action]:text-xs [&_.label]:text-xs">
		<DropdownMenu.Group>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<Edit {flexibleIOTest} />
			</DropdownMenu.Item>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<Delete {flexibleIOTest} />
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
