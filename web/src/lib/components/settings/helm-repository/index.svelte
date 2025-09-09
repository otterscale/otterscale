<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Update from './update.svelte';

	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Layout from '$lib/components/settings/layout';
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const environmentClient = createClient(EnvironmentService, transport);

	const urls = writable<string[]>([]);
	let isUrlsLoading = $state(true);

	onMount(async () => {
		try {
			await environmentClient.getConfigHelmRepositories({}).then((response) => {
				urls.set(response.urls);
				isUrlsLoading = false;
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isUrlsLoading}
	<Layout.Root>
		<Layout.Title>{m.repository()}</Layout.Title>
		<Layout.Description>
			{m.setting_helm_repository_description()}
		</Layout.Description>
		<Layout.Controller>
			<Update {urls} />
		</Layout.Controller>
		<Layout.Viewer>
			{#if $urls}
				<div class="grid grid-cols-4 gap-4">
					{#each $urls as url}
						<Card.Root class="group">
							<Card.Content class="flex size-fit items-center gap-2">
								<div class="bg-muted group-hover:bg-muted-foreground size-fit rounded-full p-2">
									<Icon icon="ph:cloud" class="text-muted-foreground group-hover:text-muted size-6" />
								</div>
								<div>
									<p class="text-base font-medium">Repository</p>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<p class="text-muted-foreground max-w-3xs truncate text-sm">
												{url}
											</p>
										</Tooltip.Trigger>
										<Tooltip.Content>
											{url}
										</Tooltip.Content>
									</Tooltip.Root>
								</div>
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
			{/if}
		</Layout.Viewer>
	</Layout.Root>
{/if}
