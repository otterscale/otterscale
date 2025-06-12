import type { InputType as MultipleInputType } from '../multiple/types'

type InputType = Exclude<MultipleInputType, 'color'>
type BooleanOption = {
    value: any,
    label: any,
    icon: string
}

export type {
    BooleanOption, InputType
}
