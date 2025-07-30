<script lang="ts">
	import type { Snippet } from 'svelte';
	import { AppSidebar } from '$lib/components/layout';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { getPath, homePath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import type { LayoutData } from './$types';

	interface Props {
		data: LayoutData;
		children: Snippet;
	}

	let { data, children }: Props = $props();

	// Computed values
	const currentPath = $derived(getPath($breadcrumb.current));
	const filteredParents = $derived(
		$breadcrumb.parents.filter((parent) => getPath(parent).url !== homePath)
	);
</script>

<svelte:head>
	<title>{currentPath.title} | OtterScale 🦦</title>
</svelte:head>

<Sidebar.Provider>
	<AppSidebar user={data.user} />

	<Sidebar.Inset>
		<header class="flex h-16 shrink-0 items-center gap-2">
			<div class="flex items-center gap-2 px-4">
				<Sidebar.Trigger class="-ml-1" />
				<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />

				<nav aria-label="Breadcrumb">
					<Breadcrumb.Root>
						<Breadcrumb.List>
							{#each filteredParents as parent}
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href={getPath(parent).url}>
										{getPath(parent).title}
									</Breadcrumb.Link>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
							{/each}

							<Breadcrumb.Item>
								<Breadcrumb.Page>{currentPath.title}</Breadcrumb.Page>
							</Breadcrumb.Item>
						</Breadcrumb.List>
					</Breadcrumb.Root>
				</nav>
			</div>
		</header>

		<main class="flex flex-1 flex-col gap-4 py-4 md:gap-6 md:py-6">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
