import Menubar from './tabs-menubar.svelte';
import MenuItem from './tabs-menuitem.svelte';
import Root from './tabs.svelte';

import {
	Trigger,
	Menu,
	Content as MenuContent,
	Separator as MenuSeparator
} from '$lib/components/ui/menubar/index';
import { Content } from '$lib/components/ui/tabs';

export { Root, Content, Trigger, Menubar, Menu, MenuContent, MenuItem, MenuSeparator };
