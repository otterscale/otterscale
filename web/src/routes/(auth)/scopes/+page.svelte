<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { resolve } from '$app/paths';
	import { type Scope, ScopeService } from '$lib/api/scope/v1/scope_pb';
	import SquareGridImage from '$lib/assets/square-grid.svg';
	import DialogCreateScope from '$lib/components/layout/dialog-create-scope.svelte';
	import { scopeIcon } from '$lib/components/scopes/icon';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	const EXCLUDED_SCOPES = ['cos', 'cos-dev', 'cos-lite'];

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	let scopes = $state<Scope[]>([]);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes = response.scopes.filter((scope) => !EXCLUDED_SCOPES.includes(scope.name));
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	function getCardColumnClass(index: number, scopeCount: number): string {
		const startColumns: Record<number, string> = {
			1: 'col-start-4',
			2: 'col-start-3',
			3: 'col-start-2'
		};
		return index === 0 ? startColumns[scopeCount] || '' : '';
	}

	function getScopeIndex(scopeName: string): number {
		return scopes.findIndex((scope) => scope.name === scopeName);
	}

	let open = $state(false);
	let mounted = $state(false);
	onMount(async () => {
		await fetchScopes();
		mounted = true;
	});
</script>

<svelte:head>
	<title>{m.welcome_to({ name: 'OtterScale' })}</title>
</svelte:head>

<main
	class="relative flex min-h-[calc(100vh-24*var(--spacing))] flex-col overflow-hidden px-2 py-20 md:px-4 md:py-24"
>
	<!-- Background Image -->
	<div class="absolute inset-x-0 top-0 flex h-full w-full items-center justify-center opacity-100">
		<img
			src={SquareGridImage}
			alt="square-grid"
			class="[mask-image:radial-gradient(75%_75%_at_center,white,transparent)] opacity-90"
		/>
	</div>

	<!-- Header -->
	{#if mounted}
		<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">
			{m.scope_selector()}
		</h2>
		<p class="mt-4 text-center text-lg text-muted-foreground">{m.scope_selector_description()}</p>
	{:else}
		<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">
			{m.scope_selector_loading()}
		</h2>
		<Button
			class="mt-4 text-center text-lg text-muted-foreground"
			variant="ghost"
			size="lg"
			disabled
		>
			<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
			{m.scope_selector_loading_description()}
		</Button>
	{/if}
	<!-- Scopes Grid -->
	<div class="z-10 mx-auto grid w-full grid-cols-8 gap-4 px-4 py-10 sm:gap-6 xl:px-0 2xl:w-3/4">
		{#if scopes.length === 0}
			<!-- Add Scope Card -->
			<button onclick={() => (open = true)} class="group col-span-2 col-start-4 cursor-pointer">
				<Card.Root class="transition-all duration-200 hover:scale-105 hover:shadow-lg">
					<Card.Header class="gap-0">
						<div class="flex items-center gap-4">
							<!-- Add Icon -->
							<div class="flex size-10 items-center justify-center rounded-lg bg-primary">
								<Icon icon="ph:plus-bold" class="size-6 text-card" />
							</div>

							<!-- Add Scope Info -->
							<div class="grid -space-y-1">
								<Card.Description class="capitalize">{m.create()}</Card.Description>
								<Card.Title class="text-2xl text-nowrap">New Scope</Card.Title>
							</div>
						</div>

						<!-- Action -->
						<Card.Action class="overflow-hidden group-hover:self-center">
							<div
								class="hidden size-10 items-center justify-center rounded-full text-muted-foreground transition-colors group-hover:inline-flex group-hover:bg-primary group-hover:text-primary-foreground"
							>
								<Icon icon="ph:arrow-right" class="size-6" />
							</div>
						</Card.Action>
					</Card.Header>
				</Card.Root>
			</button>
			<DialogCreateScope bind:open />
		{:else}
			{#each scopes as scope, index (scope.name)}
				<a
					href={resolve(`/(auth)/scope/[scope]`, { scope: scope.name })}
					class="group col-span-2 cursor-pointer {getCardColumnClass(index, scopes.length)}"
				>
					<Card.Root>
						<Card.Header class="gap-0">
							<div class="flex items-center gap-4">
								<!-- Scope Icon -->
								<div class="flex size-10 items-center justify-center rounded-lg bg-primary">
									<Icon
										icon="{scopeIcon(getScopeIndex(scope.name))}-fill"
										class="size-6 text-primary-foreground"
									/>
								</div>

								<!-- Scope Info -->
								<div class="grid -space-y-1">
									<Card.Description class="capitalize">
										{scope.status}
									</Card.Description>
									<Card.Title class="text-2xl text-nowrap">{scope.name}</Card.Title>
								</div>
							</div>

							<!-- Scope Stats and Action -->
							<Card.Action class="overflow-hidden group-hover:self-center">
								<Badge variant="outline" class="hidden group-hover:hidden lg:block">
									<span class="text-green-600">
										{m.machines()}: {scope.machineCount > 0 ? scope.machineCount : '-'} /
										{m.unit()}: {scope.unitCount > 0 ? scope.unitCount : '-'}
									</span>
								</Badge>
								<div
									class="hidden size-10 items-center justify-center rounded-full text-muted-foreground transition-colors group-hover:inline-flex group-hover:bg-primary group-hover:text-primary-foreground"
								>
									<Icon icon="ph:arrow-right" class="size-6" />
								</div>
							</Card.Action>
						</Card.Header>
					</Card.Root>
				</a>
			{/each}
		{/if}
	</div>
</main>
