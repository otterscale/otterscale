<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { fly } from 'svelte/transition';

	import { version } from '$app/environment';
	import { env } from '$env/dynamic/public';
	import { BootstrapService, type WatchStatusResponse } from '$lib/api/bootstrap/v1/bootstrap_pb';
	import SquareGridImage from '$lib/assets/square-grid.svg';
	import * as Code from '$lib/components/custom/code';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';

	// Constants
	const INSTALL_CODE = `bash -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/refs/tags/${version}/scripts/install.sh)" -- url=${env.PUBLIC_API_URL}`;
	const RETRY_DELAY = 2000;

	const SetupPhase = {
		Loading: 0,
		Waiting: 1,
		Installing: 2,
		Completed: 3
	} as const;

	// State
	const statusStore = writable<{
		phase: number;
		messages: WatchStatusResponse[];
		newURL: string;
	}>({
		phase: SetupPhase.Loading,
		messages: [],
		newURL: ''
	});

	const { phase, messages, newURL } = $derived($statusStore);
	let terminal: HTMLDivElement | undefined = $state();
	let mounted = $state(false);

	// Context
	const transport: Transport = getContext('transport');
	const client = createClient(BootstrapService, transport);

	// Functions
	function scrollToBottom() {
		terminal?.scrollTo({ top: terminal.scrollHeight });
	}

	function determinePhase(status: WatchStatusResponse): number {
		if (status.newUrl !== '') return SetupPhase.Completed;
		if (status.phase !== '') return SetupPhase.Installing;
		return SetupPhase.Waiting;
	}

	async function watchStatus() {
		while (true) {
			try {
				for await (const status of client.watchStatus({})) {
					statusStore.update((state) => ({
						...state,
						phase: determinePhase(status),
						messages: [...state.messages.slice(-199), status],
						newURL: status.newUrl || state.newURL
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
	onMount(() => {
		watchStatus();
	});

	// Effects
	$effect(() => {
		mounted = messages.length > 0;
		if (mounted) scrollToBottom();
	});
</script>

<main class="relative flex max-h-screen min-h-screen flex-col overflow-hidden bg-sidebar">
	<div class="absolute inset-x-0 top-0 flex h-full w-full items-center justify-center opacity-100">
		<img
			src={SquareGridImage}
			alt="square-grid"
			class="mask-[radial-gradient(75%_75%_at_center,white,transparent)] opacity-90"
		/>
	</div>

	<div class="z-10 mt-16 flex flex-col items-center justify-center p-4 md:mt-18">
		<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">
			{phase === SetupPhase.Completed ? m.setup_environment_complete() : m.setup_environment()}
		</h2>

		{#if phase === SetupPhase.Loading}
			{@render loadingView()}
		{:else if phase === SetupPhase.Waiting}
			{@render waitingView()}
		{:else if phase === SetupPhase.Installing}
			{@render installingView()}
		{:else if phase === SetupPhase.Completed}
			{@render completionView()}
		{/if}
	</div>
</main>

{#snippet loadingView()}
	<Button class="mt-4 text-center text-lg text-muted-foreground" variant="ghost" size="lg" disabled>
		<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
		{m.setup_environment_loading()}
	</Button>
{/snippet}

{#snippet waitingView()}
	<Button class="mt-4 text-center text-lg text-muted-foreground" variant="ghost" size="lg" disabled>
		<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
		{m.setup_environment_waiting()}
	</Button>

	<div class="relative mx-auto flex w-full max-w-6xl grow flex-col sm:mt-48">
		<p class="mt-4 text-center text-lg text-muted-foreground">
			{m.setup_environment_curl_description()}
		</p>

		<div class="w-full p-6">
			<Code.Root lang="bash" class="w-full" variant="secondary" code={INSTALL_CODE} hideLines>
				<Code.CopyButton />
			</Code.Root>
		</div>

		{@render decorativeArrows()}
	</div>
{/snippet}

{#snippet installingView()}
	<Button class="mt-4 text-center text-lg text-muted-foreground" variant="ghost" size="lg" disabled>
		<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
		{m.setup_environment_installing()}
	</Button>

	<div
		class="dark mt-6 aspect-video w-full max-w-6xl flex-col rounded-xl border border-border bg-secondary font-mono text-sm text-card-foreground shadow-sm"
	>
		<div class="flex border-b border-inherit p-4">
			<div class="flex items-center gap-2">
				<div class="size-3 rounded-full bg-[#ff605c]"></div>
				<div class="size-3 rounded-full bg-[#ffbd44]"></div>
				<div class="size-3 rounded-full bg-[#00ca4e]"></div>
			</div>
		</div>
		<div bind:this={terminal} class="h-[calc(100%-64px)] flex-col overflow-auto p-4">
			{#each messages as msg, i (i)}
				{@const isLastMessage = i === messages.length - 1}
				{#if msg.message !== ''}
					<span class="block" transition:fly={{ y: -5, duration: 500 }}>
						<span class="flex space-x-1 {isLastMessage ? '' : 'text-green-500'}">
							<Icon
								icon={isLastMessage ? 'ph:spinner-gap' : 'ph:check-bold'}
								class="size-5 {isLastMessage ? 'animate-spin' : ''}"
							/>
							<span>[{msg.phase}] {msg.message}</span>
						</span>
					</span>
				{/if}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet completionView()}
	<p class="mx-auto mt-4 max-w-2xl text-lg text-muted-foreground">
		{m.setup_environment_complete_description()}
	</p>
	<div class="mt-8 text-center">
		<Button href={newURL} class="inline-flex items-center gap-2">
			{m.goto()}
		</Button>
	</div>
{/snippet}

{#snippet decorativeArrows()}
	<div class="absolute end-0 top-0 hidden translate-x-20 -translate-y-12 md:block">
		<svg
			class="h-auto w-16 text-orange-500"
			width={121}
			height={135}
			viewBox="0 0 121 135"
			fill="none"
			aria-hidden="true"
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
			aria-hidden="true"
		>
			<path
				d="M4 82.4591C54.7956 92.8751 30.9771 162.782 68.2065 181.385C112.642 203.59 127.943 78.57 122.161 25.5053C120.504 2.2376 93.4028 -8.11128 89.7468 25.5053C85.8633 61.2125 130.186 199.678 180.982 146.248L214.898 107.02C224.322 95.4118 242.9 79.2851 258.6 107.02C274.299 134.754 299.315 125.589 309.861 117.539L343 93.4426"
				stroke="currentColor"
				stroke-width={7}
				stroke-linecap="round"
			/>
		</svg>
	</div>
{/snippet}
