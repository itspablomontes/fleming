import { type UseMutationOptions, useMutation } from "@tanstack/react-query";

import { approveConsent, denyConsent, revokeConsent } from "../api";

export const useApproveConsent = (
	options?: UseMutationOptions<void, Error, string>,
) =>
	useMutation<void, Error, string>({
		mutationFn: approveConsent,
		...options,
	});

export const useDenyConsent = (
	options?: UseMutationOptions<void, Error, string>,
) =>
	useMutation<void, Error, string>({
		mutationFn: denyConsent,
		...options,
	});

export const useRevokeConsent = (
	options?: UseMutationOptions<void, Error, string>,
) =>
	useMutation<void, Error, string>({
		mutationFn: revokeConsent,
		...options,
	});
