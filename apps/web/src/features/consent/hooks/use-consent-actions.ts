import { type UseMutationOptions, useMutation, useQueryClient } from "@tanstack/react-query";

import { approveConsent, denyConsent, revokeConsent } from "../api";

export const useApproveConsent = (
	options?: UseMutationOptions<void, Error, string>,
) => {
	const queryClient = useQueryClient();

	return useMutation<void, Error, string>({
		mutationFn: approveConsent,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["active-grants"] });
			queryClient.invalidateQueries({ queryKey: ["my-grants"] });
			queryClient.invalidateQueries({ queryKey: ["audit-entries"] });
		},
		...options,
	});
};

export const useDenyConsent = (
	options?: UseMutationOptions<void, Error, string>,
) => {
	const queryClient = useQueryClient();

	return useMutation<void, Error, string>({
		mutationFn: denyConsent,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["pending-requests"] });
			queryClient.invalidateQueries({ queryKey: ["audit-entries"] });
		},
		...options,
	});
};

export const useRevokeConsent = (
	options?: UseMutationOptions<void, Error, string>,
) => {
	const queryClient = useQueryClient();

	return useMutation<void, Error, string>({
		mutationFn: revokeConsent,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["active-grants"] });
			queryClient.invalidateQueries({ queryKey: ["my-grants"] });
			queryClient.invalidateQueries({ queryKey: ["audit-entries"] });
		},
		...options,
	});
};
