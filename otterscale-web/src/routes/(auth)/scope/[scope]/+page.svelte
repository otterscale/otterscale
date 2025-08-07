<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { HomeCard, HomeCell } from '$lib/components/scopes';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths, pathDisabled } from '$lib/path';
	import { currentCeph, currentKubernetes } from '$lib/stores';

	const cards = $derived([
		{
			background: 'bg-[#1c77c3]/30',
			path: dynamicPaths.models(page.params.scope),
			description: m.models_description()
		},
		{
			background: 'bg-[#39a9db]/30',
			path: dynamicPaths.applications(page.params.scope),
			description: m.applications_description()
		},
		{
			background: 'bg-[#f39237]/30',
			path: dynamicPaths.storage(page.params.scope),
			description: m.storage_description()
		},
		{
			background: 'bg-[#d63230]/30',
			path: dynamicPaths.machines(page.params.scope),
			description: m.machines_description()
		}
	]);
</script>

<!-- just-in-time  -->
<dummy class="bg-[#1c77c3]/30"></dummy>
<dummy class="bg-[#39a9db]/30"></dummy>
<dummy class="bg-[#f39237]/30"></dummy>
<dummy class="bg-[#d63230]/30"></dummy>
<dummy class="text-[#f0424d]"></dummy>
<dummy class="text-[#326de6]"></dummy>
<dummy class="text-[#c72c48]"></dummy>
<dummy class="text-[#dd4813]"></dummy>
<dummy class="text-[#ea4b71]"></dummy>

<div class="flex h-full flex-col justify-between">
	<div class="mx-auto flex max-w-5xl px-4 xl:px-0">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
			{#each cards as card}
				<HomeCard
					background={card.background}
					path={card.path}
					description={card.description}
					disabled={pathDisabled(
						$currentCeph?.name,
						$currentKubernetes?.name,
						page.params.scope,
						card.path.url
					)}
				/>
			{/each}
		</div>
	</div>

	<div class="bg-muted -mx-2 -my-4 hidden py-6 md:-mx-4 md:-my-6 md:py-8 lg:block">
		<div class="mx-auto max-w-5xl px-6">
			<div class="grid items-center gap-4 sm:grid-cols-2">
				<div class="dark:bg-muted/50 relative mx-auto w-fit">
					<div class="to-muted absolute inset-0 z-10 bg-radial from-transparent to-75%"></div>
					<div class="mx-auto mb-2 flex w-fit justify-center gap-2">
						<HomeCell icon="logos:postgresql" />
						<HomeCell icon="ph:circle-dashed" color="#f0424d" />
						<HomeCell icon="ph:circle-dashed" color="#326de6" />
						<HomeCell icon="ph:circle-dashed" color="#f0424d" />
					</div>
					<div class="mx-auto my-2 flex w-fit justify-center gap-2">
						<HomeCell icon="simple-icons:n8n" color="#ea4b71" />
						<HomeCell icon="logos:kubernetes" />
						<HomeCell icon="fluent-emoji-flat:otter" />
						<HomeCell icon="simple-icons:ceph" color="#f0424d" />
						<HomeCell icon="simple-icons:minio" color="#c72c48" />
					</div>
					<div class="mx-auto flex w-fit justify-center gap-2">
						<HomeCell icon="ph:circle-dashed" color="#326de6" />
						<HomeCell icon="simple-icons:maas" color="#dd4813" />
						<HomeCell icon="logos:juju" />
					</div>
				</div>
				<div class="mx-auto mt-6 max-w-lg space-y-4 text-center sm:mt-0 sm:text-left">
					<h2 class="text-3xl font-semibold text-balance">
						{m.home_integration()}
					</h2>
					<p class="text-muted-foreground">
						{m.home_integration_description()}
					</p>

					<Button href={dynamicPaths.applications(page.params.scope).url}>
						{m.home_get_started()}
						<Icon icon="ph:cursor-click-fill" />
					</Button>
				</div>
			</div>
		</div>
	</div>
</div>
