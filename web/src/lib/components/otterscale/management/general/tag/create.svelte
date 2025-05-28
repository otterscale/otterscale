<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { TagService, type Tag, type CreateTagRequest } from '$gen/api/tag/v1/tag_pb';
	import { toast } from 'svelte-sonner';

	let { tags = $bindable() }: { tags: Tag[] } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(TagService, transport);

	const DEFAULT_REQUEST = {} as CreateTagRequest;

	let createTagRequest = $state(DEFAULT_REQUEST);

	function reset() {
		createTagRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="ghost">
			<Icon icon="ph:plus" /> Tag
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add Tag</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="flex flex-col gap-4 rounded-lg border p-4">
					<div class="grid gap-2">
						<Label>Name</Label>
						<Input bind:value={createTagRequest.name} />
					</div>
					<div class="grid gap-2">
						<Label>Comment</Label>
						<Input bind:value={createTagRequest.comment} />
					</div>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.createTag(createTagRequest).then((r) => {
						toast.success(`Create ${createTagRequest.name}`);
						client.listTags({}).then((r) => {
							tags = r.tags;
						});
					});
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
