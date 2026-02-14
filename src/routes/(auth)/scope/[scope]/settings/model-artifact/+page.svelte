<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ModelArtifact } from '$lib/components/settings/model-artifact';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! })
		},
		{
			title: m.model_artifact(),
			url: resolve('/(auth)/scope/[scope]/settings/model-artifact', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<ModelArtifact scope={page.params.scope!} />
{/key}
