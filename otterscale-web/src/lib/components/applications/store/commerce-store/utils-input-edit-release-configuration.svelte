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
	import DynamicInput from './utils/dynamic-input.svelte';

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
		valuesMap = {
			nodePort: '',
			storageClasses: []
		} as any;
	}

	const transport: Transport = getContext('transport');

	const chartMetadata = writable<Application_Chart_Metadata>();
	const plugins: Plugin[] = [
		gfmPlugin(),
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { shell, yaml } }]
		}
	];
	let values = $state('');
	let readme = $state('');
	let open = $state(false);
	let basicRequest = $state({
		nodePort: '',
		storageClasses: []
	});
	let isChartMetadataLoading = $state(true);

	const client = createClient(ApplicationService, transport);

	async function fetchChartMetadata(chartRef: string) {
		try {
			await client
				.getChartMetadata({
					chartRef: chartRef
				})
				.then((response) => {
					chartMetadata.set(response);
				});

			values = $chartMetadata.valuesYaml;
			readme = $chartMetadata.readmeMd;

			isChartMetadataLoading = false;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}
	onMount(async () => {
		try {
			await fetchChartMetadata(chartRef);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		View/Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content class="h-[70vh] max-w-[70vw] min-w-[70vw]">
		<Tabs.Root value="basic" class="w-full">
			<AlertDialog.Header>
				<div class="flex items-start justify-between gap-2">
					<div>
						<AlertDialog.Title>Edit Profile</AlertDialog.Title>
						<AlertDialog.Description>
							<div class="flex items-center gap-1">
								For more information, view the chart documentation at
								<a
									href={chartRef}
									target="_blank"
									rel="noopener noreferrer"
									class="flex items-center gap-1"
								>
									{chartRef}
									<Icon icon="ph:arrow-square-out" />
								</a>
							</div>
						</AlertDialog.Description>
					</div>
					<div class="flex justify-end">
						<Tabs.List>
							<Tabs.Trigger
								value="basic"
								disabled={!valuesMap || Object.keys(valuesMap).length === 0}
							>
								Basic
							</Tabs.Trigger>
							<Tabs.Trigger value="advance">Advance</Tabs.Trigger>
						</Tabs.List>
					</div>
				</div>
			</AlertDialog.Header>
			<Resizable.PaneGroup direction="horizontal" class="flex items-center justify-center">
				<Resizable.Pane defaultSize={50} class="h-[50vh]">
					{#if isChartMetadataLoading}
						<div class="flex-col space-y-4 pr-4">
							{#each Array(3) as _}
								<Skeleton class="h-[40px] w-full" />
								<Skeleton class="h-[20px] w-3/4" />
								<Skeleton class="h-[20px] w-1/2" />
								<Skeleton class="h-[20px] w-1/2" />
							{/each}
						</div>
					{:else}
						<div class="markdown h-full max-w-[35vw] overflow-auto">
							<Markdown {plugins} md={readme.substring(readme.indexOf('#'))} />
						</div>
					{/if}
				</Resizable.Pane>
				<Resizable.Handle withHandle class="h-[50vh]" />
				<Resizable.Pane defaultSize={50} class="h-[50vh]">
					<Tabs.Content value="basic">
						<div class="h-full overflow-auto rounded-lg p-2">
							<DynamicInput bind:data={valuesMap} />
						</div>
					</Tabs.Content>
					<Tabs.Content value="advance">
						<div class="h-[100vh]">
							<Monaco
								options={{
									language: 'yaml',
									automaticLayout: true,
									padding: { top: 32, bottom: 32 },
									overviewRulerBorder: false,
									hideCursorInOverviewRuler: true
								}}
								theme="vs-dark"
								bind:value={values}
							/>
						</div>
					</Tabs.Content>
				</Resizable.Pane>
			</Resizable.PaneGroup>
			<AlertDialog.Footer>
				<AlertDialog.Cancel class="mr-auto">Cancel</AlertDialog.Cancel>
				<AlertDialog.Action
					onclick={() => {
						open = false;
						valuesYaml = values;
					}}
				>
					Confirm
				</AlertDialog.Action>
			</AlertDialog.Footer>
		</Tabs.Root>
	</AlertDialog.Content>
</AlertDialog.Root>

<!-- <style lang="postcss">
	.markdown :global(h1) {
		@apply mt-6 mb-4 border-b border-gray-200 pb-2 text-3xl font-bold;
	}
	.markdown :global(h2) {
		@apply mt-6 mb-4 border-b border-gray-200 pb-2 text-2xl font-bold;
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
</style> -->
