<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import { TagService, type Tag } from '$lib/api/tag/v1/tag_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Alert from '$lib/components/ui/alert';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import CreateBootImage from './create-boot-image.svelte';
	import CreateTag from './create-tag.svelte';
	import DeleteTag from './delete-tag.svelte';
	import ImportBootImage from './import-boot-image.svelte';
	import * as Layout from './layout';
	import ReadArchitectures from './read-architectures.svelte';
	import SetBootImageAsDefault from './set-boot-image-as-default.svelte';
	import SingleSignOn from './single-sign-on.svelte';
	import UpdateNTPServer from './update-ntp-server.svelte';
	import UpdatePackageRepository from './update-package-repository.svelte';
	import { Item, Items } from './utils';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	let configuration = $state(writable<Configuration>());
	let isConfigurationLoading = $state(true);

	const tagClient = createClient(TagService, transport);
	const configurationClient = createClient(ConfigurationService, transport);
	let tags = $state(writable<Tag[]>());
	let isTagLoading = $state(true);
	let isMounted = $state(false);

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
					<Layout.Title>Address</Layout.Title>
					<Layout.Description>
						NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces,
						to be used as time references for MAAS itself, the machines MAAS deploys, and devices
						that make use of MAAS's DHCP services.
					</Layout.Description>
					<Layout.Actions>
						<UpdateNTPServer bind:configuration />
					</Layout.Actions>
					<Layout.Controller>
						<Card.Root>
							<Card.Content>
								{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
									{#each $configuration.ntpServer.addresses as address}
										<Badge>{address}</Badge>
									{/each}
								{/if}
							</Card.Content>
						</Card.Root>
					</Layout.Controller>
				{/if}
			</Item>

			<Item icon="ph:cloud" name="Package Repository" type="system" value="package_repository">
				{#if !isConfigurationLoading}
					<Layout.Title>Repository</Layout.Title>
					<Layout.Description>
						Package repositories contain software packages that can be installed on machines. These
						repositories can be configured with custom URLs and enabled/disabled as needed to
						control software sources and updates for your infrastructure.
					</Layout.Description>
					<Layout.Controller>
						<div class="rounded-lg border shadow-sm">
							<Table.Root>
								<Table.Header class="bg-muted rounded-lg shadow-sm">
									<Table.Row class="*:px-4">
										<Table.Head>NAME</Table.Head>
										<Table.Head>URL</Table.Head>
										<Table.Head>ENABLED</Table.Head>
										<Table.Head></Table.Head>
									</Table.Row>
								</Table.Header>
								<Table.Body>
									{#each $configuration.packageRepositories as packageRepository}
										<Table.Row class="*:px-4">
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
												<div class="flex items-center justify-end">
													<UpdatePackageRepository bind:configuration {packageRepository} />
												</div>
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					</Layout.Controller>
				{/if}
			</Item>

			<Item icon="ph:squares-four" name="Boot Image" type="system" value="boot_image">
				{#if !isConfigurationLoading}
					<Layout.Title>Image</Layout.Title>
					<Layout.Description>
						Boot images are operating system files used for machine deployment. These images can be
						configured with different architectures and distribution series, and can be set as
						default for automatic deployment of machines.
					</Layout.Description>
					<Layout.Actions>
						<CreateBootImage bind:configuration />
						<ImportBootImage bind:configuration />
					</Layout.Actions>
					<Layout.Controller>
						<div class="rounded-lg border shadow-sm">
							<Table.Root>
								<Table.Header class="bg-muted rounded-lg">
									<Table.Row class="*:px-4">
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
												<span class="flex items-center justify-end gap-1">
													<ReadArchitectures {bootImage} />
													{Object.keys(bootImage.architectureStatusMap).length}
												</span>
											</Table.Cell>
											<Table.Cell>
												<div class="flex items-center justify-end">
													<SetBootImageAsDefault {bootImage} bind:configuration />
												</div>
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					</Layout.Controller>
				{/if}
			</Item>

			<Item icon="ph:tag" name="Tag" type="machine" value="tag">
				{#if !isConfigurationLoading}
					<Layout.Title>Tag</Layout.Title>

					<Layout.Description>
						Tags are identifiable labels that can be assigned to machines for filtering and group
						management. These tags help in organizing and managing machines based on their
						characteristics or purposes.
					</Layout.Description>
					<Layout.Actions>
						<CreateTag bind:tags />
					</Layout.Actions>
					<Layout.Controller>
						<div class="rounded-lg border shadow-sm">
							<Table.Root>
								<Table.Header class="bg-muted rounded-lg">
									<Table.Row class="*:px-4">
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
													<DeleteTag {tag} bind:tags />
												</div>
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					</Layout.Controller>
				{/if}
			</Item>

			<Item icon="ph:key" name="Single Sign On" type="secret" value="single-sign-on">
				{#if !isConfigurationLoading}
					<Layout.Title>Single Sign On</Layout.Title>
					<Layout.Description>
						Single Sign-On (SSO) allows users to access multiple applications with one set of
						credentials. Configure your SSO provider details here to enable centralized
						authentication across your infrastructure management system.
					</Layout.Description>
					<Layout.Controller>
						<SingleSignOn />
					</Layout.Controller>
				{/if}
			</Item>
		</Items>
	{/if}
</main>
