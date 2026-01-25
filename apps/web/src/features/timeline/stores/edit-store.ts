import { create } from "zustand";
import type { TimelineEvent } from "@/features/timeline/types";

interface EditState {
	editEvent: TimelineEvent | null;
	isEditing: boolean;
	startEdit: (event: TimelineEvent) => void;
	setEditEvent: (event: TimelineEvent | null) => void;
	cancelEdit: () => void;
}

export const useEditStore = create<EditState>((set) => ({
	editEvent: null,
	isEditing: false,
	startEdit: (event) => set({ editEvent: event, isEditing: true }),
	setEditEvent: (event) => set({ editEvent: event, isEditing: Boolean(event) }),
	cancelEdit: () => set({ editEvent: null, isEditing: false }),
}));
