<script lang="ts">
	import { Check, ChevronsUpDown, Globe } from '@lucide/svelte';
	import type { ComponentProps } from '@sjsf/form';
	import { onMount } from 'svelte';

	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils';

	// Get full IANA timezone list from browser - K8s compatible format
	// This includes all ~400+ IANA timezones
	function getAllTimezones(): string[] {
		try {
			// Intl.supportedValuesOf is available in modern browsers
			const timezones = Intl.supportedValuesOf('timeZone');
			// Ensure UTC is at the top
			return ['Etc/UTC', ...timezones.filter((tz) => tz !== 'Etc/UTC')];
		} catch {
			// Fallback for older browsers
			return [
				'Etc/UTC',
				'Africa/Cairo',
				'Africa/Johannesburg',
				'Africa/Lagos',
				'Africa/Nairobi',
				'America/Anchorage',
				'America/Chicago',
				'America/Denver',
				'America/Los_Angeles',
				'America/New_York',
				'America/Sao_Paulo',
				'America/Toronto',
				'Asia/Dubai',
				'Asia/Hong_Kong',
				'Asia/Kolkata',
				'Asia/Seoul',
				'Asia/Shanghai',
				'Asia/Singapore',
				'Asia/Taipei',
				'Asia/Tokyo',
				'Australia/Melbourne',
				'Australia/Sydney',
				'Europe/Amsterdam',
				'Europe/Berlin',
				'Europe/London',
				'Europe/Moscow',
				'Europe/Paris',
				'Pacific/Auckland',
				'Pacific/Honolulu'
			];
		}
	}

	const TIMEZONES = getAllTimezones();

	let { value = $bindable(), config }: ComponentProps['textWidget'] = $props();

	// Popover state
	let open = $state(false);
	let searchQuery = $state('');
	let triggerRef = $state<HTMLButtonElement | null>(null);

	// Get browser's default timezone
	function getBrowserTimezone(): string {
		try {
			return Intl.DateTimeFormat().resolvedOptions().timeZone;
		} catch {
			return 'Etc/UTC'; // Fallback to UTC
		}
	}

	// Set default timezone from browser on mount if value is empty
	onMount(() => {
		if (!value) {
			const browserTz = getBrowserTimezone();
			// Use browser timezone if it's in our list, otherwise default to UTC
			value = TIMEZONES.includes(browserTz) ? browserTz : 'Etc/UTC';
		}
	});

	// Filter timezones based on search query
	const filteredTimezones = $derived(
		searchQuery
			? TIMEZONES.filter((tz) => tz.toLowerCase().includes(searchQuery.toLowerCase()))
			: TIMEZONES
	);

	// Cache UTC offsets to avoid creating Intl.DateTimeFormat on every render
	const utcOffsetCache = new Map<string, string>();

	// Get UTC offset for a timezone (lazy cached)
	function getUtcOffset(timezone: string): string {
		const cached = utcOffsetCache.get(timezone);
		if (cached !== undefined) return cached;

		try {
			const date = new Date();
			const formatter = new Intl.DateTimeFormat('en-US', {
				timeZone: timezone,
				timeZoneName: 'shortOffset'
			});
			const parts = formatter.formatToParts(date);
			const offsetPart = parts.find((p) => p.type === 'timeZoneName');
			const result = offsetPart?.value || '';
			utcOffsetCache.set(timezone, result);
			return result;
		} catch {
			utcOffsetCache.set(timezone, '');
			return '';
		}
	}

	// Display text for the button
	const displayText = $derived(
		value ? `${value} (${getUtcOffset(value as string)})` : 'Select timezone...'
	);

	function handleSelect(timezone: string) {
		value = timezone;
		open = false;
		searchQuery = '';
	}

	// Check if the form is disabled or read-only
	const isDisabled = $derived(config.schema.readOnly === true);
</script>

<div class="timezone-select-widget">
	<Popover.Root bind:open>
		<Popover.Trigger bind:ref={triggerRef} class="w-full">
			{#snippet child({ props })}
				<Button
					variant="outline"
					class={cn('w-full justify-between', !value && 'text-muted-foreground')}
					disabled={isDisabled}
					{...props}
				>
					<span class="flex items-center gap-2 truncate">
						<Globe class="h-4 w-4 shrink-0" />
						<span class="truncate">{displayText}</span>
					</span>
					<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-75 p-0" align="start">
			<Command.Root shouldFilter={false}>
				<Command.Input
					placeholder="Search timezones..."
					value={searchQuery}
					oninput={(e) => (searchQuery = e.currentTarget.value)}
					onkeydown={(e) => e.stopPropagation()}
				/>
				<Command.List>
					<Command.Empty>No timezones found.</Command.Empty>
					<Command.Group>
						{#each filteredTimezones as timezone (timezone)}
							<Command.Item value={timezone} onSelect={() => handleSelect(timezone)}>
								<Check class={cn('mr-2 h-4 w-4', value !== timezone && 'text-transparent')} />
								<div class="flex flex-col">
									<span class="font-medium">{timezone}</span>
									<span class="text-xs text-muted-foreground">{getUtcOffset(timezone)}</span>
								</div>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
</div>
