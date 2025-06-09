type OptionType = {
    value: any;
    label: string;
    icon?: string;
    information?: string;
    enabled?: boolean;
};

type valueSetterType = (value: any) => void;
type valueGetterType = () => any;

export type {
    OptionType,
    valueSetterType,
    valueGetterType,
}