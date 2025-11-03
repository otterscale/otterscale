<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
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
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);
	let expandedArchitectures = new SvelteMap<string, boolean>();

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
			<div class="rounded-lg border shadow-sm w-full">
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
									<div class="flex items-center justify-end gap-1">
										<div class="flex flex-wrap gap-1">
											{#if !expandedArchitectures.get(bootImage.name)}
												{#each bootImage.architectures.slice(0, 3) as architecture}
													<Badge variant="outline">{architecture}</Badge>
												{/each}
												{#if bootImage.architectures.length > 3}
													<Badge variant="outline" class="h-fit w-fit">
														+{bootImage.architectures.length - 3}
													</Badge>
												{/if}
											{:else}
												{#each bootImage.architectures as architecture}
													<Badge variant="outline">{architecture}</Badge>
												{/each}
											{/if}
										</div>
										{#if bootImage.architectures.length > 3}
											<Button
												variant="outline"
												size="icon"
												class="size-6"
												onclick={() => {
													expandedArchitectures.set(
														bootImage.name,
														!expandedArchitectures.get(bootImage.name),
													);
												}}
											>
												<Icon
													icon="ph:caret-left"
													class={cn(
														'size-4 transition-all',
														expandedArchitectures.get(bootImage.name)
															? 'rotate-90'
															: '-rotate-90',
													)}
												/>
											</Button>
										{/if}
									</div>
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
