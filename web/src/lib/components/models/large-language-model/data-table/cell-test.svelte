<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

	import { onMount } from 'svelte';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import Prompting from '$lib/components/prompting/index.svelte';

	let { model, serviceUri, scope }: { model: Model; serviceUri: string; scope: string } = $props();

	const readyPods = $derived(
		model.pods.filter((pod) => {
			const match = pod.ready.match(/^(\d+)\/(\d+)$/);
			if (!match) return false;
			return Number(match[1]) / Number(match[2]) === 1;
		})
	);
	const isReady = $derived(readyPods.length > 0);

	let tags = $state([] as string[]);
	onMount(async () => {
		const response = await fetch(`https://huggingface.co/api/models/${model.id}`);
		const body = await response.json();
		tags = body.tags;
	});
</script>

{#if tags.includes('text-generation')}
	<Prompting {serviceUri} {model} {scope} disabled={!isReady} />
{/if}
