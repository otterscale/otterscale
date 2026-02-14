<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { MachineService, type Tag } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import Create from './create.svelte';
	import Delete from './delete.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);
	const tags = writable<Tag[]>();

	async function fetch() {
		try {
			const response = await client.listTags({});
			tags.set(response.tags);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}

	let isTagLoading = $state(true);
	onMount(async () => {
		await fetch();
		isTagLoading = false;
	});
</script>

{#if !isTagLoading}
	<Layout.Root>
		<Layout.Title>{m.tags()}</Layout.Title>
		<Layout.Description>
			{m.setting_machine_tag_description()}
		</Layout.Description>
		<Layout.Controller>
			<Create {tags} />
		</Layout.Controller>
		<Layout.Viewer>
			<div class="w-full rounded-lg border shadow-sm">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.tag()}</Table.Head>
							<Table.Head>{m.comment()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $tags as tag (tag.name)}
							<Table.Row>
								<Table.Cell>{tag.name}</Table.Cell>
								<Table.Cell>
									<p class={cn(tag.comment ? 'text-primary' : 'text-muted-foreground')}>
										{tag.comment ? tag.comment : 'No comments available.'}
									</p>
								</Table.Cell>
								<Table.Cell>
									<div class="flex items-center justify-end">
										<Delete {tag} {tags} />
									</div>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
