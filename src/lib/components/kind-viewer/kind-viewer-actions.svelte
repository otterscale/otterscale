<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import Ellipsis from '@lucide/svelte/icons/ellipsis';
	import type { Row } from '@tanstack/table-core';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { cn } from '$lib/utils';

	let { row }: { row: Row<Record<string, JsonValue>> } = $props();

	let actionsOpen = $state(false);
	let deleteModalOpen = $state(false);
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
			<DropdownMenu.Item disabled>
				View
				<DropdownMenu.Shortcut>ctrl V</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
			<DropdownMenu.Item disabled>
				Edit
				<DropdownMenu.Shortcut>ctrl E</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<DropdownMenu.Item
			disabled
			onSelect={(e) => {
				e.preventDefault();
				deleteModalOpen = true;
			}}
		>
			<AlertDialog.Root
				bind:open={deleteModalOpen}
				onOpenChangeComplete={(isOpen) => {
					if (!isOpen) {
						actionsOpen = false;
					}
				}}
			>
				<AlertDialog.Trigger class={cn('text-destructive focus:text-destructive')}>
					Delete
				</AlertDialog.Trigger>
				<AlertDialog.Content class="max-w-xl">
					<AlertDialog.Header>
						<AlertDialog.Title>Delete</AlertDialog.Title>
						<AlertDialog.Description>{row.original.Name}</AlertDialog.Description>
					</AlertDialog.Header>
					<p>
						Are you sure you want to delete <span class="font-semibold">{row.original.Name}</span>?
						This action cannot be undone.
					</p>
					<AlertDialog.Footer>
						<AlertDialog.Cancel class={buttonVariants({ variant: 'destructive' })}>
							Cancel
						</AlertDialog.Cancel>
						<AlertDialog.Action
							onclick={() => {
								deleteModalOpen = false;
							}}
						>
							Confirm
						</AlertDialog.Action>
					</AlertDialog.Footer>
				</AlertDialog.Content>
			</AlertDialog.Root>
			<DropdownMenu.Shortcut>ctrl D</DropdownMenu.Shortcut>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
