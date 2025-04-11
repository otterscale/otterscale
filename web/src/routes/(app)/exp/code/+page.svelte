<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

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
</script>

<div class="markdown-body">
	<Markdown {md} {plugins} />
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
		<div class="grid h-[60vh] grid-cols-2 items-center gap-4">
			{#if isLoading}
				<Skeleton class="h-[20px] w-[100px] rounded-full" />
			{:else}
				<div class="markdown-body h-full overflow-auto border">
					<Markdown {plugins} md={readme} />
				</div>
			{/if}
			<div class="h-full">
				<Monaco
					options={{ language: 'yaml', automaticLayout: true }}
					theme="vs-dark"
					on:ready={(event) => console.log(event.detail)}
					bind:value={values}
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button type="submit">Save changes</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
