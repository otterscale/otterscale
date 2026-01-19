<script lang="ts">
	import Plus from '@lucide/svelte/icons/plus';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';

	let { resource }: { resource: string } = $props();

	let createModalOpen = $state(false);
</script>

<AlertDialog.Root bind:open={createModalOpen}>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'ghost' })}>
		<Plus class="opacity-60" size={16} />
		Create
	</AlertDialog.Trigger>
	<AlertDialog.Content class="max-w-xl">
		<AlertDialog.Header>
			<AlertDialog.Title>Create Resource "{resource}"</AlertDialog.Title>
		</AlertDialog.Header>

		<AlertDialog.Footer>
			<AlertDialog.Cancel class={buttonVariants({ variant: 'destructive' })}>
				Cancel
			</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					createModalOpen = false;
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
