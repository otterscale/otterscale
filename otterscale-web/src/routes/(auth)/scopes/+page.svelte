<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { goto } from '$app/navigation';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import SquareGridImage from '$lib/assets/square-grid.svg';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { scopeIcon } from '$lib/components/scopes/icon';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const scopes = writable<Scope[]>([]);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	onMount(fetchScopes);
</script>

<svelte:head>
	<title>{m.welcome_to({ name: 'OtterScale 🦦' })}</title>
</svelte:head>

<main
	class="bg-sidebar relative flex min-h-screen flex-col overflow-hidden px-2 py-20 md:px-4 md:py-24"
>
	<div class="absolute inset-x-0 top-0 flex h-full w-full items-center justify-center opacity-100">
		<img
			src={SquareGridImage}
			alt="square-grid"
			class="[mask-image:radial-gradient(75%_75%_at_center,white,transparent)] opacity-90"
		/>
	</div>

	<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.scope_selector()}</h2>
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.scope_selector_description()}
	</p>

	<div class="z-10 mx-auto grid w-full grid-cols-3 gap-4 px-4 py-10 sm:gap-6 xl:px-0 2xl:w-3/5">
		{#each $scopes as scope}
			<Card.Root
				class="group cursor-pointer"
				onclick={() => {
					goto(dynamicPaths.scope(scope.name).url);
				}}
			>
				<Card.Header class="gap-0">
					<div class="flex items-center gap-4">
						<div class="bg-primary flex size-10 items-center justify-center rounded-lg">
							<Icon
								icon="{scopeIcon($scopes.findIndex((s) => s.name === scope.name))}-fill"
								class="text-primary-foreground size-6"
							/>
						</div>
						<div class="grid -space-y-1">
							<Card.Description class="capitalize">
								{scope.status}
							</Card.Description>
							<Card.Title class="text-2xl text-nowrap">{scope.name}</Card.Title>
						</div>
					</div>

					<Card.Action class="overflow-hidden group-hover:self-center">
						<Badge variant="outline" class="hidden group-hover:hidden lg:block">
							<span class="text-green-600">
								{m.machines()}: {scope.machineCount > 0 ? scope.machineCount : '-'} /
								{m.unit()}: {scope.unitCount > 0 ? scope.unitCount : '-'}
							</span>
						</Badge>
						<div
							class="text-muted-foreground group-hover:bg-primary group-hover:text-primary-foreground hidden size-10 items-center justify-center rounded-full transition-colors group-hover:inline-flex"
						>
							<Icon icon="ph:arrow-right" class="size-6" />
						</div>
					</Card.Action>
				</Card.Header>
			</Card.Root>
		{/each}
	</div>
</main>
