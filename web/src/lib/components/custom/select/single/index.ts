import type { OptionType, valueSetterType } from './types'

import { OptionManager } from './utils.svelte';

import { Popover as PopoverPrimitive } from "bits-ui";
import Root from './select.svelte'
import Trigger from './select-trigger.svelte'
import Content from './select-content.svelte'
import Options from './select-options.svelte'
import Group from './select-group.svelte'
import Input from './select-input.svelte'
import Empty from './select-empty.svelte'
import List from './select-list.svelte'
import Item from './select-item.svelte'
import Check from './select-check.svelte'
import Shortcut from './select-shortcut.svelte'
// const Root = PopoverPrimitive.Root;
const Close = PopoverPrimitive.Close;

export {
    Root,
    Trigger,
    Content,
    Options,
    Group,
    Input,
    Empty,
    List,
    Item,
    Shortcut,
    Check,
    Close,
    //
    type OptionType,
    type valueSetterType,
    //
    OptionManager,
};
