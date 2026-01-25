import { create } from "zustand";
import type { TimelineEvent } from "@/features/timeline/types";
import { useEditStore } from "./edit-store";
import { useUploadStore } from "./upload-store";

interface CoordinatorState {
	canUpload: () => boolean;
	canEdit: () => boolean;
	isAnyOperationActive: () => boolean;
	startEdit: (event: TimelineEvent) => void;
	startUpload: () => void;
	resetAll: () => void;
}

export const useTimelineCoordinator = create<CoordinatorState>(() => ({
	canUpload: () => !useEditStore.getState().isEditing,
	canEdit: () => !useUploadStore.getState().isUploading,
	isAnyOperationActive: () =>
		useEditStore.getState().isEditing || useUploadStore.getState().isUploading,
	startEdit: (event) => {
		useUploadStore.getState().reset();
		useEditStore.getState().startEdit(event);
	},
	startUpload: () => {
		useEditStore.getState().cancelEdit();
	},
	resetAll: () => {
		useUploadStore.getState().reset();
		useEditStore.getState().cancelEdit();
	},
}));
