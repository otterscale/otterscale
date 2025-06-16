<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	import ActionEdit from './action-edit.svelte';
	import ActionAuthorize from './action-authorize.svelte';
	import ActionDelete from './action-delete.svelte';
	import { ActionSubvolume } from './action-subvolume';
	import type { FileSystem } from './types';

	let { fileSystem }: { fileSystem: FileSystem } = $props();
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
			<DropdownMenu.Label class="label">Views</DropdownMenu.Label>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionSubvolume />
			</DropdownMenu.Item>
			<DropdownMenu.Label class="label">Actions</DropdownMenu.Label>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionEdit {fileSystem} />
			</DropdownMenu.Item>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionAuthorize {fileSystem} />
			</DropdownMenu.Item>
			<DropdownMenu.Item class="action" onSelect={(e) => e.preventDefault()}>
				<ActionDelete {fileSystem} />
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
