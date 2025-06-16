<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import Icon from '@iconify/svelte';
	import type { Namespace } from './types';
	import ActionEdit from './action-edit.svelte';
	import ActionCopy from './action-copy.svelte';
	import ActionDelete from './action-delete.svelte';

	let { namespace }: { namespace: Namespace } = $props();
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
			<DropdownMenu.Label class="label">Actions</DropdownMenu.Label>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionEdit {namespace} />
			</DropdownMenu.Item>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionCopy {namespace} />
			</DropdownMenu.Item>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionDelete {namespace} />
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
