<script lang="ts">
	import {
		ApplicationService,
		type Application_Chart_Metadata
	} from '$lib/api/application/v1/application_pb';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Resizable from '$lib/components/ui/resizable';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import shell from 'highlight.js/lib/languages/shell';
	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';
	import rehypeHighlight from 'rehype-highlight';
	import { getContext, onMount } from 'svelte';
	import type { Plugin } from 'svelte-exmarkdown';
	import Markdown from 'svelte-exmarkdown';
	import { gfmPlugin } from 'svelte-exmarkdown/gfm';
	import Monaco from 'svelte-monaco';
	import { writable } from 'svelte/store';
	import Unstruct from './utils/dynamic-input.svelte';

	let {
		chartRef,
		valuesYaml = $bindable(),
		valuesMap = $bindable()
	}: {
		chartRef: string;
		valuesYaml: string;
		valuesMap?: { [key: string]: any };
	} = $props();

	if (!valuesMap) {
		valuesMap = { nodePort: '', storageClasses: [] };
	}

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const chartMetadataStore = writable<Application_Chart_Metadata>();
	const chartMetadataLoading = writable(true);
	async function fetchChartMetadata(chartRef: string) {
		try {
			const response = await client.getChartMetadata({
				chartRef: chartRef
			});
			chartMetadataStore.set(response);
			values = $chartMetadataStore.valuesYaml;
			readme = $chartMetadataStore.readmeMd;
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			chartMetadataLoading.set(false);
		}
	}

	const plugins: Plugin[] = [
		gfmPlugin(),
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { shell, yaml } }]
		}
	];
	let tab = $state(!valuesMap || Object.keys(valuesMap).length === 0 ? 'advance' : 'basic');

	let values = $state('');
	let readme = $state('');

	let open = $state(false);

	onMount(async () => {
		try {
			await fetchChartMetadata(chartRef);
			values = $chartMetadataStore.valuesYaml;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		{m.view_edit()}
	</AlertDialog.Trigger>
	<AlertDialog.Content class="max-w-[78vw] min-w-[62vw]">
		<AlertDialog.Header>
			<div class="flex items-end justify-between gap-2">
				<div>
					<AlertDialog.Title>
						{m.edit_release_configuration()}
					</AlertDialog.Title>
					<AlertDialog.Description class="flex items-center gap-1">
						{m.edit_release_configuration_documentation()}
						<a
							href={chartRef}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-center gap-1 underline hover:no-underline"
						>
							{chartRef}
						</a>
					</AlertDialog.Description>
				</div>
				<div class="flex items-center gap-2">
					<Tabs.Root bind:value={tab}>
						<Tabs.List>
							<Tabs.Trigger value="basic">
								{m.basic()}
							</Tabs.Trigger>
							<Tabs.Trigger value="advance">
								{m.advance()}
							</Tabs.Trigger>
						</Tabs.List></Tabs.Root
					>
				</div>
			</div>
		</AlertDialog.Header>
		<Resizable.PaneGroup direction="horizontal">
			<Resizable.Pane defaultSize={50} class="h-[70vh]">
				{#if $chartMetadataLoading}
					<div class="flex-col space-y-4 pr-4">
						{#each Array(3) as _}
							<Skeleton class="h-[40px] w-full" />
							<Skeleton class="h-[20px] w-3/4" />
							<Skeleton class="h-[20px] w-1/2" />
							<Skeleton class="h-[20px] w-1/2" />
						{/each}
					</div>
				{:else}
					<div class="markdown h-full overflow-auto">
						<Markdown {plugins} md={readme.substring(readme.indexOf('#'))} />
					</div>
				{/if}
			</Resizable.Pane>
			<Resizable.Handle withHandle />
			<Resizable.Pane defaultSize={50} class="h-[70vh]">
				<Tabs.Root value={tab}>
					<Tabs.Content value="basic">
						<div class="grid h-full max-h-[calc(70vh_-_40px)] gap-4 overflow-auto px-2">
							<Unstruct bind:data={valuesMap} />
						</div>
					</Tabs.Content>
					<Tabs.Content value="advance" class="h-[70vh]">
						<div class="h-[70vh]">
							<Monaco
								options={{
									language: 'yaml',
									padding: { top: 32, bottom: 8 },
									automaticLayout: true
								}}
								theme="vs-dark"
								bind:value={values}
							/>
						</div>
					</Tabs.Content>
				</Tabs.Root>
			</Resizable.Pane>
		</Resizable.PaneGroup>
		<AlertDialog.Footer>
			<AlertDialog.Cancel class="mr-auto">
				{m.cancel()}
			</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					open = false;
					valuesYaml = values;
				}}
			>
				{m.confirm()}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<style lang="postcss">
	@reference "tailwindcss";

	.markdown :global(h1) {
		@apply mt-8 mb-4 border-b border-gray-200 pb-2 text-3xl font-bold;
	}
	.markdown :global(h2) {
		@apply mt-8 mb-4 border-b border-gray-200 pb-2 text-2xl font-bold;
	}
	.markdown :global(h3) {
		@apply mt-5 mb-2 text-xl font-semibold;
	}
	.markdown :global(h4) {
		@apply mt-4 mb-2 text-lg font-semibold;
	}
	.markdown :global(h5) {
		@apply mt-4 mb-2 text-base font-semibold;
	}
	.markdown :global(h6) {
		@apply mt-4 mb-2 text-base font-semibold text-gray-600;
	}
	.markdown :global(p) {
		@apply mb-4;
	}
	.markdown :global(a) {
		@apply break-words text-blue-600 underline transition-colors hover:text-blue-800;
	}
	.markdown :global(ul),
	.markdown :global(ol) {
		@apply mb-4 pl-8;
	}
	.markdown :global(ul) > :global(li) {
		@apply list-disc;
	}
	.markdown :global(ol) > :global(li) {
		@apply list-decimal;
	}
	.markdown :global(li) {
		@apply mb-1;
	}
	.markdown :global(blockquote) {
		@apply mb-4 rounded border-l-4 border-gray-300 bg-gray-50 pl-4 text-gray-700;
	}
	.markdown :global(pre) {
		@apply mb-4 overflow-x-auto rounded bg-gray-100 p-4 text-sm leading-relaxed;
	}
	.markdown :global(code) {
		@apply rounded bg-gray-200 px-1 py-0.5 font-mono text-sm;
	}
	.markdown :global(pre) :global(code) {
		@apply m-0 rounded-none bg-transparent p-0 text-inherit;
	}
	.markdown :global(table) {
		@apply my-6 w-full border-collapse;
	}
	.markdown :global(th),
	.markdown :global(td) {
		@apply border border-gray-200 px-4 py-2;
	}
	.markdown :global(th) {
		@apply bg-gray-50 font-semibold;
	}
	.markdown :global(img) {
		@apply h-auto max-w-full rounded;
	}
	.markdown :global(hr) {
		@apply my-8 border-t border-gray-200;
	}
	.markdown :global(strong) {
		@apply font-bold;
	}
	.markdown :global(em) {
		@apply italic;
	}
	.markdown :global(del) {
		@apply line-through;
	}
</style>
