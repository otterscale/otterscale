<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentCeph, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.setupScope(page.params.scope)
	});
</script>

<!-- just-in-time  -->
<dummy class="bg-[#326de6]"></dummy>
<dummy class="bg-[#f0424d]"></dummy>

<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.setup_scope()}</h2>

{#if $currentKubernetes || $currentCeph}
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.setup_scope_configured_description()}
	</p>
	<div class="mx-auto max-w-5xl px-4 py-10 xl:px-0">
		<div class="bg-card rounded-xl border shadow-sm">
			<div class="rounded-xl p-4 lg:p-8">
				<div class="grid min-w-2xl grid-cols-1 items-center gap-x-12 gap-y-20 sm:grid-cols-2">
					<div
						class="before:bg-border relative text-center before:absolute before:start-1/2 before:-top-full before:mt-3.5 before:h-20 before:w-px before:-translate-x-1/2 before:rotate-[60deg] before:transform before:shadow-sm first:before:hidden sm:before:-start-6 sm:before:top-1/2 sm:before:mt-0 sm:before:-translate-x-0 sm:before:-translate-y-1/2 sm:before:rotate-12"
					>
						<div class="space-y-2">
							<Icon icon="simple-icons:ceph" class="mx-auto size-14 shrink-0 text-[#f0424d]" />
							<h3 class="text-lg font-semibold sm:text-2xl">Ceph</h3>
						</div>
						<Button
							href={dynamicPaths.setupScopeCeph(page.params.scope).url}
							variant="ghost"
							class="text-muted-foreground text-sm sm:text-base"
						>
							<Icon icon="ph:wrench" class="size-5" />
							{$currentCeph ? $currentCeph.name : '-'}
						</Button>
					</div>
					<div
						class="before:bg-border relative text-center before:absolute before:start-1/2 before:-top-full before:mt-3.5 before:h-20 before:w-px before:-translate-x-1/2 before:rotate-[60deg] before:transform before:shadow-sm first:before:hidden sm:before:-start-6 sm:before:top-1/2 sm:before:mt-0 sm:before:-translate-x-0 sm:before:-translate-y-1/2 sm:before:rotate-12"
					>
						<div class="space-y-2">
							<Icon
								icon="simple-icons:kubernetes"
								class="mx-auto size-14 shrink-0 text-[#326de6]"
							/>
							<h3 class="text-lg font-semibold sm:text-2xl">Kubernetes</h3>
						</div>
						<Button
							href={dynamicPaths.setupScopeKubernetes(page.params.scope).url}
							variant="ghost"
							class="text-muted-foreground text-sm sm:text-base"
						>
							<Icon icon="ph:wrench" class="size-5" />
							{$currentKubernetes ? $currentKubernetes.name : '-'}
						</Button>
					</div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<p class="text-muted-foreground mt-4 text-center text-lg">
		{m.error_generic()}
	</p>
{/if}
