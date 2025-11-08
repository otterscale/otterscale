<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Create from './create.svelte';
	import Delete from './delete.svelte';

	import { MachineService, type Tag } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);

	const tags = writable<Tag[]>();
	let isTagLoading = $state(true);

	onMount(async () => {
		try {
			await client.listTags({}).then((response) => {
				tags.set(response.tags);
				isTagLoading = false;
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
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
						<Table.Row
							class="*:px-4 [&_th]:bg-muted [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg"
						>
							<Table.Head>{m.tag()}</Table.Head>
							<Table.Head>{m.comment()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $tags as tag}
							<Table.Row class="*:px-4">
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
