<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Nexus, type Tag, type DeleteTagRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';

	let {
		tag
	}: {
		tag: Tag;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { name: tag.name } as DeleteTagRequest;

	let deleteTagRequest = $state(DEFAULT_REQUEST);

	function reset() {
		deleteTagRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="ghost">
			<Icon icon="ph:trash" class="p-0" /> Delete
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete Tag {tag.name}</AlertDialog.Title>
			<AlertDialog.Description class="rounded-lg bg-muted/50 p-4">
				Are you sure you want to delete this tag? This action cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.deleteTag(deleteTagRequest).then((r) => {
						toast.info(`Delete ${deleteTagRequest.name}`);
					});
					// toast.info(`Delete ${deleteTagRequest.name}`);
					console.log(deleteTagRequest.name);
					reset();
					close();
				}}
			>
				Delete
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
