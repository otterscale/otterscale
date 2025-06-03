import type { OptionType, AncestralOptionType } from './types'

import { OptionManager } from './utils.svelte'

import { DropdownMenu as DropdownMenuPrimitive } from "bits-ui";
import Root from './select.svelte';
import Trigger from './select-trigger.svelte';
import Content from './select-content.svelte';
import Group from './select-group.svelte';
import GroupHeading from './select-group-heading.svelte';
import Label from './select-label.svelte';
import Item from './select-item.svelte';
import Check from './select-check.svelte';
import Shortcut from './select-shortcut.svelte';
import Separator from './select-separator.svelte';
import SubTrigger from './select-sub-trigger.svelte';
import SubContent from './select-sub-content.svelte';

const Sub = DropdownMenuPrimitive.Sub;

export {
    Root,
    Sub,
    Trigger,
    Content,
    Group,
    GroupHeading,
    Label,
    Item,
    Shortcut,
    Check,
    Separator,
    SubTrigger,
    SubContent,
    //
    type OptionType,
    type AncestralOptionType,
    //
    OptionManager
};
