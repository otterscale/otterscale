<script lang="ts">
	import Plus from '@lucide/svelte/icons/plus';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';

	let {
		kind,
		onsuccess,
		...rest
	}: {
		cluster: string;
		namespace: string;
		group: string;
		version: string;
		kind: string;
		resource: string;
		onsuccess?: () => void;
	} = $props();

	// Suppress unused variable warning
	void rest;

	let createModalOpen = $state(false);
</script>

<AlertDialog.Root bind:open={createModalOpen}>
	<AlertDialog.Trigger disabled class={buttonVariants({ variant: 'outline' })}>
		<Plus class="opacity-60" size={16} />
	</AlertDialog.Trigger>
	<AlertDialog.Content class="max-w-xl">
		<AlertDialog.Header>
			<AlertDialog.Title>Create {kind}</AlertDialog.Title>
		</AlertDialog.Header>

		<AlertDialog.Footer>
			<AlertDialog.Cancel class={buttonVariants({ variant: 'destructive' })}>
				Cancel
			</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					createModalOpen = false;
					onsuccess?.();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
