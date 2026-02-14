<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import Ellipsis from '@lucide/svelte/icons/ellipsis';
	import type { Row } from '@tanstack/table-core';

	import Deletor from '$lib/components/dynamic-form/simpleapp/delete-dialog.svelte';
	import Editor from '$lib/components/dynamic-form/simpleapp/edit-sheet.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	let { row, schema, object }: { row: Row<Record<string, JsonValue>>; schema: any; object: any } =
		$props();

	let actionsOpen = $state(false);
</script>

<DropdownMenu.Root bind:open={actionsOpen}>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="flex justify-end">
				<Button size="icon" variant="ghost" class="shadow-none" aria-label="Edit item" {...props}>
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
				<Editor
					name={String(row.original['Name'])}
					{schema}
					{object}
					onOpenChangeComplete={() => {
						if (actionsOpen) {
							actionsOpen = false;
						}
					}}
				/>
			</DropdownMenu.Item>
			<DropdownMenu.Item
				onSelect={(e) => {
					e.preventDefault();
				}}
			>
				<Deletor
					name={String(row.original['Name'])}
					onOpenChangeComplete={() => {
						if (actionsOpen) {
							actionsOpen = false;
						}
					}}
				/>
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
