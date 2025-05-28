<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import * as Select from '$lib/components/ui/select/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import {
		ConfigurationService,
		type Configuration,
		type Configuration_BootImageSelection,
		type CreateBootImageRequest
	} from '$gen/api/configuration/v1/configuration_pb';

	let { configuration = $bindable() }: { configuration: Configuration } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	const bootImageSelectionsStore = writable<Configuration_BootImageSelection[]>([]);
	const bootImageSelectionsLoading = writable(true);
	async function fetchBootImageSelections() {
		try {
			const response = await client.listBootImageSelections({});
			bootImageSelectionsStore.set(response.bootImageSelections);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			bootImageSelectionsLoading.set(false);
		}
	}

	const DEFAULT_ARCHITECTURES = [] as string[];
	const DEFAULT_REQUEST = { architectures: DEFAULT_ARCHITECTURES } as CreateBootImageRequest;

	let createBootImageRequest = $state(DEFAULT_REQUEST);
	let candidateArchitectures = $state(DEFAULT_ARCHITECTURES);

	function refreshCandidateArchitecturesByDistroSeries(image: Configuration_BootImageSelection) {
		candidateArchitectures = image.architectures;
	}

	function reset() {
		createBootImageRequest = DEFAULT_REQUEST;
		resetCandidateArchitectures();
	}
	function resetSelectedArchitectures() {
		createBootImageRequest.architectures = DEFAULT_ARCHITECTURES;
	}
	function resetCandidateArchitectures() {
		candidateArchitectures = DEFAULT_ARCHITECTURES;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchBootImageSelections();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="ghost">
			<span class="flex items-center gap-2">
				<Icon icon="ph:plus" />
				Boot Image
			</span>
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Create Boot Image</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-4">
				<div class="grid gap-2">
					<Label for="name">Distro Series</Label>
					<Select.Root type="single" bind:value={createBootImageRequest.distroSeries}>
						<Select.Trigger>
							{createBootImageRequest.distroSeries === ''
								? 'Select'
								: createBootImageRequest.distroSeries}
						</Select.Trigger>
						<Select.Content>
							{#each $bootImageSelectionsStore as image}
								<Select.Item
									value={image.distroSeries}
									onclick={() => {
										refreshCandidateArchitecturesByDistroSeries(image);
										resetSelectedArchitectures();
									}}
								>
									{image.name}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid gap-2">
					<Label for="source">Architecture</Label>
					<span class="flex flex-wrap items-center gap-1">
						{#each createBootImageRequest.architectures as architecture}
							<Badge
								variant="outline"
								onclick={() => {
									createBootImageRequest.architectures =
										createBootImageRequest.architectures.filter((a) => a !== architecture);
								}}
							>
								{architecture}
								<Icon icon="ph:x" class="p-0" />
							</Badge>
						{/each}
					</span>
					<Select.Root
						bind:value={createBootImageRequest.architectures}
						type="multiple"
						disabled={candidateArchitectures.length == 0}
					>
						<Select.Trigger>
							{candidateArchitectures.length > 0 ? 'Select' : 'Choose a Distro Series'}
						</Select.Trigger>
						<Select.Content>
							{#each candidateArchitectures as architecture}
								<Select.Item value={architecture}>
									{architecture}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.createBootImage(createBootImageRequest)
						.then((r) => {
							toast.success(
								`Create boot images ${createBootImageRequest.distroSeries}: ${createBootImageRequest.architectures.join(', ')}`
							);
							client.getConfiguration({}).then((r) => {
								configuration = r;
							});
						})
						.catch((e) => {
							`Create boot images fail`;
						});
					console.log(createBootImageRequest);
					reset();
					close();
				}}>Create</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
