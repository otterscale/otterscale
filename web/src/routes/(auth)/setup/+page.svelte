<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { fly } from 'svelte/transition';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { env } from '$env/dynamic/public';
	import {
		EnvironmentService,
		type WatchStatusResponse
	} from '$lib/api/environment/v1/environment_pb';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { Button } from '$lib/components/ui/button';
	import * as Code from '$lib/components/custom/code';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths, staticPaths, urlIcon } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Constants
	const INSTALL_CODE = `sh -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/main/scripts/install.sh" -- url=${env.PUBLIC_API_URL})`;
	const RETRY_DELAY = 2000;
	const services = (scope: Scope) => [
		{ path: dynamicPaths.models(scope.name), description: m.models_description() },
		{ path: dynamicPaths.databases(scope.name), description: m.databases_description() },
		{
			path: dynamicPaths.applications(scope.name),
			description: m.applications_description()
		},
		{ path: dynamicPaths.storage(scope.name), description: m.storage_description() },
		{ path: dynamicPaths.machines(scope.name), description: m.machines_description() },
		{ path: dynamicPaths.settings(scope.name), description: m.settings_description() }
	];

	// Types
	interface SetupState {
		started: boolean;
		finished: boolean;
		messages: WatchStatusResponse[];
	}

	// State
	const statusStore = writable<SetupState>({ started: false, finished: false, messages: [] });
	const { started, finished, messages } = $derived($statusStore);

	let terminal: HTMLDivElement | undefined = $state();
	let mounted = $state(false);

	// Context
	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const environmentClient = createClient(EnvironmentService, transport);
	const scopes = writable<Scope[]>([]);

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [staticPaths.home], current: staticPaths.setup });

	// Functions
	function scrollToBottom() {
		if (terminal) {
			terminal.scrollTop = terminal.scrollHeight;
		}
	}

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function watchStatus() {
		while (true) {
			try {
				for await (const status of environmentClient.watchStatus({})) {
					statusStore.update((state) => ({
						...state,
						started: status.started,
						finished: status.finished,
						messages: [...state.messages, status]
					}));
				}
				break;
			} catch (error) {
				console.error('Error watching statuses:', error);
				await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY));
			}
		}
	}

	// Lifecycle
	onMount(async () => {
		await watchStatus();
	});

	// Effects
	$effect(() => {
		mounted = messages.length > 0;
		if (mounted) {
			scrollToBottom();
		}

		if (finished) {
			fetchScopes();
		}
	});
</script>

<main class="flex flex-1 flex-col px-2 py-20 md:px-4 md:py-24">
	<div class="flex flex-col items-center justify-center">
		<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">
			{m.setup_environment()}
		</h2>

		{#if mounted}
			{#if finished}
				<!-- Completion View -->
				<p class="text-muted-foreground mx-auto mt-4 max-w-2xl text-lg">
					{m.setup_environment_complete_description()}
				</p>

				<div class="mx-auto max-w-4xl sm:mt-12">
					<div class="grid grid-cols-1 gap-x-12 gap-y-6 md:grid-cols-2">
						{#if $scopes.length > 0}
							{#each services($scopes[0]) as service}
								<div
									class="group hover:border-primary/20 hover:bg-muted/50 flex items-center gap-4 rounded-lg border border-transparent p-3 transition-all duration-300"
								>
									<div class="bg-primary/10 rounded-full p-3">
										<Icon icon={urlIcon(service.path.url)} class="size-6" />
									</div>
									<div>
										<h3 class="font-medium">{service.path.title}</h3>
										<p class="text-muted-foreground text-sm">
											{service.description}
										</p>
									</div>
									<div class="ml-auto">
										<a
											href={service.path.url}
											class="text-muted-foreground group-hover:bg-primary group-hover:text-primary-foreground inline-flex h-7 w-7 items-center justify-center rounded-full transition-colors"
										>
											<Icon icon="ph:arrow-right" class="size-4" />
										</a>
									</div>
								</div>
							{/each}
						{/if}
					</div>

					<div class="mt-12 text-center">
						<Button href={staticPaths.home.url} class="inline-flex items-center gap-2">
							{m.back_to_home()}
						</Button>
					</div>
				</div>
			{:else if started}
				<!-- Installation Progress View -->
				<Button
					class="text-muted-foreground mt-4 text-center text-lg"
					variant="ghost"
					size="lg"
					disabled
				>
					<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
					{m.setup_environment_installing()}
				</Button>

				<div
					class="border-border dark bg-card text-card-foreground m-6 aspect-video w-full max-w-6xl flex-col rounded-xl border font-mono text-sm shadow-sm"
				>
					<div class="flex border-b border-inherit p-4">
						<div class="flex items-center gap-2">
							<div class="size-3 rounded-full bg-[#ff605c]"></div>
							<div class="size-3 rounded-full bg-[#ffbd44]"></div>
							<div class="size-3 rounded-full bg-[#00ca4e]"></div>
						</div>
					</div>
					<div bind:this={terminal} class="h-[calc(100%-64px)] flex-col overflow-auto p-4">
						{#each messages as msg, i}
							{@const isLastMessage = messages.length === i + 1}
							{@const iconName = isLastMessage ? 'ph:spinner-gap' : 'ph:check-bold'}
							{@const iconClass = isLastMessage ? 'animate-spin' : ''}
							{@const textClass = isLastMessage ? '' : 'text-green-500'}

							{#if msg.message !== ''}
								<span class="block" transition:fly={{ y: -5, duration: 500 }}>
									<span class="flex space-x-1 {textClass}">
										<Icon icon={iconName} class="size-5 {iconClass}" />
										<span>[{msg.phase}] {msg.message}</span>
									</span>
								</span>
							{/if}
						{/each}
					</div>
				</div>
			{:else}
				<!-- Installation View -->
				<Button
					class="text-muted-foreground mt-4 text-center text-lg"
					variant="ghost"
					size="lg"
					disabled
				>
					<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
					{m.setup_environment_waiting()}
				</Button>

				<div class="relative mx-auto flex w-full max-w-6xl flex-grow flex-col sm:mt-48">
					<p class="text-muted-foreground mt-4 text-center text-lg">
						{m.setup_environment_curl_description()}
					</p>

					<div class="w-full p-6">
						<Code.Root lang="bash" class="w-full" variant="secondary" code={INSTALL_CODE} hideLines>
							<Code.CopyButton />
						</Code.Root>
					</div>

					<!-- Decorative elements -->
					<div class="absolute end-0 top-0 hidden translate-x-20 -translate-y-12 md:block">
						<svg
							class="h-auto w-16 text-orange-500"
							width={121}
							height={135}
							viewBox="0 0 121 135"
							fill="none"
						>
							<path
								d="M5 16.4754C11.7688 27.4499 21.2452 57.3224 5 89.0164"
								stroke="currentColor"
								stroke-width={10}
								stroke-linecap="round"
							/>
							<path
								d="M33.6761 112.104C44.6984 98.1239 74.2618 57.6776 83.4821 5"
								stroke="currentColor"
								stroke-width={10}
								stroke-linecap="round"
							/>
							<path
								d="M50.5525 130C68.2064 127.495 110.731 117.541 116 78.0874"
								stroke="currentColor"
								stroke-width={10}
								stroke-linecap="round"
							/>
						</svg>
					</div>

					<div class="absolute start-0 bottom-0 hidden -translate-x-32 translate-y-10 md:block">
						<svg
							class="h-auto w-40 text-cyan-500"
							width={347}
							height={188}
							viewBox="0 0 347 188"
							fill="none"
						>
							<path
								d="M4 82.4591C54.7956 92.8751 30.9771 162.782 68.2065 181.385C112.642 203.59 127.943 78.57 122.161 25.5053C120.504 2.2376 93.4028 -8.11128 89.7468 25.5053C85.8633 61.2125 130.186 199.678 180.982 146.248L214.898 107.02C224.322 95.4118 242.9 79.2851 258.6 107.02C274.299 134.754 299.315 125.589 309.861 117.539L343 93.4426"
								stroke="currentColor"
								stroke-width={7}
								stroke-linecap="round"
							/>
						</svg>
					</div>
				</div>
			{/if}
		{:else}
			<!-- Loading View -->
			<Button
				class="text-muted-foreground mt-4 text-center text-lg"
				variant="ghost"
				size="lg"
				disabled
			>
				<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
				{m.setup_environment_loading()}
			</Button>
		{/if}
	</div>
</main>
