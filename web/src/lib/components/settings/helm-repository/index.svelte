<script lang="ts" module>
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import Update from './update.svelte';
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
			{m.setting_html_repository_description()}
		</Layout.Description>
		<Layout.Controller>
			<Update {urls} />
		</Layout.Controller>
		<Layout.Viewer>
			<Card.Root>
				<Card.Content>
					{#if $urls}
						{#each $urls as url}
							<Badge>{url}</Badge>
						{/each}
					{/if}
				</Card.Content>
			</Card.Root>
		</Layout.Viewer>
	</Layout.Root>
{/if}
