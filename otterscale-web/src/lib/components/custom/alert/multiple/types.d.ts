import type { AlertType } from '../single/types'
import type { AlertVariant } from './alert.svelte'

type VariantGetterType = (level: any) => AlertVariant
type ValueType = {
    index: number
}

export type {
    AlertType,
    AlertVariant, ValueType, VariantGetterType
}
