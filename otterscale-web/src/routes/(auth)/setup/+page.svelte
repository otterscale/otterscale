<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { homePath, setupPath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import * as Code from '$lib/components/custom/code';
	import { m } from '$lib/paraglide/messages';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [homePath], current: setupPath });

	const code = `sh -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/main/scripts/install.sh" -- url=${PUBLIC_API_URL})`;
</script>

<div class="flex flex-col items-center">
	<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.setup_environment()}</h2>
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.setup_environment_description()}
	</p>
	<div class="min-w-6xl p-6">
		<Code.Root lang="bash" class="w-full" variant="secondary" {code} hideLines>
			<Code.CopyButton />
		</Code.Root>
	</div>
</div>
