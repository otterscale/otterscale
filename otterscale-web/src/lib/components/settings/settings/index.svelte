<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import { TagService, type Tag } from '$lib/api/tag/v1/tag_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import UpdateNTPServer from './update-ntp-server.svelte';
	import UpdatePackageRepository from './update-package-repository.svelte';
	import CreateTag from './create-tag.svelte';
	import DeleteTag from './delete-tag.svelte';
	import * as Alert from '$lib/components/ui/alert';
	import ReadArchitectures from './read-architectures.svelte';
	import SetBootImageAsDefault from './set-boot-image-as-default.svelte';
	import CreateBootImage from './create-boot-image.svelte';
	import ImportBootImage from './import-boot-image.svelte';
	import Items from './utils/items.svelte';
	import Item from './utils/item.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const tagClient = createClient(TagService, transport);
	const configurationClient = createClient(ConfigurationService, transport);

	let configuration = $state(writable<Configuration>());
	let isConfigurationLoading = $state(true);
	async function fetchConfiguration() {
		try {
			configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let tags = $state(writable<Tag[]>());
	let isTagLoading = $state(true);

	async function fetchTags() {
		try {
			tagClient.listTags({}).then((response) => {
				tags.set(response.tags);
				isTagLoading = false;
			});
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchConfiguration();
			await fetchTags();
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main class="mx-auto w-7xl space-y-6">
	<Alert.Root variant="destructive" class="bg-destructive/5 border-destructive/10">
		<Icon icon="ph:warning" />
		<Alert.Description>
			These settings will apply across all scopes in your infrastructure. Please review changes
			carefully before applying them.</Alert.Description
		>
	</Alert.Root>
	{#if isMounted}
		<Items>
			<Item icon="ph:network" name="NTP Servers" type="network" value="ntp_servers">
				{#if !isConfigurationLoading}
					<span class="flex justify-between gap-2">
						<h1 class="text-lg font-bold">Address</h1>
						<UpdateNTPServer bind:configuration />
					</span>
					<p class="text-muted-foreground text-sm">
						NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces,
						to be used as time references for MAAS itself, the machines MAAS deploys, and devices
						that make use of MAAS's DHCP services.
					</p>
					<div class="bg-muted rounded-lg p-4">
						{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
							{#each $configuration.ntpServer.addresses as address}
								<Badge>{address}</Badge>
							{/each}
						{/if}
					</div>
				{/if}
			</Item>

			<Item icon="ph:cloud" name="Package Repository" type="system" value="package_repository">
				{#if !isConfigurationLoading}
					<h1 class="text-lg font-bold">Repository</h1>
					<p class="text-muted-foreground text-sm">
						Package repositories contain software packages that can be installed on machines. These
						repositories can be configured with custom URLs and enabled/disabled as needed to
						control software sources and updates for your infrastructure.
					</p>
					<div>
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head>NAME</Table.Head>
									<Table.Head>URL</Table.Head>
									<Table.Head>ENABLED</Table.Head>
									<Table.Head></Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each $configuration.packageRepositories as packageRepository}
									<Table.Row class="border-none *:text-xs">
										<Table.Cell>{packageRepository.name}</Table.Cell>
										<Table.Cell>
											<span class="flex items-center gap-1">
												{packageRepository.url}
												<Button
													variant="ghost"
													target="_blank"
													href={packageRepository.url}
													class="flex items-start gap-1"
												>
													<Icon icon="ph:arrow-square-out" />
												</Button>
											</span>
										</Table.Cell>
										<Table.Cell>
											<Icon icon={packageRepository.enabled ? 'ph:circle' : 'ph:x'} />
										</Table.Cell>
										<Table.Cell>
											<UpdatePackageRepository bind:configuration {packageRepository} />
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				{/if}
			</Item>

			<Item icon="ph:squares-four" name="Boot Image" type="system" value="boot_image">
				{#if !isConfigurationLoading}
					<span class="flex justify-between gap-2">
						<h1 class="text-lg font-bold">Image</h1>
						<span>
							<CreateBootImage bind:configuration />
							<ImportBootImage bind:configuration />
						</span>
					</span>

					<p class="text-muted-foreground text-sm">
						Boot images are operating system files used for machine deployment. These images can be
						configured with different architectures and distribution series, and can be set as
						default for automatic deployment of machines.
					</p>

					<div>
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head>NAME</Table.Head>
									<Table.Head>SOURCE</Table.Head>
									<Table.Head>DISTRO SERIES</Table.Head>
									<Table.Head>DEFAULT</Table.Head>
									<Table.Head class="text-right">ARCHITECTURES</Table.Head>
									<Table.Head></Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each $configuration.bootImages as bootImage}
									<Table.Row>
										<Table.Cell>{bootImage.name}</Table.Cell>
										<Table.Cell>{bootImage.source}</Table.Cell>
										<Table.Cell>
											<Badge variant="outline">{bootImage.distroSeries}</Badge>
										</Table.Cell>
										<Table.Cell>
											<Icon icon={bootImage.default ? 'ph:circle' : 'ph:x'} />
										</Table.Cell>
										<Table.Cell>
											<span class="flex items-center justify-end gap-1">
												<ReadArchitectures {bootImage} />
												{Object.keys(bootImage.architectureStatusMap).length}
											</span>
										</Table.Cell>
										<Table.Cell>
											<span class="flex justify-end">
												<SetBootImageAsDefault {bootImage} bind:configuration />
											</span>
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				{/if}
			</Item>

			<Item icon="ph:tag" name="Tag" type="machine" value="tag">
				{#if !isConfigurationLoading}
					<span class="flex justify-between gap-2">
						<h1 class="text-lg font-bold">Tag</h1>
						<CreateTag bind:tags />
					</span>
					<p class="text-muted-foreground text-sm">
						Tags are identifiable labels that can be assigned to machines for filtering and group
						management. These tags help in organizing and managing machines based on their
						characteristics or purposes.
					</p>
					<div>
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head>TAG</Table.Head>
									<Table.Head>COMMENT</Table.Head>
									<Table.Head></Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each $tags as tag}
									<Table.Row>
										<Table.Cell>{tag.name}</Table.Cell>
										<Table.Cell>
											<p class={cn(tag.comment ? 'text-primary' : 'text-muted-foreground')}>
												{tag.comment ? tag.comment : 'No comments available.'}
											</p>
										</Table.Cell>
										<Table.Cell>
											<DeleteTag {tag} bind:tags />
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				{/if}
			</Item>
		</Items>
	{/if}
</main>
