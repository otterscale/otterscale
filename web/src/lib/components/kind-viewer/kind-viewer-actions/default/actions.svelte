<!-- <script lang="ts">
	import Ellipsis from '@lucide/svelte/icons/ellipsis';

	import Button from '$lib/components/ui/button/button.svelte';
</script>

<Button size="icon" variant="ghost" aria-label="Edit item" disabled>
	<Ellipsis size={16} aria-hidden="true" />
</Button> -->

<script lang="ts">
	import { EllipsisIcon, PencilIcon, Trash2Icon } from '@lucide/svelte';

	import Button from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Item from '$lib/components/ui/item';

	import View from './view.svelte';

	let { schema, object }: { schema: any; object: any } = $props();

	let actionsOpen = $state(false);
</script>

<DropdownMenu.Root bind:open={actionsOpen}>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="flex justify-end">
				<Button size="icon" variant="ghost" class="shadow-none" aria-label="Edit item" {...props}>
					<EllipsisIcon size={16} aria-hidden="true" />
				</Button>
			</div>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Group>
			<DropdownMenu.Item
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<View {schema} {object} />
			</DropdownMenu.Item>
			<DropdownMenu.Item
				disabled
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<Item.Root class="p-0 text-xs" size="sm">
					<Item.Media>
						<PencilIcon />
					</Item.Media>
					<Item.Content>
						<Item.Title>Update</Item.Title>
					</Item.Content>
				</Item.Root>
			</DropdownMenu.Item>
			<DropdownMenu.Item
				disabled
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<Item.Root class="p-0 text-xs **:text-destructive" size="sm">
					<Item.Media>
						<Trash2Icon />
					</Item.Media>
					<Item.Content>
						<Item.Title>Delete</Item.Title>
					</Item.Content>
				</Item.Root>
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
