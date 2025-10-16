<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy } from 'svelte';
	import { fly } from 'svelte/transition';

	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	let { namespace, pod }: { namespace: string; pod: Application_Pod } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	let open = $state(false);
	let messages = $state<{ message: string; phase: string }[]>([]);
	let shouldStop = $state(false);
	let terminal = $state<HTMLDivElement>();

	async function watchLogs() {
		try {
			messages = [];
			shouldStop = false;

			const stream = applicationClient.watchLogs({
				scope: $currentKubernetes?.scope || '',
				facility: $currentKubernetes?.name || '',
				namespace,
				podName: pod.name,
				containerName: '',
			});
			for await (const response of stream) {
				if (shouldStop) break;

				messages = [...messages, { message: response.log, phase: 'LOG' }];
				scrollToBottom();

				// Add delay for better UX
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
		} catch (error) {
			// Ignore abort-related errors
			const isAbortError =
				error instanceof Error &&
				(error.name === 'AbortError' ||
					error.message.includes('aborted') ||
					error.message.includes('canceled'));

			if (!isAbortError) {
				console.error('Error watching logs:', error);
			}
		}
	}

	function scrollToBottom() {
		if (terminal) {
			terminal.scrollTop = terminal.scrollHeight;
		}
	}

	function stopWatching() {
		shouldStop = true;
		messages = [];
	}

	// Stop watching when sheet closes
	$effect(() => {
		if (!open && shouldStop === false) {
			stopWatching();
		}
	});

	onDestroy(() => {
		stopWatching();
	});
</script>

<Sheet.Root bind:open>
	<Sheet.Trigger onclick={watchLogs}>
		<span class="flex items-center gap-1">
			<Icon icon="ph:file-text" />
			{m.log()}
		</span>
	</Sheet.Trigger>

	<Sheet.Content class="rounded-l-lg border-none bg-transparent sm:max-w-9/10">
		{#if open}
			<div
				class="border-border bg-secondary dark text-card-foreground size-full flex-col rounded-l-lg border font-mono text-sm shadow-sm"
			>
				<div bind:this={terminal} class="h-full flex-col overflow-auto p-4">
					{#each messages as msg, i}
						{@const isLastMessage = messages.length === i + 1}
						{@const textClass = isLastMessage ? '' : 'text-green-500'}

						{#if msg.message}
							<span class="block" transition:fly={{ y: -5, duration: 500 }}>
								<span class="flex space-x-1 {textClass}">
									<span>[{msg.phase}] {msg.message}</span>
								</span>
							</span>
						{/if}
					{/each}
				</div>
			</div>
		{/if}
	</Sheet.Content>
</Sheet.Root>
