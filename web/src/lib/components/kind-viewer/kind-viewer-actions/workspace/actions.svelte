<script lang="ts">
	import Ellipsis from '@lucide/svelte/icons/ellipsis';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';

	import Button from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	import Delete from './delete.svelte';
	import Update from './update.svelte';
	import View from './view.svelte';

	let { schema, object }: { schema: any; object: TenantOtterscaleIoV1Alpha1Workspace } = $props();

	let actionsOpen = $state(false);
</script>

<DropdownMenu.Root bind:open={actionsOpen}>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="flex justify-end">
				<Button size="icon" variant="ghost" class="shadow-none" aria-label="Actions" {...props}>
					<Ellipsis size={16} aria-hidden="true" />
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
				<View
					{schema}
					{object}
					onOpenChangeComplete={() => {
						actionsOpen = false;
					}}
				/>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			<DropdownMenu.Item
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<Update
					{schema}
					{object}
					onOpenChangeComplete={() => {
						actionsOpen = false;
					}}
				/>
			</DropdownMenu.Item>
			<DropdownMenu.Item
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<Delete
					{schema}
					onOpenChangeComplete={() => {
						actionsOpen = false;
					}}
				/>
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
