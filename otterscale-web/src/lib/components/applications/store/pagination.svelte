<script lang="ts" module>
	import { type Application_Chart } from '$lib/api/application/v1/application_pb';
	import * as Pagination from '$lib/components/ui/pagination';
	import Icon from '@iconify/svelte';
	import type { PaginationManager } from './utils.svelte';
</script>

<script lang="ts">
	let { paginationManager }: { paginationManager: PaginationManager } = $props();
</script>

<Pagination.Root
	class={paginationManager.count === 0 ? 'hidden' : 'visible'}
	count={paginationManager.count}
	perPage={paginationManager.perPage}
	siblingCount={paginationManager.siblingCount}
>
	{#snippet children({ pages, currentPage })}
		<Pagination.Content>
			<Pagination.Item>
				<Pagination.PrevButton
					onclick={() => {
						paginationManager.activePage = Math.max(currentPage - 1, 1) - 1;
					}}
				>
					<Icon icon="ph:caret-left" class="hidden sm:block" />
				</Pagination.PrevButton>
			</Pagination.Item>
			{#each pages as page (page.key)}
				{#if page.type === 'ellipsis'}
					<Pagination.Item>
						<Pagination.Ellipsis />
					</Pagination.Item>
				{:else}
					<Pagination.Item
						onclick={() => {
							paginationManager.activePage = currentPage - 1;
						}}
					>
						<Pagination.Link
							size="icon"
							class="h-7 w-7"
							{page}
							isActive={currentPage === page.value}
						>
							{page.value}
						</Pagination.Link>
					</Pagination.Item>
				{/if}
			{/each}
			<Pagination.Item>
				<Pagination.NextButton
					onclick={() => {
						paginationManager.activePage =
							Math.min(
								currentPage + 1,
								Math.ceil(paginationManager.count / paginationManager.perPage)
							) - 1;
					}}
				>
					<Icon icon="ph:caret-right" class="hidden sm:block" />
				</Pagination.NextButton>
			</Pagination.Item>
		</Pagination.Content>
	{/snippet}
</Pagination.Root>
