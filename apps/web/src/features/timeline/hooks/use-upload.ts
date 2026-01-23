import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { addEvent } from "../api";

export function useUpload() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: addEvent,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["timeline"] });
			toast.success("Document uploaded and encrypted successfully");
		},
		onError: (error) => {
			console.error("Upload error:", error);
			toast.error("Failed to upload document. Please try again.");
		},
	});
}
