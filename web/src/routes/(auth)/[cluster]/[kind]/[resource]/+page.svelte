<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Component } from 'svelte';

	import { page } from '$app/state';
	import { type GetRequest, ResourceService } from '$lib/api/resource/v1/resource_pb';
	import LogoImage from '$lib/assets/logo.svg';
	import { getResourceInspector } from '$lib/components/dynamical-table/inspectors';
	import ResourceInspector from '$lib/components/dynamical-table/resource-Inspector.svelte';

	const cluster = $derived(page.params.cluster ?? '');
	const resource = $derived(page.params.resource ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const kind = $derived(page.url.searchParams.get('kind') ?? '');
	const namespace = $derived(page.url.searchParams.get('namespace') ?? '');
	const name = $derived(page.url.searchParams.get('name') ?? '');

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	// eslint-disable-next-line
	let object = $state<any>(undefined);

	let getAbortController: AbortController | null = null;

	let isMounted = $state(false);
	let isGetting = $state(false);
	let isDestroyed = false;

	async function GetResource() {
		if (isGetting || isDestroyed) return;

		isGetting = true;
		getAbortController = new AbortController();

		try {
			const response = await resourceClient.get(
				{
					cluster,
					namespace,
					group,
					version,
					resource,
					name
				} as GetRequest,
				{ signal: getAbortController.signal }
			);
			console.log(response);

			object = response.object;
		} catch (error) {
			if (error instanceof Error && error.name === 'ConnectError') {
				if (error.cause === 'Aborted due to component destroyed.') {
					return;
				}
			}

			console.error('Failed to get resource:', error);

			return null;
		} finally {
			isGetting = false;
			getAbortController = null;
		}
	}

	onMount(async () => {
		await GetResource();

		isMounted = true;
	});
	onDestroy(() => {
		isDestroyed = true;
		if (getAbortController) {
			getAbortController.abort('Aborted due to component destroyed.');
			getAbortController = null;
		}
	});

	const Inspector: Component<{ object: any }> = $derived(getResourceInspector(resource));
</script>

{#if isMounted}
	<ResourceInspector {cluster} {namespace} {group} {version} {kind} {resource} {name} {object}>
		<Inspector {object} />
	</ResourceInspector>
{:else}
	<main
		class="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4"
	>
		<!-- Decorative background elements -->
		<div
			class="absolute top-20 left-10 h-40 w-40 rounded-full border border-border/20 opacity-50"
		></div>
		<div
			class="absolute right-10 bottom-20 h-52 w-52 rounded-full border border-border/10 opacity-30"
		></div>
		<div class="relative z-10 w-full max-w-md space-y-8 text-center">
			<!-- Main animated icon -->
			<div class="flex justify-center">
				<div class="relative h-32 w-32">
					<!-- Outer rotating ring -->
					<div
						class="absolute inset-0 animate-spin rounded-full border-2 border-transparent border-t-primary"
					></div>
					<!-- Middle ring - counter spin -->
					<div
						class="absolute inset-3 animate-spin rounded-full border border-primary/40"
						style="animation-direction: reverse; animation-duration: 3s;"
					></div>
					<!-- First orbital layer - 3 dots -->
					<div class="absolute inset-0 flex items-center justify-center">
						<div class="relative h-20 w-20">
							{#each [0, 1, 2] as i (i)}
								<div
									class="absolute h-1.5 w-1.5 rounded-full bg-primary"
									style="animation: orbit 3s linear infinite; transform-origin: 40px 40px; transform: rotate({i *
										120}deg) translateX(40px);"
								></div>
							{/each}
						</div>
					</div>
					<!-- Second orbital layer - 4 dots, smaller, faster, offset -->
					<div class="absolute inset-0 flex items-center justify-center">
						<div class="relative h-14 w-14">
							{#each [0, 1, 2, 3] as i (i)}
								<div
									class="absolute h-1 w-1 rounded-full bg-primary/70"
									style="animation: orbit 2s linear infinite; animation-direction: reverse; transform-origin: 28px 28px; transform: rotate({i *
										90}deg) translateX(28px);"
								></div>
							{/each}
						</div>
					</div>
					<!-- Third orbital layer - 5 dots, tiny, very fast -->
					<div class="absolute inset-0 flex items-center justify-center">
						<div class="relative h-9 w-9">
							{#each [0, 1, 2, 3, 4] as i (i)}
								<div
									class="absolute h-0.5 w-0.5 rounded-full bg-primary/50"
									style="animation: orbit 1.2s linear infinite; transform-origin: 18px 18px; transform: rotate({i *
										72}deg) translateX(18px);"
								></div>
							{/each}
						</div>
					</div>
					<!-- Center icon -->
					<div class="absolute inset-0 flex items-center justify-center">
						<img src={LogoImage} alt="logo" class="size-16" />
					</div>
				</div>
			</div>
			<!-- Main heading - fixed height to prevent floating -->
			<div class="flex h-24 flex-col justify-center space-y-4">
				<h1 class="animate-breathe text-4xl font-bold tracking-tight text-foreground">Loading</h1>
				<div class="mx-auto h-px w-12 bg-primary"></div>
				<p class="animate-breathe text-sm text-muted-foreground" style="animation-delay: 0.2s;">
					Fetching resource
				</p>
			</div>
			<!-- Additional description -->
			<div class="space-y-3 py-4">
				<p class="text-xs leading-relaxed text-muted-foreground">
					We're fetching the latest configuration data from ourservers.
				</p>
				<p class="text-xs text-muted-foreground/70">
					This typically takes a few moments. Thank you for your patience.
				</p>
			</div>
			<!-- Animated horizontal bars -->
			<div class="flex h-30 items-center justify-center gap-1.5 py-6">
				{#each [0, 1, 2, 3] as i (i)}
					<div
						class="w-1 rounded-full bg-primary/60"
						style="height: 24px; animation: pulse-bar 1.6s ease-in-out infinite; animation-delay: {i *
							150}ms;"
					></div>
				{/each}
			</div>
			<!-- Subtle status indicators -->
			<div class="flex items-center justify-center gap-6 text-xs text-muted-foreground/70">
				<div class="flex items-center gap-1.5">
					<div class="h-1.5 w-1.5 animate-pulse rounded-full bg-primary/40"></div>
					<span>Active</span>
				</div>
				<div class="h-4 w-px bg-border/50"></div>
				<div class="text-muted-foreground/50">Please wait</div>
			</div>
		</div>
		<style>
			@keyframes orbit {
				from {
					transform: rotate(0deg) translateX(var(--translate-distance, 0));
				}
				to {
					transform: rotate(360deg) translateX(var(--translate-distance, 0));
				}
			}
			@keyframes pulse-bar {
				0%,
				100% {
					opacity: 0.4;
					height: 8px;
				}
				50% {
					opacity: 1;
					height: 32px;
				}
			}
			@keyframes breathe {
				0%,
				100% {
					opacity: 1;
				}
				50% {
					opacity: 0.6;
				}
			}
			.animate-breathe {
				animation: breathe 3s ease-in-out infinite;
			}
		</style>
	</main>
{/if}
