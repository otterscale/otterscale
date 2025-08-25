<script lang="ts">
	import { page } from '$app/state';
	import { TagService, type Tag } from '$lib/api/tag/v1/tag_pb';
	import * as Table from '$lib/components/custom/table';
	import CreateTag from '$lib/components/settings/general/create-tag.svelte';
	import DeleteTag from '$lib/components/settings/general/delete-tag.svelte';
	import * as Layout from '$lib/components/settings/general/layout';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.machine_tag(), url: '' }
	});

	const transport: Transport = getContext('transport');
	const tagClient = createClient(TagService, transport);

	const tags = writable<Tag[]>();
	let isTagLoading = $state(true);

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await tagClient.listTags({}).then((response) => {
				tags.set(response.tags);
				isTagLoading = false;
			});
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isTagLoading}
	<Layout.Title>Tag</Layout.Title>

	<Layout.Description>
		Tags are identifiable labels that can be assigned to machines for filtering and group
		management. These tags help in organizing and managing machines based on their characteristics
		or purposes.
	</Layout.Description>
	<Layout.Actions>
		<CreateTag {tags} />
	</Layout.Actions>
	<Layout.Controller>
		<div class="rounded-lg border shadow-sm">
			<Table.Root>
				<Table.Header>
					<Table.Row
						class="*:bg-muted *:rounded-t-lg *:px-4 *:first:rounded-tl-lg *:last:rounded-tr-lg"
					>
						<Table.Head>TAG</Table.Head>
						<Table.Head>COMMENT</Table.Head>
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
									<DeleteTag {tag} {tags} />
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</Layout.Controller>
{/if}
