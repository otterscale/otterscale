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

	// Derived breadcrumb values
	const breadcrumbData = $derived.by(() => {
		const parent = getPath($breadcrumb.parent);
		const current = getPath($breadcrumb.current);

		return {
			parentURL: parent?.url,
			parentTitle: parent?.title,
			currentTitle: current?.title,
			showParentBreadcrumb: parent?.url !== homePath
		};
	});
</script>

<svelte:head>
	<title>{breadcrumbData.currentTitle} | OtterScale 🦦</title>
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
							{#if breadcrumbData.showParentBreadcrumb}
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href={breadcrumbData.parentURL}>
										{breadcrumbData.parentTitle}
									</Breadcrumb.Link>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
							{/if}
							<Breadcrumb.Item>
								<Breadcrumb.Page>{breadcrumbData.currentTitle}</Breadcrumb.Page>
							</Breadcrumb.Item>
						</Breadcrumb.List>
					</Breadcrumb.Root>
				</nav>
			</div>
		</header>

		<main class="flex flex-1 flex-col gap-4 p-4 pt-0">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
