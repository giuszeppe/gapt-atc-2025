export interface SelectionItem {
    title: string;
    options: { title: string; icon: string; value: Mode | AdvancementType | InputType | Scenario | Role }[];
}

export interface SimulationStep {
    role: "aircraft" | "tower",
    text: string,
    index: number,
}

export interface ChatMessage {
    role: "aircraft" | "tower";
    text: string;
    formattedText?: string;
}

export enum Mode {
    Singleplayer = "singleplayer",
    Multiplayer = "multiplayer",
}

export enum AdvancementType {
    Continuous = "continuous",
    ClickToStep = "click to step",
}

export enum InputType {
    Block = "block",
    Text = "text",
    Speech = "speech",
}

export enum Scenario {
    Takeoff = "takeoff",
    Enroute = "enroute",
    Landing = "landing",
}

export enum Role {
    Aircraft = "aircraft",
    Tower = "tower",
}