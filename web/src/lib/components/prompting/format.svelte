<div
	role="group"
	data-slot="field"
	data-orientation="vertical"
	class="group/field [&amp;&gt;*]:w-full [&amp;&gt;.sr-only]:w-auto flex w-full flex-col gap-3 data-[invalid=true]:text-destructive"
>
	<label
		data-slot="field-label"
		class="group/field-label peer/field-label has-[&gt;[data-slot=field]]:w-full has-[&gt;[data-slot=field]]:flex-col has-[&gt;[data-slot=field]]:rounded-md has-[&gt;[data-slot=field]]:border [&amp;&gt;*]:data-[slot=field]:p-4 sr-only flex w-fit items-center gap-2 text-sm leading-snug font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 group-data-[disabled=true]/field:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50 has-data-[state=checked]:border-primary has-data-[state=checked]:bg-primary/5 dark:has-data-[state=checked]:bg-primary/10"
		for="notion-prompt">Prompt</label
	>
	<div
		data-slot="input-group"
		role="group"
		class="group/input-group has-[&gt;textarea]:h-auto has-[&gt;[data-align=inline-start]]:[&amp;&gt;input]:pl-2 has-[&gt;[data-align=inline-end]]:[&amp;&gt;input]:pr-2 has-[&gt;[data-align=block-start]]:h-auto has-[&gt;[data-align=block-start]]:flex-col has-[&gt;[data-align=block-start]]:[&amp;&gt;input]:pb-3 has-[&gt;[data-align=block-end]]:h-auto has-[&gt;[data-align=block-end]]:flex-col has-[&gt;[data-align=block-end]]:[&amp;&gt;input]:pt-3 relative flex h-9 w-full min-w-0 items-center rounded-md border border-input shadow-xs transition-[color,box-shadow] outline-none has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-[3px] has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50 has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-destructive/20 dark:bg-input/30 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40"
	>
		<textarea
			data-slot="input-group-control"
			class="flex field-sizing-content min-h-16 w-full flex-1 resize-none rounded-none border-0 border-input bg-transparent px-3 py-3 text-base shadow-none transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-0 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:bg-transparent dark:aria-invalid:ring-destructive/40"
			id="notion-prompt"
			placeholder="Ask, search, or make anything..."
		></textarea>
		<div
			role="group"
			data-slot="input-group-addon"
			data-align="block-start"
			class="[&amp;&gt;svg:not([class*='size-'])]:size-4 [&amp;&gt;kbd]:rounded-[calc(var(--radius)-5px)] group-has-[&gt;input]/input-group:pt-2.5 order-first flex h-auto w-full cursor-text items-center justify-start gap-2 px-3 py-1.5 pt-3 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [.border-b]:pb-3"
		>
			<button
				data-slot="tooltip-trigger"
				class="[&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 [&amp;_svg]:shrink-0 has-[&gt;svg]:px-2.5 flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-full border bg-background px-2.5 py-2 text-sm font-medium whitespace-nowrap shadow-none transition-transform outline-none hover:bg-accent hover:text-accent-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:border-input dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40"
				type="button"
				data-size="sm"
				aria-haspopup="dialog"
				aria-expanded="false"
				aria-controls="radix-_r_oq_"
				data-state="closed"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="tabler-icon tabler-icon-at"
					><path d="M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path
						d="M16 12v1.5a2.5 2.5 0 0 0 5 0v-1.5a9 9 0 1 0 -5.5 8.28"
					></path></svg
				> Add context</button
			>
			<div class="no-scrollbar -m-1.5 flex gap-1 overflow-y-auto p-1.5"></div>
		</div>
		<div
			role="group"
			data-slot="input-group-addon"
			data-align="block-end"
			class="[&amp;&gt;svg:not([class*='size-'])]:size-4 [&amp;&gt;kbd]:rounded-[calc(var(--radius)-5px)] group-has-[&gt;input]/input-group:pb-2.5 order-last flex h-auto w-full cursor-text items-center justify-start gap-1 px-3 py-1.5 pb-3 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [.border-t]:pt-3"
		>
			<button
				data-slot="tooltip-trigger"
				class="[&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 [&amp;_svg]:shrink-0 has-[&gt;svg]:p-0 flex size-8 shrink-0 items-center justify-center gap-2 rounded-full p-0 text-sm font-medium whitespace-nowrap shadow-none transition-all outline-none hover:bg-accent hover:text-accent-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:hover:bg-accent/50 dark:aria-invalid:ring-destructive/40"
				type="button"
				data-size="icon-sm"
				aria-label="Attach file"
				data-state="closed"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="tabler-icon tabler-icon-paperclip"
					><path
						d="M15 7l-6.5 6.5a1.5 1.5 0 0 0 3 3l6.5 -6.5a3 3 0 0 0 -6 -6l-6.5 6.5a4.5 4.5 0 0 0 9 9l6.5 -6.5"
					></path></svg
				></button
			><button
				data-slot="tooltip-trigger"
				class="[&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 [&amp;_svg]:shrink-0 has-[&gt;svg]:px-2.5 flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-full px-2.5 py-2 text-sm font-medium whitespace-nowrap shadow-none transition-all outline-none hover:bg-accent hover:text-accent-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:hover:bg-accent/50 dark:aria-invalid:ring-destructive/40"
				type="button"
				data-size="sm"
				id="radix-_r_ot_"
				aria-haspopup="menu"
				aria-expanded="false"
				data-state="closed">Auto</button
			><button
				data-slot="dropdown-menu-trigger"
				class="[&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 [&amp;_svg]:shrink-0 has-[&gt;svg]:px-2.5 flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-full px-2.5 py-2 text-sm font-medium whitespace-nowrap shadow-none transition-all outline-none hover:bg-accent hover:text-accent-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:hover:bg-accent/50 dark:aria-invalid:ring-destructive/40"
				type="button"
				data-size="sm"
				id="radix-_r_p0_"
				aria-haspopup="menu"
				aria-expanded="false"
				data-state="closed"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="tabler-icon tabler-icon-world"
					><path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path><path d="M3.6 9h16.8"></path><path
						d="M3.6 15h16.8"
					></path><path d="M11.5 3a17 17 0 0 0 0 18"></path><path d="M12.5 3a17 17 0 0 1 0 18"
					></path></svg
				> All Sources</button
			><button
				data-slot="button"
				class="[&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 [&amp;_svg]:shrink-0 has-[&gt;svg]:p-0 ml-auto flex size-8 shrink-0 items-center justify-center gap-2 rounded-full bg-primary p-0 text-sm font-medium whitespace-nowrap text-primary-foreground shadow-none transition-all outline-none hover:bg-primary/90 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40"
				type="button"
				data-size="icon-sm"
				aria-label="Send"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="tabler-icon tabler-icon-arrow-up"
					><path d="M12 5l0 14"></path><path d="M18 11l-6 -6"></path><path d="M6 11l6 -6"
					></path></svg
				></button
			>
		</div>
	</div>
</div>
