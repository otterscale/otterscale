<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Resizable from '$lib/components/ui/resizable';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Nexus, type Application_Chart_Metadata } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { get, writable } from 'svelte/store';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Monaco from 'svelte-monaco';

	import type { Plugin } from 'svelte-exmarkdown';
	import Markdown from 'svelte-exmarkdown';

	import rehypeHighlight from 'rehype-highlight';

	import { gfmPlugin } from 'svelte-exmarkdown/gfm';

	import shell from 'highlight.js/lib/languages/shell';
	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';
	import Unstruct from './unstruct.svelte';

	let {
		chartRef,
		valuesYaml = $bindable(),
		valuesMap = $bindable()
	}: {
		chartRef: string;
		valuesYaml: string;
		valuesMap: { [key: string]: string };
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

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
	let md = $state('');

	let values = $state('');
	let readme = $state('');

	let open = $state(false);

	let basicRequest = $state({
		nodePort: '',
		storageClasses: []
	});

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
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>View/Edit</AlertDialog.Trigger
	>
	<AlertDialog.Content class="min-w-[62vw] max-w-[78vw]">
		<AlertDialog.Header>
			<AlertDialog.Title>Edit Profile</AlertDialog.Title>
			<AlertDialog.Description class="flex items-center gap-1">
				For more information, view the chart documentation at
				<a href={chartRef} target="_blank" rel="noopener noreferrer" class="flex items-center gap-1"
					>{chartRef} <Icon icon="ph:arrow-square-out" /></a
				>
			</AlertDialog.Description>
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
				<Tabs.Root
					value={!valuesMap || Object.keys(valuesMap).length === 0 ? 'advance' : 'basic'}
					class="p-2"
				>
					<Tabs.List class="h-[40px] w-fit">
						<Tabs.Trigger value="basic" disabled={!valuesMap || Object.keys(valuesMap).length === 0}
							>Basic</Tabs.Trigger
						>
						<Tabs.Trigger value="advance">Advance</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="basic" class="">
						<div class="grid h-full max-h-[calc(70vh_-_40px)] gap-4 overflow-auto p-4">
							<Unstruct bind:data={valuesMap} />
						</div>
					</Tabs.Content>
					<Tabs.Content value="advance" class="h-[70vh]">
						<Monaco
							options={{ language: 'yaml', automaticLayout: true }}
							theme="vs-dark"
							bind:value={values}
						/>
					</Tabs.Content>
				</Tabs.Root>
			</Resizable.Pane>
		</Resizable.PaneGroup>
		<AlertDialog.Footer>
			<AlertDialog.Cancel class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					open = false;
					valuesYaml = values;
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<style lang="postcss">
	.markdown :global(h1) {
		@apply mb-4 mt-6 border-b border-gray-200 pb-2 text-3xl font-bold;
	}
	.markdown :global(h2) {
		@apply mb-4 mt-6 border-b border-gray-200 pb-2 text-2xl font-bold;
	}
	.markdown :global(h3) {
		@apply mb-2 mt-5 text-xl font-semibold;
	}
	.markdown :global(h4) {
		@apply mb-2 mt-4 text-lg font-semibold;
	}
	.markdown :global(h5) {
		@apply mb-2 mt-4 text-base font-semibold;
	}
	.markdown :global(h6) {
		@apply mb-2 mt-4 text-base font-semibold text-gray-600;
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
