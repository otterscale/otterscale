<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { Nexus, type Facility_Action } from '$gen/api/nexus/v1/nexus_pb';
	import { DoAction } from './index';

	let {
		scopeUuid,
		facilityName
	}: {
		scopeUuid: string;
		facilityName: string;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const actionsStore = writable<Facility_Action[]>([]);
	const actionsLoading = writable(true);
	async function fetchActions() {
		try {
			const response = await client.listActions({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			});
			actionsStore.set(response.actions);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			actionsLoading.set(false);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchActions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});
</script>

{#if isMounted}
	{#if $actionsStore.length > 0}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				<Button variant="ghost" class="flex items-center gap-1">
					<Icon icon="tabler:settings" class="size-4" />
					Actions
				</Button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content>
				{#each $actionsStore as action}
					<DropdownMenu.Item>
						<span class="flex h-full w-full items-center justify-start gap-2">
							<DoAction {action} />
							<span class="flex h-full w-full items-center justify-between gap-2">
								<p class="text-sm">{action.name}</p>
								<HoverCard.Root openDelay={13}>
									<HoverCard.Trigger>
										<Icon icon="ph:info" class="size-4 text-blue-800" />
									</HoverCard.Trigger>
									<HoverCard.Content class="text-sm">
										{action.description}
									</HoverCard.Content>
								</HoverCard.Root>
							</span>
						</span>
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{/if}
{:else}
	<ComponentLoading />
{/if}
