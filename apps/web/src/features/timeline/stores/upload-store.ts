import { create } from "zustand";

interface UploadState {
	file: File | null;
	isUploading: boolean;
	uploadStatus: string;
	error: string | null;
	setFile: (file: File | null) => void;
	startUpload: () => void;
	finishUpload: () => void;
	setStatus: (status: string) => void;
	setError: (error: string | null) => void;
	reset: () => void;
}

export const useUploadStore = create<UploadState>((set) => ({
	file: null,
	isUploading: false,
	uploadStatus: "",
	error: null,
	setFile: (file) => set({ file }),
	startUpload: () => set({ isUploading: true }),
	finishUpload: () => set({ isUploading: false, uploadStatus: "" }),
	setStatus: (status) => set({ uploadStatus: status }),
	setError: (error) => set({ error }),
	reset: () =>
		set({
			file: null,
			isUploading: false,
			uploadStatus: "",
			error: null,
		}),
}));
