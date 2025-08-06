<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import Button from '$lib/components/ui/button/button.svelte';
	import changelogRead from '$lib/stores/changelog';
	import Icon from '@iconify/svelte';
	import type { PageData } from './$types';
	import { m } from '$lib/paraglide/messages';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	changelogRead.set(true);

	const CHANGE_TYPES = {
		feat: {
			title: 'New Features',
			icon: 'ph:check-circle',
			colors: {
				bg: 'bg-green-50 dark:bg-green-950/30',
				text: 'text-green-700 dark:text-green-400',
				border: 'border-green-200 dark:border-green-900'
			}
		},
		fix: {
			title: 'Bug Fixes',
			icon: 'ph:warning',
			colors: {
				bg: 'bg-red-50 dark:bg-red-950/30',
				text: 'text-red-700 dark:text-red-400',
				border: 'border-red-200 dark:border-red-900'
			}
		},
		perf: {
			title: 'Improvements',
			icon: 'ph:arrow-circle-up-right',
			colors: {
				bg: 'bg-blue-50 dark:bg-blue-950/30',
				text: 'text-blue-700 dark:text-blue-400',
				border: 'border-blue-200 dark:border-blue-900'
			}
		},
		refactor: {
			title: 'Refactoring',
			icon: 'ph:recycle',
			colors: {
				bg: 'bg-purple-50 dark:bg-purple-950/30',
				text: 'text-purple-700 dark:text-purple-400',
				border: 'border-purple-200 dark:border-purple-900'
			}
		},
		test: {
			title: 'Tests',
			icon: 'ph:test-tube',
			colors: {
				bg: 'bg-yellow-50 dark:bg-yellow-950/30',
				text: 'text-yellow-700 dark:text-yellow-400',
				border: 'border-yellow-200 dark:border-yellow-900'
			}
		},
		style: {
			title: 'Styling',
			icon: 'ph:palette',
			colors: {
				bg: 'bg-pink-50 dark:bg-pink-950/30',
				text: 'text-pink-700 dark:text-pink-400',
				border: 'border-pink-200 dark:border-pink-900'
			}
		},
		docs: {
			title: 'Documentation',
			icon: 'ph:book',
			colors: {
				bg: 'bg-indigo-50 dark:bg-indigo-950/30',
				text: 'text-indigo-700 dark:text-indigo-400',
				border: 'border-indigo-200 dark:border-indigo-900'
			}
		},
		chore: {
			title: 'Miscellaneous Task',
			icon: 'ph:gear',
			colors: {
				bg: 'bg-gray-50 dark:bg-gray-950/30',
				text: 'text-gray-700 dark:text-gray-400',
				border: 'border-gray-200 dark:border-gray-900'
			}
		}
	} as const;

	function extractPRNumber(url: string): string {
		return `#${url.split('/').pop()}`;
	}

	function hasChanges(release: any, changeKey: keyof typeof CHANGE_TYPES): boolean {
		return release.changes[changeKey]?.length > 0;
	}
</script>

<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.changelog()}</h2>

<p class="text-muted-foreground mt-4 text-center text-lg">
	{m.changelog_description()}
</p>

<Accordion.Root type="single" class="mx-auto w-full py-10 md:max-w-[60%]" value="latest">
	{#each data.releases as release}
		<Accordion.Item value={release.latest ? 'latest' : release.tag_name}>
			<Accordion.Trigger
				class="hover:bg-accent items-center hover:no-underline [[data-state=open]]:rounded-b-none [[data-state=open]]:border-b [&>svg:last-child]:hidden"
			>
				<div class="flex w-full flex-col flex-wrap space-y-1 px-6">
					<span class="text-foreground flex items-center space-x-1 text-lg font-medium">
						{release.name}
					</span>
					<div class="flex space-x-2">
						<div class="flex items-center space-x-1">
							<Icon icon="ph:tag" class="text-muted-foreground size-3.5" />
							<span class="text-muted-foreground text-xs tracking-tight">{release.tag_name}</span>
						</div>
						<div class="flex items-center space-x-1">
							<Icon icon="ph:clock" class="text-muted-foreground size-3.5" />
							<span class="text-muted-foreground text-xs tracking-tight">{release.created_at}</span>
						</div>
					</div>
				</div>

				{#if release.prerelease}
					<Badge variant="outline">Pre-release</Badge>
				{/if}
				{#if release.latest}
					<Badge>Latest</Badge>
				{/if}

				<Button href={release.html_url} variant="ghost" size="icon">
					<Icon icon="ph:arrow-square-out" class="size-4" />
				</Button>

				<Icon
					icon="ph:caret-down"
					class="text-muted-foreground pointer-events-none mr-4 size-5 shrink-0 translate-y-0.5 transition-transform duration-200"
				/>
			</Accordion.Trigger>

			<Accordion.Content class="flex flex-col space-y-4 p-6 text-balance">
				{#each Object.entries(CHANGE_TYPES) as [key, config]}
					{#if hasChanges(release, key as keyof typeof CHANGE_TYPES)}
						<Card.Root class="gap-0 p-0 {config.colors.border}">
							<Card.Header class="rounded-t-xl py-3 {config.colors.bg}">
								<Card.Title
									class="flex items-center space-x-1 text-sm font-semibold {config.colors.text}"
								>
									<Icon icon={config.icon} class="size-5" />
									<span>{config.title}</span>
								</Card.Title>
							</Card.Header>

							<Card.Content class="space-y-2 rounded-b-xl py-6">
								{#each release.changes[key as keyof typeof CHANGE_TYPES] as item}
									<p class="flex items-center space-x-2 text-sm">
										<Icon icon="ph:circle-fill" class="size-1.5" />
										<span>
											{item.description} by
											<a
												target="_blank"
												href="https://github.com/{item.author}"
												class="font-semibold underline underline-offset-4"
											>
												{item.author}
											</a>
											in
											<a
												target="_blank"
												href={item.pull_request}
												class="text-blue-600 underline underline-offset-4 hover:text-blue-800"
											>
												{extractPRNumber(item.pull_request)}
											</a>
										</span>
									</p>
								{/each}
							</Card.Content>
						</Card.Root>
					{/if}
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/each}
</Accordion.Root>
