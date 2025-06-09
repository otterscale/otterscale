type OptionType = {
    value: any;
    label: string;
    icon?: string;
    information?: string;
    enabled?: boolean;
};

type valueSetterType = (option: OptionType) => void;

export type {
    OptionType,
    valueSetterType,
}