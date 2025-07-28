import type { InputType as MultipleInputType } from '../multiple/types'

type InputType = MultipleInputType | 'password'
type BooleanOption = {
    value: any,
    label: any,
    icon: string
}
type UnitType = {
    value: any,
    label: any,
    icon?: string
}

export type {
    BooleanOption, InputType, UnitType
}
