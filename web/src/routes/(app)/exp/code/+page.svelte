<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Resizable from '$lib/components/ui/resizable';

	import { KubeService } from '$gen/api/kube/v1/kube_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { get, writable } from 'svelte/store';

	import Monaco from 'svelte-monaco';

	import type { Plugin } from 'svelte-exmarkdown';
	import Markdown from 'svelte-exmarkdown';

	import rehypeHighlight from 'rehype-highlight';

	import { gfmPlugin } from 'svelte-exmarkdown/gfm';

	import shell from 'highlight.js/lib/languages/shell';
	import yaml from 'highlight.js/lib/languages/yaml';
	import 'highlight.js/styles/github.css';

	let md = $state(`# Hello world!
 - 1123
 - 1123
 - 1123
 - 1123
 `);

	const plugins: Plugin[] = [
		gfmPlugin(),
		{
			rehypePlugin: [rehypeHighlight, { ignoreMissing: true, languages: { shell, yaml } }]
		}
	];

	// Get the transport out of context
	const transport: Transport = getContext('transport');
	const client = createClient(KubeService, transport);

	let isLoading = $state(true);

	let values = $state('');
	let readme = $state('');

	onMount(async () => {
		await client
			.getChartInfo({
				chartRef: 'oci://registry-1.docker.io/bitnamicharts/nginx:19.0.1'
			})
			.then((res) => {
				values = res.values;
				readme = res.readme;
			})
			.finally(() => {
				isLoading = false;
			});
	});

	let mdx = `# H1
## H2
### H3`;
</script>

<div class="markdown">
	<Markdown md={mdx} {plugins} />
</div>

<Dialog.Root>
	<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>Edit Profile</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[1280px]">
		<Dialog.Header>
			<Dialog.Title>Edit profile</Dialog.Title>
			<Dialog.Description>
				Make changes to your profile here. Click save when you're done.
			</Dialog.Description>
		</Dialog.Header>
		<Resizable.PaneGroup direction="horizontal">
			<Resizable.Pane defaultSize={50} class="h-[70vh]">
				{#if isLoading}
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
				<Monaco
					options={{ language: 'yaml', automaticLayout: true }}
					theme="vs-dark"
					on:ready={(event) => console.log(event.detail)}
					bind:value={values}
				/>
			</Resizable.Pane>
		</Resizable.PaneGroup>
		<Dialog.Footer>
			<Button type="submit">Save changes</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

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
