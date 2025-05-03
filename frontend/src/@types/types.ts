export interface SelectionItem {
    title: string;
    options: { title: string; icon: string; }[];
}

export interface SimulationItem {
    name: string,
    index: number,
    steps: SimulationStep[],
}

export interface SimulationStep {
    role: Role,
    text: string,
    index: number,
}

export interface ChatMessage {
    role: Role;
    text: string;
    formattedText?: string;
}

export type InputType = "block" | "text" | "speech";

export type Role = "aircraft" | "tower";
