<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import Actions from './cell-actions.svelte';
	import Create from './create.svelte';
	import Import from './import.svelte';
	import ReadArchitectures from './read-architectures.svelte';

	import { ConfigurationService, type Configuration } from '$lib/api/configuration/v1/configuration_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);

	const reloadManager = new ReloadManager(() => {
		configurationClient.getConfiguration({}).then((response) => {
			configuration.set(response);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(async () => {
		try {
			await configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		reloadManager.start();
	});

	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if !isConfigurationLoading}
	<Layout.Root>
		<Layout.Title>{m.image()}</Layout.Title>
		<Layout.Description>
			{m.setting_boot_image_description()}
		</Layout.Description>
		<Layout.Controller>
			<Create {configuration} />
			<Import {configuration} />
		</Layout.Controller>
		<Layout.Viewer>
			<div class="rounded-lg border shadow-sm">
				<Table.Root>
					<Table.Header>
						<Table.Row class="[&_th]:bg-muted *:px-4 [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg">
							<Table.Head>{m.name()}</Table.Head>
							<Table.Head>{m.source()}</Table.Head>
							<Table.Head>{m.distro_series()}</Table.Head>
							<Table.Head>{m.default_value()}</Table.Head>
							<Table.Head class="text-right">{m.architecture()}</Table.Head>
							<Table.Head class="text-right">{m.status()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $configuration.bootImages as bootImage}
							<Table.Row class="*:px-4">
								<Table.Cell>{bootImage.name}</Table.Cell>
								<Table.Cell>{bootImage.source}</Table.Cell>
								<Table.Cell>
									<Badge variant="outline">{bootImage.distroSeries}</Badge>
								</Table.Cell>
								<Table.Cell>
									<Icon icon={bootImage.default ? 'ph:circle' : 'ph:x'} />
								</Table.Cell>
								<Table.Cell>
									<span class="mr-1 flex items-center justify-end gap-1">
										{#each bootImage.architectures as architecture}
											<Badge variant="outline" class="">{architecture}</Badge>
										{/each}
									</span>
								</Table.Cell>
								<Table.Cell>
									<span class="flex items-center justify-end">
										{Object.keys(bootImage.architectureStatusMap).length}
										<ReadArchitectures {bootImage} />
									</span>
								</Table.Cell>
								<Table.Cell class="p-0">
									<Actions {bootImage} {configuration} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
