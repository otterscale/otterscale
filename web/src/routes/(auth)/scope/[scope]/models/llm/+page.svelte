<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import LLM from '$lib/components/models/large-language-model/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.models(),
			url: resolve('/(auth)/scope/[scope]/models', { scope: page.params.scope! })
		},
		{
			title: m.llm(),
			url: resolve('/(auth)/scope/[scope]/models/llm', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<LLM scope={page.params.scope!} />
{/key}
