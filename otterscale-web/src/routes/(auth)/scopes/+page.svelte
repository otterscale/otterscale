<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import SquareGridImage from '$lib/assets/square-grid.svg';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { scopeIcon } from '$lib/components/scopes/icon';
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

	<div
		class="mx-auto grid gap-4 px-4 py-10 sm:gap-6 xl:px-0
	{$scopes.length > 4 ? 'sm:grid-cols-2 lg:grid-cols-4' : 'grid-cols-3'}"
	>
		{#each $scopes as scope}
			<a
				class="bg-card group text-card-foreground relative flex cursor-pointer flex-col gap-6 overflow-hidden rounded-xl shadow-sm hover:shadow-md"
				href={dynamicPaths.scope(scope.name).url}
			>
				<div
					class="text-primary/5 absolute -top-4 -right-4 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				>
					{scope.name}
				</div>
				<div class="relative flex space-x-4 p-6">
					<div class="bg-primary flex size-10 items-center justify-center rounded-lg">
						<Icon
							icon="{scopeIcon($scopes.findIndex((s) => s.name === scope.name))}-fill"
							class="text-primary-foreground size-6"
						/>
					</div>
					<div class="flex flex-col">
						<p class="text-muted-foreground text-xs tracking-wide uppercase">Scope</p>
						<h3 class="text-xl font-medium tracking-wide sm:text-2xl">{scope.name}</h3>
					</div>
					<div class="ml-auto">
						<div
							class="text-muted-foreground group-hover:bg-primary group-hover:text-primary-foreground inline-flex size-10 items-center justify-center rounded-full transition-colors"
						>
							<Icon icon="ph:arrow-right" class="size-6" />
						</div>
					</div>
				</div>
			</a>
		{/each}
	</div>
</main>
