import { create } from "zustand";
import type { TimelineEvent } from "@/features/timeline/types";
import { useEditStore } from "./edit-store";
import { useUploadStore } from "./upload-store";

interface LinkModeState {
	isLinkMode: boolean;
	linkSource: TimelineEvent | null;
	pendingTarget: TimelineEvent | null;
}

interface CoordinatorState extends LinkModeState {
	canUpload: () => boolean;
	canEdit: () => boolean;
	isAnyOperationActive: () => boolean;
	startEdit: (event: TimelineEvent) => void;
	startUpload: () => void;
	resetAll: () => void;
	// Link mode actions
	startLinkMode: () => void;
	cancelLinkMode: () => void;
	selectLinkSource: (event: TimelineEvent) => void;
	selectLinkTarget: (event: TimelineEvent) => void;
	clearPendingTarget: () => void;
	completeLinking: () => void;
}

export const useTimelineCoordinator = create<CoordinatorState>((set, get) => ({
	// Link mode state
	isLinkMode: false,
	linkSource: null,
	pendingTarget: null,

	// Existing computed state
	canUpload: () => !useEditStore.getState().isEditing && !get().isLinkMode,
	canEdit: () => !useUploadStore.getState().isUploading && !get().isLinkMode,
	isAnyOperationActive: () =>
		useEditStore.getState().isEditing ||
		useUploadStore.getState().isUploading ||
		get().isLinkMode,

	// Existing actions
	startEdit: (event) => {
		get().cancelLinkMode();
		useUploadStore.getState().reset();
		useEditStore.getState().startEdit(event);
	},
	startUpload: () => {
		get().cancelLinkMode();
		useEditStore.getState().cancelEdit();
	},
	resetAll: () => {
		useUploadStore.getState().reset();
		useEditStore.getState().cancelEdit();
		set({ isLinkMode: false, linkSource: null, pendingTarget: null });
	},

	// Link mode actions
	startLinkMode: () => {
		useUploadStore.getState().reset();
		useEditStore.getState().cancelEdit();
		set({ isLinkMode: true, linkSource: null, pendingTarget: null });
	},
	cancelLinkMode: () => {
		set({ isLinkMode: false, linkSource: null, pendingTarget: null });
	},
	selectLinkSource: (event) => {
		set({ linkSource: event, pendingTarget: null });
	},
	selectLinkTarget: (event) => {
		const { linkSource } = get();
		if (!linkSource || linkSource.id === event.id) return;
		set({ pendingTarget: event });
	},
	clearPendingTarget: () => {
		set({ pendingTarget: null });
	},
	completeLinking: () => {
		set({ isLinkMode: false, linkSource: null, pendingTarget: null });
	},
}));
